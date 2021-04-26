package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type DealData struct {
	Bid decimal.Decimal `json:"bid,omitempty"`
	Lot decimal.Decimal `json:"lot,omitempty"`
}

func HandleDeal(
	collaterals core.CollateralStore,
	flips core.FlipStore,
	wallets core.WalletStore,
	parliaments core.Parliament,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		f, err := From(r, flips)
		if err != nil {
			return err
		}

		cid, _ := uuid.FromString(f.CollateralID)
		c, err := cat.From(r.WithBody(cid), collaterals)
		if err != nil {
			return err
		}

		event, err := flips.FindEvent(ctx, f.TraceID, r.Version)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("flips.FindEvent")
			return err
		}

		if event.ID == 0 {
			if err := Deal(r, f); err != nil {
				return err
			}

			var transfers []*core.Transfer
			if f.Lot.IsPositive() {
				memo := core.TransferAction{
					ID:     f.TraceID,
					Source: r.Action.String(),
				}.Encode()

				transfers = append(transfers, &core.Transfer{
					TraceID:   uuid.Modify(r.TraceID, memo),
					AssetID:   c.Gem,
					Amount:    f.Lot,
					Memo:      memo,
					Threshold: 1,
					Opponents: []string{f.Guy},
				})
			}

			if err := wallets.CreateTransfers(ctx, transfers); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}

			event = &core.FlipEvent{
				CreatedAt: r.Now,
				FlipID:    f.TraceID,
				Version:   r.Version,
				Action:    r.Action,
				Bid:       f.Bid,
				Lot:       f.Lot,
			}

			if err := flips.CreateEvent(ctx, event); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.CreateEvent")
				return err
			}
		}

		if c.Version < r.Version {
			c.Litter = c.Litter.Sub(f.Tab)

			if err := collaterals.Update(ctx, c, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		if f.Version < r.Version {
			f.Action = r.Action

			if err := flips.Update(ctx, f, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.Update")
				return err
			}
		}

		if err := parliaments.FlipDeal(ctx, f); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("parliaments.FlipDeal")
			return err
		}

		return nil
	}
}

func Deal(r *maker.Request, f *core.Flip) error {
	if err := require(f.TicFinished(r.Now) || f.EndFinished(r.Now), "not-finished"); err != nil {
		return err
	}

	if err := require(f.Action != core.ActionFlipDeal, "already-dealed"); err != nil {
		return err
	}

	return nil
}
