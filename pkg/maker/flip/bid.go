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

type BidData struct {
	Bid decimal.Decimal `json:"bid,omitempty"`
	Lot decimal.Decimal `json:"lot,omitempty"`
}

func HandleBid(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Sender != "", "anonymous"); err != nil {
			return err
		}

		f, err := From(r, flips)
		if err != nil {
			return err
		}

		var lot decimal.Decimal
		if err := require(r.Scan(&lot) == nil, "bad-data"); err != nil {
			return err
		}

		vid, _ := uuid.FromString(f.VaultID)
		v, err := vat.From(r.WithBody(vid), vaults)
		if err != nil {
			return err
		}

		cid, _ := uuid.FromString(v.CollateralID)
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
			var err error
			if f.Bid.LessThan(f.Tab) {
				err = Tend(r, c, f, lot)
			} else {
				err = Dent(r, c, f, lot)
			}

			if err != nil {
				return maker.WithFlag(err, maker.FlagRefund)
			}

			var transfers []*core.Transfer

			// 退款给上一个出价的人
			if f.Bid.IsPositive() {
				memo := core.TransferAction{
					ID:     f.TraceID,
					Source: r.Action.String(),
				}.Encode()

				transfers = append(transfers, &core.Transfer{
					TraceID:   uuid.Modify(r.TraceID, memo),
					AssetID:   c.Dai,
					Amount:    f.Bid,
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
				Bid:       r.Amount,
				Lot:       lot,
			}

			if err := flips.CreateEvent(ctx, event); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.CreateEvent")
				return err
			}
		}

		if c.Version < r.Version {
			c.Debt = c.Debt.Sub(event.Bid.Sub(f.Bid))
			c.Ink = c.Ink.Add(f.Lot.Sub(event.Lot))

			if err := collaterals.Update(ctx, c, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		if v.Version < r.Version && event.Lot.LessThan(f.Lot) {
			v.Ink = v.Ink.Add(f.Lot.Sub(event.Lot))

			if err := vaults.Update(ctx, v, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		if f.Version < r.Version {
			f.Action = r.Action
			f.Bid = event.Bid
			f.Lot = event.Lot
			f.Guy = r.Sender
			f.Tic = r.Now.Unix() + c.TTL

			if err := flips.Update(ctx, f, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.Update")
				return err
			}
		}

		return nil
	}
}
