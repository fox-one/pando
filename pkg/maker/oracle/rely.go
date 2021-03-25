package oracle

import (
	"encoding/base64"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func HandleRely(oracles core.OracleStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		var (
			id        uuid.UUID
			publicKey types.RawMessage
		)

		if err := require(r.Scan(&id, &publicKey) == nil, "bad-data"); err != nil {
			return err
		}

		if err := oracles.Rely(ctx, id.String(), base64.StdEncoding.EncodeToString(publicKey)); err != nil {
			log.WithError(err).Errorln("oracles.Rely")
			return err
		}

		return nil
	}
}
