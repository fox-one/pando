package keeper

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func (w *Keeper) dealFlips(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-time.After(time.Second):
			_ = w.scanFinishFlips(ctx, t)
		}
	}
}

func (w *Keeper) scanFinishFlips(ctx context.Context, t time.Time) error {
	var (
		from  int64 = 0
		limit       = 500

		g   errgroup.Group
		sem = semaphore.NewWeighted(5)
	)

	for {
		select {
		case <-ctx.Done():
			return g.Wait()
		default:
		}

		flips, err := w.flips.List(ctx, from, limit)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("flips.List")
			return g.Wait()
		}

		for idx := range flips {
			flip := flips[idx]
			from = flip.ID

			if flip.Action == core.ActionFlipDeal {
				continue
			}

			if end := flip.EndFinished(t) || flip.TicFinished(t); !end {
				continue
			}

			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return g.Wait()
				}

				defer sem.Release(1)

				trace := uuid.Modify(flip.TraceID, "deal"+t.Truncate(time.Minute).String())
				return w.handleTransfer(ctx, trace, core.ActionFlipDeal, types.UUID(flip.TraceID))
			})
		}

		if len(flips) < limit {
			break
		}
	}

	return g.Wait()
}
