package pricesync

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
)

func New(
	assets core.AssetStore,
	assetz core.AssetService,
) *Sync {
	return &Sync{
		assets: assets,
		assetz: assetz,
	}
}

type Sync struct {
	assets core.AssetStore
	assetz core.AssetService
}

func (w *Sync) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "price syncer")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			_ = w.run(ctx)
		}
	}
}

func (w *Sync) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	assets, err := w.assetz.List(ctx)
	if err != nil {
		log.WithError(err).Error("assetz.List")
		return err
	}

	for _, asset := range assets {
		if err := w.assets.Save(ctx, asset); err != nil {
			log.WithError(err).Error("assets.Save")
			return err
		}
	}

	return nil
}
