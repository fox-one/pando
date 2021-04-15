package keeper

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
)

type Keeper struct {
	cats core.CollateralStore
}

func (w *Keeper) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "keeper")
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

func (w *Keeper) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	cats, err := w.cats.List(ctx)
	if err != nil {
		log.WithError(err).Errorln("cats.List")
		return err
	}

	// remove cat not live
	var idx int
	for _, cat := range cats {
		if cat.Live > 0 {
			cats[idx] = cat
			idx++
		}
	}

	cats = cats[:idx]
}
