package flip

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/shopspring/decimal"
)

type KickData struct {
	Tab decimal.Decimal `json:"tab,omitempty"`
	Lot decimal.Decimal `json:"lot,omitempty"`
	Bid decimal.Decimal `json:"bid,omitempty"`
}

func HandleKick(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	transactions core.TransactionStore,
	properties property.Store,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.BindUser() == nil && r.BindFollow() == nil, "bad-data"); err != nil {
			return err
		}

		v, err := vat.From(ctx, vaults, r)
		if err != nil {
			return err
		}

		cid, _ := uuid.FromString(v.CollateralID)
		c, err := cat.From(ctx, collaterals, r.WithBody(cid))
		if err != nil {
			return err
		}

		t, err := transactions.Find(ctx, r.TraceID())
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("transactions.Find")
			return err
		}

		if t.ID == 0 {
			opt, err := ReadOptions(ctx, properties)
			if err != nil {
				return err
			}

			t := r.Tx()
			t.TargetID = v.TraceID

			if f, err := Kick(r, c, v, opt); err == nil {
				if err := flips.Create(ctx, f); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("flips.Create")
					return err
				}

				t.Write(core.TxStatusSuccess, KickData{
					Tab: f.Tab,
					Lot: f.Lot,
					Bid: f.Bid,
				})
			} else {
				t.Write(core.TxStatusFailed, err)
			}

			if err := transactions.Create(ctx, t); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("transactions.Create")
				return err
			}
		}

		if err := require(t.Status == core.TxStatusSuccess, "tx-failed"); err != nil {
			return err
		}

		var data KickData
		_ = t.Data.Unmarshal(&data)

		dart := data.Tab.Div(c.Chop).Div(c.Rate)

		if v.Version < r.Version() {
			v.Art = v.Art.Sub(dart)
			v.Ink = v.Ink.Sub(data.Lot)

			if err := vaults.Update(ctx, v, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		if c.Version < r.Version() {
			c.Art = c.Art.Sub(dart)
			c.Debt = c.Debt.Sub(data.Bid)

			if err := collaterals.Update(ctx, c, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}

func Kick(r *maker.Request, c *core.Collateral, v *core.Vault, opt *Option) (*core.Flip, error) {
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

	assetID, bid := r.Payment()
	if err := require(assetID == c.Dai && bid.IsPositive() && bid.LessThanOrEqual(tab), "bid-not-match"); err != nil {
		return nil, err
	}

	return &core.Flip{
		CreatedAt: r.Now(),
		Version:   r.Version(),
		TraceID:   r.TraceID(),
		VaultID:   v.TraceID,
		Action:    r.Action,
		Tic:       r.Now().Add(opt.TTL),
		End:       r.Now().Add(opt.Tau),
		Bid:       bid,
		Lot:       dink,
		Tab:       tab,
		Guy:       r.UserID,
	}, nil
}
