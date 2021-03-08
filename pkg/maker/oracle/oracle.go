package oracle

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "Oracle/"+msg)
}

func updatePrices(r *maker.Request, collaterals core.CollateralStore, oracles core.OracleStore) error {
	ctx := r.Context()

	cats, err := collaterals.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("collaterals.List")
		return err
	}

	prices, err := oracles.ListCurrent(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("oracles.ListCurrent")
		return err
	}

	for _, c := range cats {
		if c.Version >= r.Version || c.Live == 0 {
			continue
		}

		if gem, dai := prices.Get(c.Gem), prices.Get(c.Dai); gem.IsPositive() && dai.IsPositive() {
			if price := gem.Div(dai).Truncate(12); !c.Price.Equal(price) {
				c.Price = price
				if err := collaterals.Update(ctx, c, r.Version); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
					return err
				}
			}
		}
	}

	return nil
}
