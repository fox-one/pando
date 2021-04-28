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

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		oracle, err := From(r, oracles)
		if err != nil {
			return err
		}

		if err := require(oracle.Threshold > 0 && int(oracle.Threshold) <= len(r.Governors), "threshold-not-reach"); err != nil {
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

		if err := require(ts >= oracle.PeekAt.Unix()+oracle.Hop, "not-passed"); err != nil {
			return err
		}

		oracle.Current = oracle.Next
		oracle.Next = price
		oracle.PeekAt = time.Unix(ts, 0).Truncate(time.Duration(oracle.Hop) * time.Second)
		oracle.Governors = r.Governors

		if err := oracles.Update(ctx, oracle, r.Version); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("oracles.Save")
			return err
		}

		return updatePrices(r, collaterals, oracles)
	}
}
