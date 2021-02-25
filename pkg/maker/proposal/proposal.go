package proposal

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "proposal/%s", msg)
}

func From(ctx context.Context, proposals core.ProposalStore, r *maker.Request) (*core.Proposal, error) {
	log := logger.FromContext(ctx)

	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	p, err := proposals.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("proposals.Find")
		return nil, err
	}

	if err := require(p.ID > 0, "not init"); err != nil {
		return nil, err
	}

	return p, nil
}
