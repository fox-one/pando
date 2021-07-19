package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

// HandleMove move the same dai from one to another
func HandleMove(collaterals core.CollateralStore) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		var (
			fromID, toID uuid.UUID
			amount       decimal.Decimal
		)

		if err := require(r.Scan(&fromID, &toID, &amount) == nil, "bad-data"); err != nil {
			return err
		}

		if err := require(fromID.String() != toID.String(), "internal"); err != nil {
			return err
		}

		amount = amount.Truncate(8)
		if err := require(amount.IsPositive(), "back-flow"); err != nil {
			return err
		}

		from, err := From(r.WithBody(fromID), collaterals)
		if err != nil {
			return err
		}

		to, err := From(r.WithBody(toID), collaterals)
		if err != nil {
			return err
		}

		if err := require(from.Dai == to.Dai, "same-dai"); err != nil {
			return err
		}

		if from.Version < r.Version {
			max := from.Supply.Sub(decimal.Max(from.Line, from.Debt))
			if err := require(amount.LessThanOrEqual(max), "insufficient-supply"); err != nil {
				return err
			}

			from.Supply = from.Supply.Sub(amount)
			if err := collaterals.Update(ctx, from, r.Version); err != nil {
				log.WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		if to.Version < r.Version {
			to.Supply = to.Supply.Add(amount)
			if err := collaterals.Update(ctx, to, r.Version); err != nil {
				log.WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}
