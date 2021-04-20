package keeper

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func (w *Keeper) scan(ctx context.Context) error {
	running := cache.New(cache.NoExpiration, cache.NoExpiration)
	var g errgroup.Group

	for {
		select {
		case <-ctx.Done():
			return g.Wait()
		case <-time.After(time.Second):
			cats, err := w.cats.List(ctx)
			if err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("cats.List")
				break
			}

			for idx := range cats {
				cat := cats[idx]

				if cat.Live == 0 {
					continue
				}

				if _, ok := running.Get(cat.TraceID); ok {
					continue
				}

				running.SetDefault(cat.TraceID, nil)

				g.Go(func() error {
					defer running.Delete(cat.TraceID)
					return w.scanUnsafeVaults(ctx, cat)
				})
			}
		}
	}
}

func (w *Keeper) scanUnsafeVaults(ctx context.Context, cat *core.Collateral) error {
	// v.Ink * c.Price >= v.Art * c.Rate * c.Mat
	// v.Rate = v.Art / v.Ink <= c.Price / c.Rate / c.Mat

	rate := cat.Price.Div(cat.Rate).Div(cat.Mat)

	var (
		g     errgroup.Group
		sem   = semaphore.NewWeighted(5)
		from  int64
		limit = 100
	)

	for {
		select {
		case <-ctx.Done():
			return g.Wait()
		default:
		}

		vats, err := w.vaults.List(ctx, core.ListVaultRequest{
			CollateralID: cat.TraceID,
			Rate:         rate,
			Desc:         true,
			FromID:       from,
			Limit:        limit,
		})

		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("vaults.List")
			break
		}

		for idx := range vats {
			vat := vats[idx]
			from = vat.ID

			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}

				defer sem.Release(1)

				trace := uuid.Modify(vat.TraceID, fmt.Sprintf("%s-%d", rate, vat.Version))
				return w.handleTransfer(ctx, trace, core.ActionFlipKick, types.UUID(vat.TraceID))
			})
		}

		if len(vats) < limit {
			break
		}
	}

	return g.Wait()
}

// func nextPrice(gem, dai *core.Oracle) (next decimal.Decimal, at time.Time) {
// 	if gem == nil || dai == nil {
// 		return
// 	}
//
// 	if n1, n2 := gem.NextPeekAt(), dai.NextPeekAt(); n1.Before(n2) {
// 		next = gem.Next.Div(dai.Current)
// 		at = n1
// 	} else if n1.After(n2) {
// 		next = gem.Current.Div(dai.Next)
// 		at = n2
// 	} else {
// 		next = gem.Next.Div(dai.Next)
// 		at = n2
// 	}
//
// 	return
// }
