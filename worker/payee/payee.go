package payee

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pando/worker/payee/actions"
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
	notifier core.Notifier,
	parliament core.Parliament,
	oracles core.OracleStore,
	oraclez core.OracleService,
	system *core.System,
) *Payee {
	return &Payee{
		assets:       assets,
		assetz:       assetz,
		wallets:      wallets,
		transactions: transactions,
		collaterals:  collaterals,
		vaults:       vaults,
		flips:        flips,
		proposals:    proposals,
		property:     property,
		notifier:     notifier,
		parliament:   parliament,
		oracles:      oracles,
		oraclez:      oraclez,
		system:       system,
	}
}

type Payee struct {
	assets       core.AssetStore
	assetz       core.AssetService
	wallets      core.WalletStore
	collaterals  core.CollateralStore
	vaults       core.VaultStore
	flips        core.FlipStore
	transactions core.TransactionStore
	property     property.Store
	proposals    core.ProposalStore
	notifier     core.Notifier
	parliament   core.Parliament
	oracles      core.OracleStore
	oraclez      core.OracleService
	system       *core.System
	actions      map[int]actions.Handler
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
	if oracle, err := w.oraclez.Parse(message); err == nil {
		if err := w.oracles.Create(ctx, oracle); err != nil {
			log.WithError(err).Errorln("oracles.Create")
			return err
		}
	}

	// 2. decode group action
	if member, body, err := core.DecodeMemberAction(message, w.system.Members); err == nil {
		return w.handleMemberAction(ctx, output, member, body)
	}

	// 3. decode tx message
	if body, err := mtg.Decrypt(message, w.system.PrivateKey); err == nil {
		return w.handleTransaction(ctx, output, body)
	}

	return nil
}

func (w *Payee) refundTransaction(ctx context.Context, output *core.Output, tx *core.Transaction) error {
	msg, _ := maker.ErrorMsg(tx.Status)

	transfer := &core.Transfer{
		TraceID:   uuid.Modify(output.TraceID, "refund tx"),
		Opponents: []string{tx.UserID},
		Threshold: 1,
		AssetID:   output.AssetID,
		Amount:    output.UTXO.Amount,
		Memo: core.TransferAction{
			Module: "refund",
			ID:     tx.FollowID,
			Source: msg,
		}.Encode(),
	}

	if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{transfer}); err != nil {
		logger.FromContext(ctx).WithError(err).Errorf("wallets.CreateTransfers")
		return err
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
