package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string, flags ...int) error {
	return maker.Require(condition, "Cat/"+msg, flags...)
}

func From(r *maker.Request, collaterals core.CollateralStore) (*core.Collateral, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)

	cat, err := collaterals.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return nil, err
	}

	if err := require(cat.ID > 0, "not-init"); err != nil {
		return nil, err
	}

	return cat, nil
}

func List(r *maker.Request, collaterals core.CollateralStore) ([]*core.Collateral, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)

	var cats []*core.Collateral

	if id == uuid.Zero {
		var err error
		cats, err = collaterals.List(ctx)
		if err != nil {
			log.WithError(err).Errorln("collaterals.List")
			return nil, err
		}
	} else {
		cat, err := collaterals.Find(ctx, id.String())
		if err != nil {
			log.WithError(err).Errorln("collaterals.Find")
			return nil, err
		}

		if err := require(cat.ID > 0, "not-init"); err != nil {
			return nil, err
		}

		cats = append(cats, cat)
	}

	return cats, nil
}
