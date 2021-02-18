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
		case <-time.After(time.Second):
			_ = w.run(ctx)
		}
	}
}

func (w *Sync) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	assets, err := w.assets.List(ctx)
	if err != nil {
		log.WithError(err).Error("assets.List")
		return err
	}

	for _, asset := range assets {
		r, err := w.assetz.Find(ctx, asset.ID)
		if err != nil {
			log.WithError(err).Errorf("assetz.Find(%s)", asset.Symbol)
			continue
		}

		asset.Price = r.Price
		asset.Logo = r.Logo
		if err := w.assets.Update(ctx, asset); err != nil {
			log.WithError(err).Errorf("assets.Update(%s)", asset.ID)
			return err
		}
	}

	return nil
}
