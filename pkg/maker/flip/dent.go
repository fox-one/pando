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

type DentData struct {
	Bid decimal.Decimal `json:"bid,omitempty"`
	Lot decimal.Decimal `json:"lot,omitempty"`
}

func HandleDent(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
	properties property.Store,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.BindUser() == nil && r.BindFollow() == nil, "bad-data"); err != nil {
			return err
		}

		var lot decimal.Decimal
		if err := require(r.Scan(&lot) == nil, "bad-data"); err != nil {
			return err
		}

		f, err := From(ctx, flips, r)
		if err != nil {
			return err
		}

		vid, _ := uuid.FromString(f.VaultID)
		v, err := vat.From(ctx, vaults, r.WithBody(vid))
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

		opt, err := ReadOptions(ctx, properties)
		if err != nil {
			return err
		}

		if t.ID == 0 {
			t := r.Tx()
			t.TargetID = f.TraceID

			var transfers []*core.Transfer

			if err := Dent(r, c, f, lot, opt); err == nil {
				_, bid := r.Payment()

				t.Write(core.TxStatusSuccess, TendData{
					Bid: bid,
					Lot: lot,
				})

				// 退款给上一个出价的人
				if f.Bid.IsPositive() {
					memo := core.TransferAction{
						ID:     f.TraceID,
						Source: "RefundBid",
					}.Encode()

					transfers = append(transfers, &core.Transfer{
						TraceID:   uuid.Modify(t.TraceID, memo),
						AssetID:   c.Dai,
						Amount:    f.Bid,
						Memo:      memo,
						Threshold: 1,
						Opponents: []string{f.Guy},
					})
				}
			} else {
				t.Write(core.TxStatusFailed, err)
			}

			if err := wallets.CreateTransfers(ctx, transfers); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}

			if err := transactions.Create(ctx, t); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("transactions.Create")
				return err
			}
		}

		if err := require(t.Status == core.TxStatusSuccess, "tx-failed"); err != nil {
			return err
		}

		var data TendData
		_ = t.Data.Unmarshal(&data)

		if v.Version < r.Version() {
			v.Ink = v.Ink.Add(f.Lot.Sub(data.Lot))

			if err := vaults.Update(ctx, v, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		if f.Version < r.Version() {
			f.Action = r.Action
			f.Bid = data.Bid
			f.Lot = data.Lot
			f.Guy = r.UserID
			f.Tic = r.Now().Add(opt.TTL)

			if err := flips.Update(ctx, f, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.Update")
				return err
			}
		}

		return nil
	}
}

func Dent(r *maker.Request, c *core.Collateral, f *core.Flip, lot decimal.Decimal, opt *Option) error {
	if err := require(r.Now().Before(f.Tic), "finished-tic"); err != nil {
		return err
	}

	if err := require(r.Now().Before(f.End), "finished-end"); err != nil {
		return err
	}

	assetID, bid := r.Payment()
	if err := require(assetID == c.Dai && bid.Equal(f.Bid) && bid.Equal(f.Tab), "bid-not-match"); err != nil {
		return err
	}

	if err := require(lot.IsPositive() && lot.LessThan(f.Lot), "lot-not-lower"); err != nil {
		return err
	}

	if err := require(lot.Mul(opt.Beg).LessThanOrEqual(f.Lot), "insufficient-decrease"); err != nil {
		return err
	}

	return nil
}
