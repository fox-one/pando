package oracle

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "Oracle/"+msg)
}

func From(r *maker.Request, oracles core.OracleStore) (*core.Oracle, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)

	oracle, err := oracles.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("oracles.Find")
		return nil, err
	}

	return oracle, nil
}

func List(r *maker.Request, oracles core.OracleStore) ([]*core.Oracle, error) {
	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)

	var list []*core.Oracle
	if id == uuid.Zero {
		var err error
		list, err = oracles.List(ctx)
		if err != nil {
			log.WithError(err).Errorln("oracles.List")
			return nil, err
		}
	} else {
		oracle, err := oracles.Find(ctx, id.String())
		if err != nil {
			log.WithError(err).Errorln("oracles.Find")
			return nil, err
		}

		if err := require(oracle.ID > 0, "not-init"); err != nil {
			return nil, err
		}

		list = append(list, oracle)
	}

	return list, nil
}

func updatePrices(r *maker.Request, collaterals core.CollateralStore, oracles core.OracleStore) error {
	ctx := r.Context()

	cats, err := collaterals.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("collaterals.List")
		return err
	}

	prices, err := oracles.ListCurrent(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("oracles.ListCurrent")
		return err
	}

	for _, c := range cats {
		if c.Version >= r.Version {
			continue
		}

		if gem, dai := prices.Get(c.Gem), prices.Get(c.Dai); gem.IsPositive() && dai.IsPositive() {
			if price := gem.Div(dai).Truncate(12); !c.Price.Equal(price) {
				c.Price = price
				if err := collaterals.Update(ctx, c, r.Version); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
					return err
				}
			}
		}
	}

	return nil
}
