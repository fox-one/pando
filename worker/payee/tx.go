package payee

import (
	"context"
	"errors"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pando/worker/payee/actions"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
)

func (w *Payee) handleTransaction(ctx context.Context, output *core.Output, body []byte) error {
	log := logger.FromContext(ctx)

	var (
		userID, followID uuid.UUID
		action           mtg.BitInt
		targetID         uuid.UUID
		err              error
	)

	body, err = mtg.Scan(body, &userID, &followID, &action, &targetID)
	if err != nil {
		// memo 无法解析，跳过
		return nil
	}

	handler, ok := w.actions[int(action)]
	if !ok {
		return nil
	}

	tx := &maker.Tx{
		Version:  output.ID,
		TraceID:  output.TraceID,
		AssetID:  output.AssetID,
		Amount:   output.Amount,
		Sender:   userID.String(),
		FollowID: followID.String(),
		Action:   int(action),
		TargetID: targetID.String(),
		Now:      output.CreatedAt,
	}

	t, err := w.transactions.Find(ctx, tx.TraceID)
	if err != nil {
		if !store.IsErrNotFound(err) {
			return err
		}

		t = convertTransaction(tx)

		data, err := handler.Build(ctx, tx, body)
		if err != nil {
			if errors.Is(err, actions.ErrBuildAbort) {
				return nil
			}

			if status, _, ok := maker.Unwrap(err); ok {
				t.Status = status
			} else {
				return err
			}
		}

		t.Data = data

		if err := w.wallets.CreateTransfers(ctx, tx.Transfers); err != nil {
			log.WithError(err).Errorln("wallets.CreateTransfers")
			return err
		}

		if err := w.transactions.Create(ctx, t); err != nil {
			log.WithError(err).Errorln("transactions.Create")
			return err
		}
	}

	if t.Status == core.TransactionStatusOK {
		if err := handler.Apply(ctx, tx, t.Data); err != nil {
			return err
		}
	} else {
		return w.refundTransaction(ctx, output, t)
	}

	return nil
}

func convertTransaction(tx *maker.Tx) *core.Transaction {
	return &core.Transaction{
		CreatedAt: tx.Now,
		TraceID:   tx.TraceID,
		FollowID:  tx.FollowID,
		UserID:    tx.Sender,
		Action:    tx.Action,
		TargetID:  tx.TargetID,
		Status:    core.TransactionStatusOK,
	}
}
