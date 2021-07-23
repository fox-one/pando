package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleFill(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		c, err := From(r, collaterals)
		if err != nil {
			return err
		}

		if err := require(c.Dai == r.AssetID, "dai-not-match"); err != nil {
			return err
		}

		if c.Version >= r.Version {
			return nil
		}

		c.Debt = c.Debt.Sub(r.Amount)
		if err := collaterals.Update(ctx, c, r.Version); err != nil {
			log.WithError(err).Errorln("collaterals.Update")
			return err
		}

		return nil
	}
}
