package proposal

import (
	"encoding/base64"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleMake(
	proposals core.ProposalStore,
	walletz core.WalletService,
	parliaments core.Parliament,
	system *core.System,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		var action core.Action
		if err := require(r.Scan(&action) == nil, "bad-data"); err != nil {
			return err
		}

		p := &core.Proposal{
			CreatedAt: r.Now,
			Version:   r.Version,
			TraceID:   r.TraceID,
			Creator:   r.Sender,
			AssetID:   r.AssetID,
			Amount:    r.Amount,
			Action:    action,
			Data:      base64.StdEncoding.EncodeToString(r.Body),
		}

		ctx := r.Context()
		if err := proposals.Create(ctx, p); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("proposals.Create")
			return err
		}

		if system.IsMember(p.Creator) {
			if err := parliaments.ProposalCreated(ctx, p); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("parliaments.ProposalCreated")
				return err
			}
		} else if system.IsStaff(p.Creator) {
			if err := handleProposal(r, walletz, system, core.ActionProposalShout, p); err != nil {
				return err
			}
		}

		return nil
	}
}
