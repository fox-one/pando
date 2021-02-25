package oracle

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleFeed(
	collaterals core.CollateralStore,
	oracles core.OracleStore,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		var (
			id    uuid.UUID
			price decimal.Decimal
			ts    int64
		)

		if err := require(r.Scan(&id, &price, &ts) == nil, "bad-data"); err != nil {
			return err
		}

		if err := require(price.IsPositive() && time.Unix(ts, 0).Before(r.Now()), "validate-failed"); err != nil {
			return err
		}

		if err := oracles.Create(ctx, &core.Oracle{
			CreatedAt: r.Now(),
			PeekAt:    time.Unix(ts, 0),
			TraceID:   r.TraceID(),
			AssetID:   id.String(),
			Price:     price,
		}); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("oracles.Create")
			return err
		}

		cats, err := collaterals.List(ctx)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("collaterals.List")
			return err
		}

		prices := number.Values{}

		for _, c := range cats {
			if c.Version >= r.Version() {
				continue
			}

			if c.Gem != id.String() && c.Dai != id.String() {
				continue
			}

			if _, ok := prices[c.Gem]; !ok {
				o, err := oracles.Find(ctx, c.Gem, r.Now())
				if err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("oracles.Find")
					return err
				}

				prices.Set(o.AssetID, o.Price)
			}

			if _, ok := prices[c.Dai]; !ok {
				o, err := oracles.Find(ctx, c.Dai, r.Now())
				if err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("oracles.Find")
					return err
				}

				prices.Set(o.AssetID, o.Price)
			}

			if gem, dai := prices.Get(c.Gem), prices.Get(c.Dai); gem.IsPositive() && dai.IsPositive() {
				c.Price = gem.Div(dai).Truncate(12)
				if err := collaterals.Update(ctx, c, r.Version()); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
					return err
				}
			}
		}

		return nil
	}
}
