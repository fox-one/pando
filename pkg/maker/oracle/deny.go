package oracle

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func HandleDeny(oracles core.OracleStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		var (
			id uuid.UUID
		)

		if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
			return err
		}

		if err := oracles.Deny(ctx, id.String()); err != nil {
			log.WithError(err).Errorln("oracles.Deny")
			return err
		}

		return nil
	}
}
