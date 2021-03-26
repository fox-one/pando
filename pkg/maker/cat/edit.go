package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/spf13/cast"
)

func HandleEdit(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		cats, err := List(r, collaterals)
		if err != nil {
			return err
		}

		for {
			var key, value string
			if err := require(r.Scan(&key, &value) == nil, ""); err != nil {
				break
			}

			for _, c := range cats {
				switch key {
				case "dust":
					c.Dust = number.Decimal(value)
				case "price":
					c.Price = number.Decimal(value)
				case "mat":
					c.Mat = number.Decimal(value)
				case "duty":
					c.Duty = number.Decimal(value)
				case "chop":
					c.Chop = number.Decimal(value)
				case "dunk":
					c.Dunk = number.Decimal(value)
				case "box":
					c.Box = number.Decimal(value)
				case "beg":
					c.Beg = number.Decimal(value)
				case "ttl":
					c.TTL = cast.ToInt64(value)
				case "tau":
					c.Tau = cast.ToInt64(value)
				case "line":
					if line := number.Decimal(value); line.LessThanOrEqual(c.Supply) {
						c.Line = line
					}
				case "live":
					if live := cast.ToBool(value); live {
						c.Live = 1
					} else {
						c.Live = 0
					}
				}
			}
		}

		for _, c := range cats {
			if err := collaterals.Update(ctx, c, r.Version); err != nil {
				log.WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}
