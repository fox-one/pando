package txsender

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func New(wallets core.WalletStore) *Sender {
	return &Sender{
		wallets: wallets,
	}
}

type Sender struct {
	wallets core.WalletStore
}

func (w *Sender) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "txsender")
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
				dur = 300 * time.Millisecond
			}
		}
	}
}

func (w *Sender) run(ctx context.Context) error {
	log := logger.FromContext(ctx)
	const Limit = 20

	txs, err := w.wallets.ListPendingRawTransactions(ctx, Limit)
	if err != nil {
		log.WithError(err).Errorln("list raw transactions")
		return err
	}

	if len(txs) == 0 {
		return errors.New("EOF")
	}

	var g errgroup.Group
	for _, tx := range txs {
		tx := tx
		g.Go(func() error {
			return w.handleRawTransaction(ctx, tx)
		})
	}

	return g.Wait()
}

func (w *Sender) handleRawTransaction(ctx context.Context, tx *core.RawTransaction) error {
	log := logger.FromContext(ctx).WithField("trace_id", tx.TraceID)
	ctx = logger.WithContext(ctx, log)

	if err := w.submitRawTransaction(ctx, tx.Data); err != nil {
		return err
	}
	if err := w.wallets.ExpireRawTransaction(ctx, tx); err != nil {
		log.WithError(err).Errorln("wallets.ExpireRawTransaction")
		return err
	}
	return nil
}

func (w *Sender) submitRawTransaction(ctx context.Context, raw string) error {
	log := logger.FromContext(ctx)
	ctx = mixin.WithMixinNetHost(ctx, mixin.RandomMixinNetHost())

	if tx, err := mixin.SendRawTransaction(ctx, raw); err != nil {
		if mixin.IsErrorCodes(err, mixin.InvalidSignature) {
			return nil
		}

		log.WithError(err).Errorln("SendRawTransaction failed")
		return err
	} else if tx.Snapshot != nil {
		return nil
	}

	var txHash mixin.Hash
	if tx, err := mixin.TransactionFromRaw(raw); err == nil {
		txHash, _ = tx.TransactionHash()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return errors.New("mixin net snapshot not generated")
		case <-time.After(dur):
			if tx, err := mixin.GetTransaction(ctx, txHash); err != nil {
				log.WithError(err).Errorln("GetTransaction failed")
				return err
			} else if tx.Snapshot != nil {
				return nil
			}
			dur = time.Second
		}
	}
}
