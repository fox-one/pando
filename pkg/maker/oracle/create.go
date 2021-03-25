package oracle

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleCreate(oracles core.OracleStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		var (
			id        uuid.UUID
			price     decimal.Decimal
			hop       int64
			threshold int64
			ts        int64
		)

		if err := require(r.Scan(&id, &price, &hop, &threshold, &ts) == nil, "bad-data"); err != nil {
			return err
		}

		if err := require(price.IsPositive() && hop > 0 && ts < r.Now.Unix(), "bad-data"); err != nil {
			return err
		}

		oracle := &core.Oracle{
			CreatedAt: r.Now,
			AssetID:   id.String(),
			Version:   r.Version,
			Current:   price,
			Next:      price,
			PeekAt:    time.Unix(ts, 0).Truncate(time.Duration(hop) * time.Second),
			Hop:       hop,
			Threshold: threshold,
			Governors: r.Governors,
		}

		if err := oracles.Create(ctx, oracle); err != nil {
			log.WithError(err).Errorln("oracles.Create")
			return err
		}

		return nil
	}
}
