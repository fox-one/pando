package payee

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	walletz core.WalletService,
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
		core.ActionCatCreate: cat.HandleCreate(collaterals, oracles),
		core.ActionCatSupply: cat.HandleSupply(collaterals),
		// vat
		core.ActionVatOpen:     vat.HandleOpen(collaterals, vaults, wallets),
		core.ActionVatDeposit:  vat.HandleDeposit(collaterals, vaults, wallets),
		core.ActionVatWithdraw: vat.HandleWithdraw(collaterals, vaults, wallets),
		core.ActionVatPayback:  vat.HandlePayback(collaterals, vaults, wallets),
		core.ActionVatGenerate: vat.HandleGenerated(collaterals, vaults, wallets),
		// flip
		core.ActionFlipKick: flip.HandleKick(collaterals, vaults, flips),
		core.ActionFlipBid:  flip.HandleBid(collaterals, vaults, flips, wallets),
		core.ActionFlipDeal: flip.HandleDeal(collaterals, flips, wallets),
		// oracle
		core.ActionOraclePoke: oracle.HandlePoke(collaterals, oracles),
		core.ActionOracleStep: oracle.HandleStep(oracles),
		// proposal
		core.ActionProposalMake:  proposal.HandleMake(proposals, walletz, parliaments, system),
		core.ActionProposalShout: proposal.HandleShout(proposals, parliaments, system),
		core.ActionProposalVote:  proposal.HandleVote(proposals, parliaments, walletz, system),
	}

	return &Payee{
		wallets:      wallets,
		property:     property,
		oraclez:      oraclez,
		transactions: transactions,
		system:       system,
		actions:      actions,
	}
}

type Payee struct {
	wallets      core.WalletStore
	property     property.Store
	oraclez      core.OracleService
	transactions core.TransactionStore
	system       *core.System
	actions      map[core.Action]maker.HandlerFunc
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
	req := requestFromOutput(output)

	// 1, parse oracle message
	if body, err := w.oraclez.Parse(message); err == nil {
		req.Action = core.ActionOraclePoke
		req.Body = body
		req.Gov = true
		return w.handleRequest(req.WithContext(ctx))
	}

	// 2. decode tx message
	if body, err := mtg.Decrypt(message, w.system.PrivateKey); err == nil {
		message = body
	}

	if payload, err := core.DecodeTransactionAction(message); err == nil {
		if req.Body, err = mtg.Scan(payload.Body, &req.Action); err == nil {
			if follow, _ := uuid.FromBytes(payload.FollowID); follow != uuid.Zero {
				req.FollowID = follow.String()
			}

			return w.handleRequest(req.WithContext(ctx))
		}
	}

	return nil
}

func (w *Payee) handleRequest(r *maker.Request) error {
	ctx := r.Context()
	log := logger.FromContext(ctx).WithField("action", r.Action.String())

	h, ok := w.actions[r.Action]
	if !ok {
		log.Debugf("handler not found")
		return nil
	}

	tx := r.Tx()

	if err := h(r); err != nil {
		var e maker.Error
		if !errors.As(err, &e) {
			return err
		}

		if r.Sender != "" && maker.ShouldRefund(e.Flag) {
			memo := core.TransferAction{
				ID:     r.FollowID,
				Source: e.Error(),
			}.Encode()

			transfer := &core.Transfer{
				TraceID:   uuid.Modify(r.TraceID, memo),
				AssetID:   r.AssetID,
				Amount:    r.Amount,
				Memo:      memo,
				Threshold: 1,
				Opponents: []string{r.Sender},
			}

			if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{transfer}); err != nil {
				log.WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}
		}

		tx.Status = core.TransactionStatusAbort
		tx.Message = e.Msg
	} else {
		tx.Status = core.TransactionStatusOk
	}

	tx.Parameters, _ = json.Marshal(r.Values())
	if err := w.transactions.Create(ctx, tx); err != nil {
		log.WithError(err).Errorln("transactions.Create")
		return err
	}

	if r.Next != nil {
		return w.handleRequest(r.Next)
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

func requestFromOutput(output *core.Output) *maker.Request {
	return &maker.Request{
		Now:      output.CreatedAt,
		Version:  output.ID,
		TraceID:  output.TraceID,
		Sender:   output.Sender,
		FollowID: output.TraceID,
		AssetID:  output.AssetID,
		Amount:   output.Amount,
	}
}
