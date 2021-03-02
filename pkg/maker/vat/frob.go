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

func HandleFrob(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.BindUser() == nil && r.BindFollow() == nil, "bad-data"); err != nil {
			return err
		}

		v, err := From(ctx, vaults, r)
		if err != nil {
			return err
		}

		if err := require(v.UserID == r.UserID, "not-authorized"); err != nil {
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
			t.TargetID = v.TraceID

			var dink, debt decimal.Decimal
			if err := require(r.Scan(&dink, &debt) == nil, "bad-data"); err != nil {
				return err
			}

			dink = dink.Truncate(8)
			debt = debt.Truncate(8)

			assetID, amount := r.Payment()
			if err := require(!dink.IsPositive() || (assetID == c.Gem && dink.Equal(amount)), "bad-data"); err != nil {
				return err
			}

			if err := require(!debt.IsNegative() || (assetID == c.Dai && debt.Neg().Equal(amount)), "bad-data"); err != nil {
				return err
			}

			var transfers []*core.Transfer

			dart := debt.Div(c.Rate)
			if err := frob(c, v, dink, dart); err == nil {
				t.Write(core.TxStatusSuccess, Data{
					Dink: dink,
					Debt: debt,
					Dart: dart,
				})

				// 提取抵押物
				if dink.IsNegative() {
					memo := core.TransferAction{
						ID:     r.FollowID,
						Source: "Withdraw",
					}.Encode()

					transfers = append(transfers, &core.Transfer{
						TraceID:   uuid.Modify(t.TraceID, memo),
						AssetID:   c.Gem,
						Amount:    dink.Abs(),
						Memo:      memo,
						Threshold: 1,
						Opponents: []string{r.UserID},
					})
				}

				// 借出新的币
				if debt.IsPositive() {
					memo := core.TransferAction{
						ID:     r.FollowID,
						Source: "Generate",
					}.Encode()

					transfers = append(transfers, &core.Transfer{
						TraceID:   uuid.Modify(t.TraceID, memo),
						AssetID:   c.Dai,
						Amount:    debt.Abs(),
						Memo:      memo,
						Threshold: 1,
						Opponents: []string{r.UserID},
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

		var data Data
		_ = t.Data.Unmarshal(&data)

		// update vat
		if v.Version < r.Version() {
			v.Art = v.Art.Add(data.Dart)
			v.Ink = v.Ink.Add(data.Dink)

			if err := vaults.Update(ctx, v, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		// update cat
		if c.Version < r.Version() {
			c.Art = c.Art.Add(data.Dart)
			c.Debt = c.Debt.Add(data.Debt)
			c.Ink = c.Ink.Add(data.Dink)

			if err := collaterals.Update(ctx, c, r.Version()); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}

// Frob modify a Vault
func frob(c *core.Collateral, v *core.Vault, dink, dart decimal.Decimal) error {
	if err := require(dart.IsNegative() || c.Art.Add(dart).Mul(c.Rate).LessThanOrEqual(c.Line), "ceiling-exceeded"); err != nil {
		return err
	}

	ink, art := v.Ink.Add(dink), v.Art.Add(dart)
	tab := art.Mul(c.Rate).Truncate(8)

	if err := require(ink.Mul(c.Price).GreaterThanOrEqual(tab.Mul(c.Mat)), "not-safe"); err != nil {
		return err
	}

	if err := require(tab.IsZero() || tab.GreaterThanOrEqual(c.Dust), "dust"); err != nil {
		return err
	}

	return nil
}
