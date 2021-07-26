package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleSupply(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		c, err := From(r, collaterals)
		if err != nil {
			return err
		}

		if err := require(r.AssetID == c.Dai, "dai-not-match"); err != nil {
			return err
		}

		if err := require(r.SysVersion >= 4 || c.Live > 0, "not-live"); err != nil {
			return maker.WithFlag(err, maker.FlagRefund)
		}

		if c.Version < r.Version {
			c.Supply = c.Supply.Add(r.Amount)

			if err := collaterals.Update(ctx, c, r.Version); err != nil {
				log.WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}
