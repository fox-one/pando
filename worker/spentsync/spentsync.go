package spentsync

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
)

func New(
	wallets core.WalletStore,
	notifier core.Notifier,
) *SpentSync {
	return &SpentSync{
		wallets:  wallets,
		notifier: notifier,
	}
}

type SpentSync struct {
	wallets  core.WalletStore
	notifier core.Notifier
}

func (w *SpentSync) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "SpentSync")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 500 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *SpentSync) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	const limit = 100
	transfers, err := w.wallets.ListTransfers(ctx, core.TransferStatusHandled, limit)
	if err != nil {
		log.WithError(err).Errorln("wallets.ListNotPassedTransfers")
		return err
	}

	if len(transfers) == 0 {
		return errors.New("EOF")
	}

	for _, transfer := range transfers {
		_ = w.handleTransfer(ctx, transfer)
	}

	return nil
}

func (w *SpentSync) handleTransfer(ctx context.Context, transfer *core.Transfer) error {
	log := logger.FromContext(ctx).WithField("trace", transfer.TraceID)

	output, err := w.wallets.FindSpentBy(ctx, transfer.AssetID, transfer.TraceID)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil
		}

		log.WithError(err).Errorln("wallets.ListSpentBy")
		return err
	}

	if output.State != mixin.UTXOStateSpent {
		return nil
	}

	signedTx := output.SignedTx
	if signedTx == "" {
		return nil
	}

	if err := w.notifier.Snapshot(ctx, transfer, signedTx); err != nil {
		log.WithError(err).Errorln("notifier.Snapshot")
		return err
	}

	if err := w.wallets.UpdateTransfer(ctx, transfer, core.TransferStatusPassed); err != nil {
		log.WithError(err).Errorln("wallets.UpdateTransfer")
		return err
	}

	return nil
}
