package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string, flags ...int) error {
	return maker.Require(condition, "Flip/"+msg, flags...)
}

func From(r *maker.Request, flips core.FlipStore) (*core.Flip, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	ctx := r.Context()
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
