package keeper

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

func (w *Keeper) scanVaults(ctx context.Context) error {
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

			oracles, err := w.oracles.List(ctx)
			if err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("oracles.List")
				return err
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

				nextPrice, _ := nextPrice(
					lookupOracle(oracles, cat.Gem),
					lookupOracle(oracles, cat.Dai),
				)

				g.Go(func() error {
					defer running.Delete(cat.TraceID)
					return w.scanUnsafeVaults(ctx, cat, nextPrice)
				})
			}
		}
	}
}

func (w *Keeper) scanUnsafeVaults(ctx context.Context, cat *core.Collateral, nextPrice decimal.Decimal) error {
	// v.Ink * c.Price >= v.Art * c.Rate * c.Mat
	// v.Rate = v.Art / v.Ink <= c.Price / c.Rate / c.Mat

	rate := cat.Price.Div(cat.Rate).Div(cat.Mat)

	nextRate := rate
	if nextPrice.IsPositive() {
		nextRate = nextPrice.Div(cat.Rate).Div(cat.Mat)
	}

	scanRate := decimal.Min(rate.Mul(number.Decimal("0.8")))

	var (
		from  int64
		limit = 100
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		vats, err := w.vaults.List(ctx, core.ListVaultRequest{
			CollateralID: cat.TraceID,
			Rate:         scanRate,
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

			switch {
			case vat.Rate.GreaterThan(rate):
				trace := uuid.Modify(vat.TraceID, fmt.Sprintf("%s-%d", rate, vat.Version))
				_ = w.handleTransfer(ctx, trace, core.ActionFlipKick, types.UUID(vat.TraceID))
			case vat.Rate.GreaterThan(nextRate):
				if err := w.notifier.VaultLiquidatedSoon(ctx, cat, vat); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("notifier.VaultLiquidatedSoon")
					return err
				}
			default:
				if err := w.notifier.VaultUnsafe(ctx, cat, vat); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("notifier.VaultUnsafe")
					return err
				}
			}
		}

		if len(vats) < limit {
			break
		}
	}

	return nil
}

func nextPrice(gem, dai *core.Oracle) (next decimal.Decimal, at time.Time) {
	if gem == nil || dai == nil {
		return
	}

	if n1, n2 := gem.NextPeekAt(), dai.NextPeekAt(); n1.Before(n2) {
		next = gem.Next.Div(dai.Current)
		at = n1
	} else if n1.After(n2) {
		next = gem.Current.Div(dai.Next)
		at = n2
	} else {
		next = gem.Next.Div(dai.Next)
		at = n2
	}

	return
}

func lookupOracle(oracles []*core.Oracle, assetID string) *core.Oracle {
	for _, oracle := range oracles {
		if oracle.AssetID == assetID {
			return oracle
		}
	}

	return nil
}
