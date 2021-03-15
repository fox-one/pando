package cat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
	"github.com/shopspring/decimal"
)

func HandleCreate(
	collaterals core.CollateralStore,
	oracles core.OracleStore,
	assets core.AssetStore,
	assetz core.AssetService,
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

		if err := require(gem.String() != dai.String(), "same-asset"); err != nil {
			return err
		}

		if _, err := handleAsset(ctx, assets, assetz, gem.String()); err != nil {
			return err
		}

		if _, err := handleAsset(ctx, assets, assetz, dai.String()); err != nil {
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

func handleAsset(ctx context.Context, assets core.AssetStore, assetz core.AssetService, id string) (*core.Asset, error) {
	log := logger.FromContext(ctx)

	asset, err := assets.Find(ctx, id)
	if err != nil {
		if !store.IsErrNotFound(err) {
			log.WithError(err).Errorln("assets.Find")
			return nil, err
		}

		asset, err = assetz.Find(ctx, id)
		if err != nil {
			log.WithError(err).Errorln("assets.Find")
			return nil, err
		}

		if err := require(asset.Symbol != "", "asset-not-exist"); err != nil {
			return nil, err
		}

		if err := assets.Create(ctx, asset); err != nil {
			log.WithError(err).Errorln("assets.Create")
			return nil, err
		}
	}

	if asset.ID != asset.ChainID {
		if _, err := handleAsset(ctx, assets, assetz, asset.ChainID); err != nil {
			return nil, err
		}
	}

	return asset, nil
}
