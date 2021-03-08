package proposal

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pkg/logger"
)

func HandleShout(
	proposals core.ProposalStore,
	parliaments core.Parliament,
	system *core.System,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(system.IsMember(r.Sender), "not-member"); err != nil {
			return err
		}

		p, err := From(r, proposals)
		if err != nil {
			return err
		}

		if err := parliaments.Created(ctx, p); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("parliaments.Created")
			return err
		}

		return nil
	}
}
