package flip

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/vat"
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
	vaults core.VaultStore,
	flips core.FlipStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
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

		if t.ID == 0 {
			t = r.Tx()
			t.TargetID = f.TraceID

			var transfers []*core.Transfer

			if err := Deal(r, f); err == nil {
				t.Write(core.TxStatusSuccess, DealData{
					Bid: f.Bid,
					Lot: f.Lot,
				})

				if f.Lot.IsPositive() {
					memo := core.TransferAction{
						ID:     f.TraceID,
						Source: "Deal",
					}.Encode()

					transfers = append(transfers, &core.Transfer{
						TraceID:   uuid.Modify(r.TraceID(), memo),
						AssetID:   c.Gem,
						Amount:    f.Lot,
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

		var data DealData
		_ = t.Data.Unmarshal(&data)

		if f.Version < r.Version() {
			f.Action = r.Action

			if err := flips.Update(ctx, f, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("flips.Update")
				return err
			}
		}

		return nil
	}
}

func Deal(r *maker.Request, f *core.Flip) error {
	if err := require(r.Now().After(f.Tic) || r.Now().After(f.End), "not-finished"); err != nil {
		return err
	}

	if err := require(f.Action != r.Action, "dealt"); err != nil {
		return err
	}

	return nil
}
