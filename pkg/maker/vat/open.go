package vat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type OpenData struct {
	Ink  decimal.Decimal `json:"ink,omitempty"`
	Debt decimal.Decimal `json:"debt,omitempty"`
}

func HandleOpen(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.BindUser() == nil && r.BindFollow() == nil, "bad-data"); err != nil {
			return err
		}

		c, err := cat.From(ctx, collaterals, r)
		if err != nil {
			return err
		}

		assetID, amount := r.Payment()
		if err := require(assetID == c.Gem, "gem-not-match"); err != nil {
			return err
		}

		t, err := transactions.Find(ctx, r.TraceID())
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("transactions.Find")
			return err
		}

		if t.ID == 0 {
			var debt decimal.Decimal
			if err := require(r.Scan(&debt) == nil, "bad-data"); err != nil {
				return err
			}

			debt = debt.Truncate(8)
			if err := require(debt.IsPositive(), "bad-data"); err != nil {
				return err
			}

			dink := amount
			dart := debt.Div(c.Rate)

			v := &core.Vault{
				CreatedAt:    r.Now(),
				TraceID:      r.TraceID(),
				Version:      r.Version(),
				UserID:       r.UserID,
				CollateralID: c.TraceID,
			}

			t = r.Tx()
			t.TargetID = v.TraceID

			var transfers []*core.Transfer

			if err := frob(c, v, dink, dart); err == nil {
				t.Write(core.TxStatusSuccess, OpenData{
					Ink:  dink,
					Debt: debt,
				})

				v.Ink, v.Art = dink, dart

				if err := vaults.Create(ctx, v); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("vaults.Create")
					return err
				}

				// 放款
				memo := core.TransferAction{
					ID:     t.FollowID,
					Source: "Open",
				}.Encode()

				transfers = append(transfers, &core.Transfer{
					TraceID:   uuid.Modify(t.TraceID, memo),
					AssetID:   c.Dai,
					Amount:    debt,
					Memo:      memo,
					Threshold: 1,
					Opponents: []string{r.UserID},
				})
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

		var data OpenData
		_ = t.Data.Unmarshal(&data)

		dart := data.Debt.Div(c.Rate)
		if c.Version < r.Version() {
			c.Art = c.Art.Add(dart)
			c.Debt = c.Debt.Add(data.Debt)
			c.Ink = c.Ink.Add(data.Ink)

			if err := collaterals.Update(ctx, c, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}
