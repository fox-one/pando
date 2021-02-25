package cat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleSupply(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		log := logger.FromContext(ctx)

		c, err := From(ctx, collaterals, r)
		if err != nil {
			return err
		}

		if assetID, amount := r.Payment(); assetID == c.Dai && c.Version < r.Version() {
			c.Line = c.Line.Add(amount)

			if err := collaterals.Update(ctx, c, r.Version()); err != nil {
				log.WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}
