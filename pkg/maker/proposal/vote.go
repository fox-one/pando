package proposal

import (
	"context"
	"database/sql"
	"encoding/base64"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleVote(
	proposals core.ProposalStore,
	parliaments core.Parliament,
	actions map[core.Action]maker.HandlerFunc,
	system *core.System,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(system.IsMember(r.UserID), "not-mtg-member"); err != nil {
			return err
		}

		p, err := From(ctx, proposals, r)
		if err != nil {
			return err
		}

		if voted := govalidator.IsIn(r.UserID, p.Votes...); !voted {
			p.Votes = append(p.Votes, r.UserID)

			if err := parliaments.Approved(ctx, p); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("parliaments.Approved")
				return err
			}

			if !p.PassedAt.Valid && len(p.Votes) >= int(system.Threshold) {
				p.PassedAt = sql.NullTime{
					Time:  r.Now(),
					Valid: true,
				}

				if err := parliaments.Passed(ctx, p); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("parliaments.Passed")
					return err
				}

				// execute
				if h, ok := actions[p.Action]; ok {
					r.UTXO.AssetID = p.AssetID
					r.UTXO.Amount = p.Amount
					r.Gov = true
					r.UserID = p.Creator
					r.Action = p.Action
					r.Body, _ = base64.StdEncoding.DecodeString(p.Data)

					if err := h(ctx, r); err != nil {
						return err
					}
				}
			}

			if err := proposals.Update(ctx, p, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("proposals.Update")
				return err
			}
		}

		return nil
	}
}
