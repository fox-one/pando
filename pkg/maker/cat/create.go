package cat

import (
	"context"

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
	assets core.AssetStore,
	assetz core.AssetService,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
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

		if err := require(gem.String() != dai.String(), "same-asset"); err != nil {
			return err
		}

		gemAsset, err := assetz.Find(ctx, gem.String())
		if err != nil {
			log.WithError(err).Errorln("assetz.Find")
			return err
		}

		if err := require(gemAsset.Symbol != "", "nil-asset"); err != nil {
			return err
		}

		if err := assets.Create(ctx, gemAsset); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("assets.Create")
			return err
		}

		daiAsset, err := assetz.Find(ctx, dai.String())
		if err != nil {
			return err
		}

		if err := require(daiAsset.Symbol != "", "nil-asset"); err != nil {
			return err
		}

		if err := assets.Create(ctx, daiAsset); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("assets.Create")
			return err
		}

		cat := &core.Collateral{
			CreatedAt: r.Now(),
			TraceID:   r.TraceID(),
			Version:   r.Version(),
			Name:      name,
			Gem:       gem.String(),
			Dai:       dai.String(),
			Art:       decimal.Zero,
			Rate:      number.Decimal("1"),
			Rho:       r.Now(),
			Dust:      number.Decimal("100"),
			Mat:       number.Decimal("1.5"),
			Duty:      number.Decimal("1.05"),
			Chop:      number.Decimal("1.13"),
			Dunk:      number.Decimal("5000"),
		}

		if assetID, amount := r.Payment(); assetID == cat.Dai {
			cat.Line = amount
		}

		goc, err := oracles.Find(ctx, gem.String(), r.Now())
		if err != nil {
			log.WithError(err).Errorln("oracles.Find")
			return err
		}

		doc, err := oracles.Find(ctx, dai.String(), r.Now())
		if err != nil {
			log.WithError(err).Errorln("oracles.Find")
			return err
		}

		if goc.Price.IsPositive() && doc.Price.IsPositive() {
			cat.Price = goc.Price.Div(doc.Price).Truncate(12)
		}

		if err := collaterals.Create(ctx, cat); err != nil {
			log.WithError(err).Errorln("collaterals.Create")
			return err
		}

		return nil
	}
}
