package proposal

import (
	"context"
	"encoding/base64"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func HandleInit(
	proposals core.ProposalStore,
	parliaments core.Parliament,
	system *core.System,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(system.IsMember(r.UserID), "not-mtg-member"); err != nil {
			return err
		}

		var (
			trace  uuid.UUID
			action core.Action
		)
		if err := require(r.Scan(&trace, &action) == nil, "bad-data"); err != nil {
			return err
		}

		assetID, amount := r.Payment()
		p := &core.Proposal{
			CreatedAt: r.Now(),
			Version:   r.Version(),
			TraceID:   trace.String(),
			Creator:   r.UserID,
			AssetID:   assetID,
			Amount:    amount,
			Action:    action,
			Data:      base64.StdEncoding.EncodeToString(r.Body),
		}

		if err := proposals.Create(ctx, p); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("proposals.Create")
			return err
		}

		if err := parliaments.Created(ctx, p); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("parliaments.Created")
			return err
		}

		return nil
	}
}
