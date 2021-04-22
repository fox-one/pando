package cashier

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func New(
	wallets core.WalletStore,
	walletz core.WalletService,
	system *core.System,
) *Cashier {
	return &Cashier{
		wallets: wallets,
		walletz: walletz,
		system:  system,
	}
}

type Cashier struct {
	wallets core.WalletStore
	walletz core.WalletService
	system  *core.System
}

func (w *Cashier) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "cashier")
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

func (w *Cashier) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	transfers, err := w.wallets.ListPendingTransfers(ctx)
	if err != nil {
		log.WithError(err).Errorln("list transfers")
		return err
	}

	if len(transfers) == 0 {
		return errors.New("EOF")
	}

	g := errgroup.Group{}
	sem := semaphore.NewWeighted(5)

	for idx := range transfers {
		transfer := transfers[idx]

		if err := sem.Acquire(ctx, 1); err != nil {
			return g.Wait()
		}

		g.Go(func() error {
			defer sem.Release(1)
			return w.handleTransfer(ctx, transfer)
		})
	}

	return g.Wait()
}

func (w *Cashier) handleTransfer(ctx context.Context, transfer *core.Transfer) error {
	log := logger.FromContext(ctx).WithField("transfer", transfer.TraceID)

	if valid := !transfer.Handled && transfer.Assigned; !valid {
		log.Panicln("invalid transfer")
	}

	outputs, err := w.wallets.ListSpentBy(ctx, transfer.AssetID, transfer.TraceID)
	if err != nil {
		log.WithError(err).Errorln("wallets.ListSpentBy")
		return err
	}

	if len(outputs) == 0 {
		log.Errorln("cannot spent transfer with empty outputs")
		return nil
	}

	return w.spend(ctx, outputs, transfer)
}

func (w *Cashier) spend(ctx context.Context, outputs []*core.Output, transfer *core.Transfer) error {
	if tx, err := w.walletz.Spend(ctx, outputs, transfer); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("walletz.Spend")
		return err
	} else if tx != nil {
		// signature completed, prepare to send this tx to mixin mainnet
		if err := w.wallets.CreateRawTransaction(ctx, tx); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("wallets.CreateRawTransaction")
			return err
		}
	}

	transfer.Handled = true
	if err := w.wallets.UpdateTransfer(ctx, transfer); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("wallets.UpdateTransfer")
		return err
	}

	return nil
}
