package cat

import (
	"context"
	"strconv"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
)

type ModifyData struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func HandleEdit(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		log := logger.FromContext(ctx)

		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		c, err := From(ctx, collaterals, r)
		if err != nil {
			return err
		}

		for {
			var key, value string
			if err := require(r.Scan(&key, &value) == nil, ""); err != nil {
				break
			}

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
			case "live":
				if live, _ := strconv.ParseBool(value); live {
					c.Live = 1
				} else {
					c.Live = 0
				}
			}
		}

		if err := collaterals.Update(ctx, c, r.Version()); err != nil {
			log.WithError(err).Errorln("collaterals.Update")
			return err
		}

		return nil
	}
}
