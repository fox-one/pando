package payee

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/flip"
	"github.com/fox-one/pando/pkg/maker/oracle"
	"github.com/fox-one/pando/pkg/maker/proposal"
	"github.com/fox-one/pando/pkg/maker/sys"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

const (
	checkpointKey = "outputs_checkpoint"
)

func New(
	assets core.AssetStore,
	assetz core.AssetService,
	wallets core.WalletStore,
	transactions core.TransactionStore,
	proposals core.ProposalStore,
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	property property.Store,
	parliaments core.Parliament,
	oracles core.OracleStore,
	oraclez core.OracleService,
	system *core.System,
) *Payee {
	actions := map[core.Action]maker.HandlerFunc{
		// sys
		core.ActionSysWithdraw: sys.HandleWithdraw(wallets),
		// cat
		core.ActionCatEdit:   cat.HandleEdit(collaterals),
		core.ActionCatFold:   cat.HandleFold(collaterals),
		core.ActionCatCreate: cat.HandleCreate(collaterals, oracles, assets, assetz),
		core.ActionCatSupply: cat.HandleSupply(collaterals),
		// vat
		core.ActionVatOpen:     vat.HandleOpen(collaterals, vaults, transactions, wallets),
		core.ActionVatDeposit:  vat.HandleDeposit(collaterals, vaults, transactions, wallets),
		core.ActionVatWithdraw: vat.HandleWithdraw(collaterals, vaults, transactions, wallets),
		core.ActionVatPayback:  vat.HandlePayback(collaterals, vaults, transactions, wallets),
		core.ActionVatGenerate: vat.HandleGernerate(collaterals, vaults, transactions, wallets),
		// flip
		core.ActionFlipKick: flip.HandleKick(collaterals, vaults, flips, transactions, property),
		core.ActionFlipBid:  flip.HandleBid(collaterals, vaults, flips, transactions, wallets, property),
		core.ActionFlipDeal: flip.HandleDeal(collaterals, vaults, flips, transactions, wallets),
		core.ActionFlipOpt:  flip.HandleOpt(property),
		// oracle
		core.ActionOracleFeed: oracle.HandleFeed(collaterals, oracles),
		// proposal
		core.ActionProposalMake: proposal.HandleMake(proposals, parliaments, system),
	}

	actions[core.ActionProposalVote] = proposal.HandleVote(proposals, parliaments, actions, system)

	return &Payee{
		wallets:   wallets,
		proposals: proposals,
		property:  property,
		oraclez:   oraclez,
		system:    system,
		actions:   actions,
	}
}

type Payee struct {
	wallets   core.WalletStore
	property  property.Store
	proposals core.ProposalStore
	oraclez   core.OracleService
	system    *core.System
	actions   map[core.Action]maker.HandlerFunc
}

func (w *Payee) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "payee")
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

func (w *Payee) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	v, err := w.property.Get(ctx, checkpointKey)
	if err != nil {
		log.WithError(err).Errorln("property.Get", err)
		return err
	}

	const Limit = 500
	outputs, err := w.wallets.List(ctx, v.Int64(), Limit)
	if err != nil {
		log.WithError(err).Errorln("wallets.List")
		return err
	}

	if len(outputs) == 0 {
		return errors.New("EOF")
	}

	for _, u := range outputs {
		if err := w.handleOutput(ctx, u); err != nil {
			return err
		}

		if err := w.property.Save(ctx, checkpointKey, u.ID); err != nil {
			log.WithError(err).Errorln("property.Save", checkpointKey)
			return err
		}
	}

	return nil
}

func (w *Payee) handleOutput(ctx context.Context, output *core.Output) error {
	log := logger.FromContext(ctx).WithField("output", output.TraceID)
	ctx = logger.WithContext(ctx, log)

	message := decodeMemo(output.Memo)

	// 1, parse oracle message
	if ora, err := w.oraclez.Parse(message); err == nil {
		asset, _ := uuid.FromString(ora.AssetID)
		if b, err := mtg.Encode(asset, ora.Price, ora.PeekAt.Unix()); err == nil {
			r := &maker.Request{
				UTXO:   output,
				Action: core.ActionOracleFeed,
				Body:   b,
				Gov:    true,
			}

			return w.handleRequest(ctx, r)
		}

		return nil
	}

	// 2. decode group action
	if member, body, err := core.DecodeMemberAction(message, w.system.Members); err == nil {
		var action core.Action
		if body, err := mtg.Scan(body, &action); err == nil {
			r := &maker.Request{
				UTXO:   output,
				Action: action,
				Body:   body,
				UserID: member.ClientID,
			}

			return w.handleRequest(ctx, r)
		}

		return nil
	}

	// 3. decode tx message
	if body, err := mtg.Decrypt(message, w.system.PrivateKey); err == nil {
		var action core.Action
		if body, err = mtg.Scan(body, &action); err == nil {
			r := &maker.Request{
				UTXO:   output,
				Action: action,
				Body:   body,
				Gov:    false,
			}

			return w.handleRequest(ctx, r)
		}

		return nil
	}

	return nil
}

func (w *Payee) handleRequest(ctx context.Context, r *maker.Request) error {
	log := logger.FromContext(ctx).WithField("action", r.Action.String())

	h, ok := w.actions[r.Action]
	if !ok {
		log.Debugf("handler not found")
		return nil
	}

	if err := h(ctx, r); err != nil {
		var e maker.Error
		if !errors.As(err, &e) {
			return err
		}

		// refunds
		if r.UserID != "" {
			id := r.FollowID
			if id == "" {
				id = r.TraceID()
			}

			memo := core.TransferAction{
				ID:     id,
				Source: e.Error(),
			}.Encode()

			asset, amount := r.Payment()
			transfer := &core.Transfer{
				TraceID:   uuid.Modify(r.TraceID(), "Refund"),
				AssetID:   asset,
				Amount:    amount,
				Memo:      memo,
				Threshold: 1,
				Opponents: []string{r.UserID},
			}

			if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{transfer}); err != nil {
				log.WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}
		}
	}

	return nil
}

func decodeMemo(memo string) []byte {
	if b, err := base64.StdEncoding.DecodeString(memo); err == nil {
		return b
	}

	if b, err := base64.URLEncoding.DecodeString(memo); err == nil {
		return b
	}

	return []byte(memo)
}
