package proposal

import (
	"database/sql"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleVote(
	proposals core.ProposalStore,
	parliaments core.Parliament,
	walletz core.WalletService,
	system *core.System,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		p, err := From(r, proposals)
		if err != nil {
			return err
		}

		if system.IsMember(r.Sender) {
			if voted := govalidator.IsIn(r.Sender, p.Votes...); !voted {
				p.Votes = append(p.Votes, r.Sender)

				if err := parliaments.Approved(ctx, p); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("parliaments.Approved")
					return err
				}

				if !p.PassedAt.Valid && len(p.Votes) >= int(system.Threshold) {
					p.PassedAt = sql.NullTime{
						Time:  r.Now,
						Valid: true,
					}

					if err := parliaments.Passed(ctx, p); err != nil {
						logger.FromContext(ctx).WithError(err).Errorln("parliaments.Passed")
						return err
					}
				}

				if err := proposals.Update(ctx, p, r.Version); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("proposals.Update")
					return err
				}
			}

			if p.PassedAt.Valid && p.Version == r.Version {
				r.Next = r.WithProposal(p)
			}
		} else if system.IsStaff(r.Sender) {
			if err := handleProposal(r, walletz, system, core.ActionProposalVote, p); err != nil {
				return err
			}
		}

		return nil
	}
}
