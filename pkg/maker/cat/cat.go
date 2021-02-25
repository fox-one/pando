package cat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "Cat/%s", msg)
}

func From(ctx context.Context, collaterals core.CollateralStore, r *maker.Request) (*core.Collateral, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	log := logger.FromContext(ctx)

	cat, err := collaterals.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return nil, err
	}

	if err := require(cat.ID > 0, "not-init"); err != nil {
		return nil, err
	}

	if err := require(r.Gov || cat.Live > 0, "not-live"); err != nil {
		return nil, err
	}

	return cat, nil
}
