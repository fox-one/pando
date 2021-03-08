package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleKick(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Sender != "", "anonymous"); err != nil {
			return err
		}

		v, err := vat.From(r, vaults)
		if err != nil {
			return err
		}

		cid, _ := uuid.FromString(v.CollateralID)
		c, err := cat.From(r.WithBody(cid), collaterals)
		if err != nil {
			return err
		}

		if err := require(c.Live > 0, "not-live"); err != nil {
			if c.Dai == r.AssetID {
				err = maker.WithFlag(err, maker.FlagRefund)
			}

			return err
		}

		flip, err := flips.Find(ctx, r.TraceID)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("flips.Find")
			return err
		}

		if flip.ID == 0 {
			flip, err = Kick(r, c, v)
			if err != nil {
				return err
			}

			if err := flips.CreateEvent(ctx, &core.FlipEvent{
				CreatedAt: r.Now,
				FlipID:    flip.TraceID,
				Version:   r.Version,
				Action:    r.Action,
				Bid:       flip.Bid,
				Lot:       flip.Lot,
			}); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.CreateEvent")
				return err
			}

			if err := flips.Create(ctx, flip); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.Create")
				return err
			}
		}

		if v.Version < r.Version {
			v.Art = v.Art.Sub(flip.Art)
			v.Ink = v.Ink.Sub(flip.Lot)

			if err := vaults.Update(ctx, v, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		if c.Version < r.Version {
			c.Art = c.Art.Sub(flip.Art)
			c.Debt = c.Debt.Sub(flip.Bid)
			c.Ink = c.Ink.Sub(flip.Lot)

			if err := collaterals.Update(ctx, c, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}

func Kick(r *maker.Request, c *core.Collateral, v *core.Vault) (*core.Flip, error) {
	if err := require(v.Ink.Mul(c.Price).LessThan(v.Art.Mul(c.Rate).Mul(c.Mat)), "not-unsafe"); err != nil {
		return nil, err
	}

	dart := decimal.Min(
		c.Dunk.Div(c.Rate).Div(c.Chop),
		v.Art,
	)

	dink := v.Ink
	if dart.LessThan(v.Art) {
		dink = dart.Div(v.Art).Mul(dink).Truncate(8)
	}

	if err := require(dart.IsPositive() && dink.IsPositive(), "null-auction"); err != nil {
		return nil, err
	}

	tab := dart.Mul(c.Rate).Mul(c.Chop).Truncate(8)

	flip := &core.Flip{
		CreatedAt:    r.Now,
		Version:      r.Version,
		TraceID:      r.TraceID,
		CollateralID: v.CollateralID,
		VaultID:      v.TraceID,
		Action:       r.Action,
		Tic:          0,
		End:          r.Now.Unix() + c.Tau,
		Bid:          decimal.Zero,
		Lot:          dink,
		Tab:          tab,
		Art:          dart,
		Guy:          r.Sender,
	}

	if r.AssetID == c.Dai {
		flip.Bid = r.Amount
		flip.Tic = r.Now.Unix() + c.TTL
	}

	return flip, nil
}
