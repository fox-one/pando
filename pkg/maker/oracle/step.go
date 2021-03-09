package oracle

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func HandleStep(oracles core.OracleStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		var (
			id uuid.UUID
			ts int64
		)

		if err := require(r.Scan(&id, &ts) == nil, "bad-data"); err != nil {
			return err
		}

		if err := require(ts > 0, "ts-is-zero"); err != nil {
			return err
		}

		oracle, err := oracles.Find(ctx, id.String())
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("oracles.Find")
			return err
		}

		if err := require(oracle.ID > 0, "not-init"); err != nil {
			return err
		}

		oracle.Hop = ts
		if err := oracles.Save(ctx, oracle, r.Version); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("oracles.Save")
			return err
		}

		return nil
	}
}
