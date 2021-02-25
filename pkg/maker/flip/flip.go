package flip

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "Flip/%s", msg)
}

func From(ctx context.Context, flips core.FlipStore, r *maker.Request) (*core.Flip, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	log := logger.FromContext(ctx)

	f, err := flips.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return nil, err
	}

	if err := require(f.ID > 0, "not-init"); err != nil {
		return nil, err
	}

	return f, nil
}
