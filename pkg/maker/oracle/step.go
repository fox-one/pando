package oracle

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleStep(oracles core.OracleStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		list, err := List(r, oracles)
		if err != nil {
			return err
		}

		var ts int64
		if err := require(r.Scan(&ts) == nil, "bad-data"); err != nil {
			return err
		}

		if err := require(ts > 0, "ts-is-zero"); err != nil {
			return err
		}

		for _, oracle := range list {
			oracle.Hop = ts
			if err := oracles.Save(ctx, oracle, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("oracles.Save")
				return err
			}
		}

		return nil
	}
}
