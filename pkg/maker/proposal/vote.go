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
			if handled := p.PassedAt.Valid || govalidator.IsIn(r.Sender, p.Votes...); !handled {
				p.Votes = append(p.Votes, r.Sender)

				if err := parliaments.ProposalApproved(ctx, p); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("parliaments.proposalApproved")
					return err
				}

				if len(p.Votes) >= int(system.Threshold) {
					p.PassedAt = sql.NullTime{
						Time:  r.Now,
						Valid: true,
					}

					if err := parliaments.ProposalPassed(ctx, p); err != nil {
						logger.FromContext(ctx).WithError(err).Errorln("parliaments.ProposalPassed")
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
