package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

// HandleFold modify the debt multiplier, creating / destroying corresponding debt
func HandleFold(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		cats, err := List(r, collaterals)
		if err != nil {
			return err
		}

		for _, c := range cats {
			if err := require(c.Live > 0, "not-live"); err != nil {
				continue
			}

			n := r.Now.Unix() - c.Rho.Unix()
			if n > 0 {
				const year int64 = 60 * 60 * 24 * 365
				q := decimal.NewFromInt(n).Div(decimal.NewFromInt(year))
				f := number.Pow(c.Duty, q)
				if rate := c.Rate.Mul(f).Truncate(16); rate.GreaterThan(c.Rate) {
					c.Rate = rate
					c.Rho = r.Now

					if err := collaterals.Update(ctx, c, r.Version); err != nil {
						log.WithError(err).Errorln("collaterals.Update")
						return err
					}
				}
			}
		}

		return nil
	}
}
