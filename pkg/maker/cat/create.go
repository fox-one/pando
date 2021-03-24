package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleCreate(
	collaterals core.CollateralStore,
	oracles core.OracleStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		var (
			gem, dai uuid.UUID
			name     string
		)

		if err := require(r.Scan(&gem, &dai, &name) == nil, "bad-data"); err != nil {
			return err
		}

		if err := require(gem != uuid.Zero && dai != uuid.Zero, "invalid-asset"); err != nil {
			return err
		}

		if err := require(gem.String() != dai.String(), "same-asset"); err != nil {
			return err
		}

		cat := &core.Collateral{
			CreatedAt: r.Now,
			TraceID:   r.TraceID,
			Version:   r.Version,
			Name:      name,
			Gem:       gem.String(),
			Dai:       dai.String(),
			Art:       decimal.Zero,
			Rate:      number.Decimal("1"),
			Rho:       r.Now,
			Dust:      number.Decimal("100"),
			Mat:       number.Decimal("1.5"),
			Duty:      number.Decimal("1.05"),
			Chop:      number.Decimal("1.13"),
			Dunk:      number.Decimal("5000"),
			Box:       number.Decimal("500000"),
			Beg:       number.Decimal("0.03"),
			TTL:       15 * 60,     // 15m
			Tau:       3 * 60 * 60, // 3h
		}

		prices, err := oracles.ListCurrent(ctx)
		if err != nil {
			log.WithError(err).Errorln("oracles.ListCurrent")
			return err
		}

		if gp, dp := prices.Get(cat.Gem), prices.Get(cat.Dai); gp.IsPositive() && dp.IsPositive() {
			cat.Price = gp.Div(dp).Truncate(12)
		}

		if err := collaterals.Create(ctx, cat); err != nil {
			log.WithError(err).Errorln("collaterals.Create")
			return err
		}

		return nil
	}
}
