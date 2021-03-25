package oracle

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/spf13/cast"
)

func HandleEdit(oracles core.OracleStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		list, err := List(r, oracles)
		if err != nil {
			return err
		}

		for {
			var key, value string
			if err := require(r.Scan(&key, &value) == nil, "bad-data"); err != nil {
				break
			}

			for _, oracle := range list {
				switch key {
				case "hop":
					if ts := cast.ToInt64(value); ts > 0 {
						oracle.Hop = ts
					}
				case "next":
					if next := number.Decimal(value).Truncate(12); next.IsPositive() {
						oracle.Next = next
					}
				case "threshold":
					if threshold := cast.ToInt64(value); threshold >= 0 {
						oracle.Threshold = threshold
					}
				}
			}
		}

		for _, oracle := range list {
			if err := oracles.Update(ctx, oracle, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("oracles.Update")
				return err
			}
		}

		return nil
	}
}
