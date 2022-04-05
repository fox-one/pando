package stater

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/jinzhu/now"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

const checkpointKey = "stater_transaction_checkpoint"

func New(
	collaterals core.CollateralStore,
	transactions core.TransactionStore,
	vaults core.VaultStore,
	stats core.StatStore,
	flips core.FlipStore,
	properties property.Store,
	assetz core.AssetService,
) *Stater {
	return &Stater{
		collaterals:  collaterals,
		transactions: transactions,
		vaults:       vaults,
		stats:        stats,
		flips:        flips,
		properties:   properties,
		assetz:       cacheAssetPrice(assetz),
	}
}

type Stater struct {
	collaterals  core.CollateralStore
	transactions core.TransactionStore
	vaults       core.VaultStore
	stats        core.StatStore
	flips        core.FlipStore
	properties   property.Store
	assetz       core.AssetService
}

func (w *Stater) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "stater")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (w *Stater) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	v, err := w.properties.Get(ctx, checkpointKey)
	if err != nil {
		log.WithError(err).Errorln("properties.Get", err)
		return err
	}

	var (
		fromID = v.Int64()
		limit  = 500
	)

	transactions, err := w.transactions.List(ctx, fromID, limit)
	if err != nil {
		log.WithError(err).Errorln("transaction.List")
		return err
	}

	if len(transactions) == 0 {
		return errors.New("EOF")
	}

	for _, t := range transactions {
		if t.Status != core.TransactionStatusOk {
			continue
		}

		log := log.WithField("tx", t.ID)
		ctx := logger.WithContext(ctx, log)

		if err := w.handleTransaction(ctx, t); err != nil {
			log.WithError(err).Errorln("handleTransaction")
			return err
		}
	}

	if err := w.properties.Save(ctx, checkpointKey, fromID); err != nil {
		log.WithError(err).Errorln("properties.Save")
		return err
	}

	return nil
}

func (w *Stater) handleTransaction(ctx context.Context, t *core.Transaction) error {
	log := logger.FromContext(ctx)

	var parameters []interface{}
	_ = t.Parameters.Unmarshal(&parameters)

	var (
		collateralID string
		ink, debt    decimal.Decimal
	)

	switch t.Action {
	case core.ActionVatOpen:
		collateralID = cast.ToString(parameters[0])
		debt = number.Decimal(cast.ToString(parameters[1]))
		ink = t.Amount
	case core.ActionVatDeposit, core.ActionVatWithdraw, core.ActionVatPayback, core.ActionVatGenerate:
		vaultID := cast.ToString(parameters[0])

		vault, err := w.vaults.Find(ctx, vaultID)
		if err != nil {
			log.WithError(err).Errorln("vault.Find")
			return err
		}

		if vault.ID == 0 {
			return fmt.Errorf("vault with id %s not found", vaultID)
		}

		event, err := w.vaults.FindEvent(ctx, vault.TraceID, t.Version)
		if err != nil {
			log.WithError(err).Errorln("vault.FindEvent")
			return err
		}

		if event.ID == 0 {
			return fmt.Errorf("vault event %s,%d not found", vault.TraceID, t.Version)
		}

		collateralID = vault.CollateralID
		ink = event.Dink
		debt = event.Debt
	case core.ActionFlipDeal:
		flipID := cast.ToString(parameters[0])
		flip, err := w.flips.Find(ctx, flipID)
		if err != nil {
			log.WithError(err).Errorln("flips.Find")
			return err
		}

		if flip.ID == 0 {
			return fmt.Errorf("flip with id %s not found", flipID)
		}

		collateralID = flip.CollateralID
		ink = flip.Lot.Neg()
		debt = flip.Bid.Neg()
	case core.ActionCatGain:
		collateralID = cast.ToString(parameters[0])
		debt = number.Decimal(cast.ToString(parameters[1]))
	case core.ActionCatFold:
		collateralID = cast.ToString(parameters[0])
	}

	if collateralID == "" {
		return nil
	}

	cats, err := w.listCollaterals(ctx, collateralID, t.CreatedAt)
	if err != nil {
		log.WithError(err).Errorln("listCollaterals")
		return err
	}

	if len(cats) == 0 {
		return fmt.Errorf("listCollaterals: no collaterals found for %s", collateralID)
	}

	for _, cat := range cats {
		if err := w.updateCollateral(ctx, cat, t, ink, debt); err != nil {
			log.WithError(err).Errorln("updateCollateral")
			return err
		}
	}

	return nil
}

func (w *Stater) listCollaterals(ctx context.Context, id string, at time.Time) ([]*core.Collateral, error) {
	var (
		cats []*core.Collateral
		err  error
	)

	if uuid.IsNil(id) {
		cats, err = w.collaterals.List(ctx)
	} else {
		var cat *core.Collateral
		if cat, err = w.collaterals.Find(ctx, id); err == nil {
			if cat.ID > 0 {
				cats = append(cats, cat)
			} else {
				err = fmt.Errorf("collateral with id %s not found", id)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	var idx int
	for _, cat := range cats {
		if cat.CreatedAt.After(at) {
			continue
		}

		cats[idx] = cat
		idx++
	}

	return cats[:idx], nil
}

func (w *Stater) updateCollateral(ctx context.Context, cat *core.Collateral, t *core.Transaction, ink, debt decimal.Decimal) error {
	log := logger.FromContext(ctx)

	n := now.New(t.CreatedAt)
	stat, err := w.stats.Find(ctx, cat.TraceID, n.BeginningOfDay())
	if err != nil {
		log.WithError(err).Errorln("stats.Find")
		return err
	}

	if stat.Version >= t.Version {
		return nil
	}

	inkPrice, err := w.assetz.ReadPrice(ctx, cat.Gem, n.EndOfDay())
	if err != nil {
		log.WithError(err).Errorf("assetz.ReadPrice %s", cat.Gem)
		return err
	}

	debtPrice, err := w.assetz.ReadPrice(ctx, cat.Dai, n.EndOfDay())
	if err != nil {
		log.WithError(err).Errorf("assetz.ReadPrice %s", cat.Dai)
		return err
	}

	s := &core.Stat{
		Version:      t.Version,
		CollateralID: cat.TraceID,
		Date:         n.BeginningOfDay(),
		Gem:          cat.Gem,
		Dai:          cat.Dai,
		Ink:          stat.Ink.Add(ink),
		Debt:         stat.Debt.Add(debt),
		InkPrice:     inkPrice,
		DebtPrice:    debtPrice,
	}

	if err := w.stats.Save(ctx, s); err != nil {
		log.WithError(err).Errorln("stats.Save")
		return err
	}

	return nil
}
