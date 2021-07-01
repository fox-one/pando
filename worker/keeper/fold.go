package keeper

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func (w *Keeper) foldCats(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-time.After(time.Second):
			cats, err := w.cats.List(ctx)
			if err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("cats.List")
				break
			}

			if len(cats) == 0 {
				break
			}

			fold := false
			for _, cat := range cats {
				if t.Sub(cat.Rho) > time.Second*90 {
					fold = true
					break
				}
			}

			if fold {
				trace := uuid.MD5(t.Truncate(time.Minute).Format(time.RFC3339Nano))
				_ = w.handleTransfer(ctx, trace, core.ActionCatFold, uuid.Zero)
			}
		}
	}
}
