package keeper

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/worker/keeper/wallet"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func New(
	cats core.CollateralStore,
	oracles core.OracleStore,
	vaults core.VaultStore,
	walletz core.WalletService,
	notifier core.Notifier,
	system *core.System,
) *Keeper {
	return &Keeper{
		cats:     cats,
		oracles:  oracles,
		vaults:   vaults,
		walletz:  wallet.FilterTrace(walletz, time.Minute),
		notifier: notifier,
		system:   system,
	}
}

type Keeper struct {
	cats     core.CollateralStore
	oracles  core.OracleStore
	vaults   core.VaultStore
	walletz  core.WalletService
	notifier core.Notifier
	system   *core.System
}

func (w *Keeper) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "keeper")
	ctx = logger.WithContext(ctx, log)

	g := errgroup.Group{}

	g.Go(func() error {
		return w.foldAll(ctx)
	})

	g.Go(func() error {
		return w.scan(ctx)
	})

	return g.Wait()
}
