package oracle

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleFeed(
	collaterals core.CollateralStore,
	oracles core.OracleStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		oracle, err := From(r, oracles)
		if err != nil {
			return err
		}

		var price decimal.Decimal
		if err := require(r.Scan(&price) == nil, "bad-data"); err != nil {
			return err
		}

		price = price.Truncate(12)
		if err := require(price.IsPositive(), "bad-data"); err != nil {
			return err
		}

		if oracle.ID == 0 {
			oracle.Hop = 60 * 60 // an hour
		}

		oracle.Current = price
		oracle.Next = price
		oracle.PeekAt = r.Now.Truncate(time.Duration(oracle.Hop) * time.Second)

		if err := oracles.Save(ctx, oracle, r.Version); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("oracles.Save")
			return err
		}

		return updatePrices(r, collaterals, oracles)
	}
}
