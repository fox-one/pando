package oracle

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandlePoke(
	collaterals core.CollateralStore,
	oracles core.OracleStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Gov || r.Oracle, "not-authorized"); err != nil {
			return err
		}

		oracle, err := From(r, oracles)
		if err != nil {
			return err
		}

		var (
			price decimal.Decimal
			ts    int64
		)

		if err := require(r.Scan(&price, &ts) == nil, "bad-data"); err != nil {
			return err
		}

		price = price.Truncate(12)
		if err := require(price.IsPositive() && ts < r.Now.Unix(), "bad-data"); err != nil {
			return err
		}

		if oracle.ID == 0 {
			oracle.Hop = 30 * 60 // 30m
			oracle.Next = price
		}

		if err := require(ts > oracle.PeekAt.Unix(), "out-of-date"); err != nil {
			return err
		}

		passed := ts > oracle.PeekAt.Unix()+oracle.Hop

		// Gov 可以随意更改 Next Price
		if err := require(r.Gov || passed, "not-passed"); err != nil {
			return err
		}

		oracle.Next = price
		if passed {
			oracle.Current = oracle.Next
			oracle.PeekAt = time.Unix(ts, 0).Truncate(time.Duration(oracle.Hop) * time.Second)
		}

		if err := oracles.Save(ctx, oracle, r.Version); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("oracles.Save")
			return err
		}

		return updatePrices(r, collaterals, oracles)
	}
}
