package vat

import (
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
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		v, err := From(r, vaults)
		if err != nil {
			return err
		}

		cid, _ := uuid.FromString(v.CollateralID)
		c, err := cat.From(r.WithBody(cid), collaterals)
		if err != nil {
			return err
		}

		if err := require(c.Live > 0, "not-live"); err != nil {
			return maker.WithFlag(err, maker.FlagRefund)
		}

		event, err := vaults.FindEvent(ctx, v.TraceID, r.Version)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("vaults.FindEvent")
			return err
		}

		if event.ID == 0 {
			var dink, debt decimal.Decimal
			if err := require(r.Scan(&dink, &debt) == nil, "bad-data"); err != nil {
				return err
			}

			dink = dink.Truncate(8)
			debt = debt.Truncate(8)

			if dink.IsPositive() { // 增加抵押
				if err := require(r.AssetID == c.Gem && dink.Equal(r.Amount), "gem-not-match"); err != nil {
					return err
				}
			}

			if dink.IsNegative() { // 提取抵押物
				if err := require(r.Sender == v.UserID, "not-authorized"); err != nil {
					return err
				}
			}

			if debt.IsPositive() { // 增加借贷
				if err := require(r.Sender == v.UserID, "not-authorized"); err != nil {
					return err
				}
			}

			if debt.IsNegative() { // 还贷
				if err := require(r.AssetID == c.Dai && debt.Abs().Equal(r.Amount), "dai-not-match"); err != nil {
					return err
				}
			}

			if dink.IsZero() && debt.IsZero() {
				return nil
			}

			dart := debt.Div(c.Rate)
			if dart.IsNegative() && v.Art.Add(dart).Mul(c.Rate).Truncate(8).IsZero() {
				dart = v.Art.Neg()
			} else if dart.IsPositive() && dart.Mul(c.Rate).LessThan(debt) {
				dart = dart.Add(decimal.New(1, -16))
			}

			if err := frob(c, v, dink, dart); err != nil {
				return maker.WithFlag(err, maker.FlagRefund)
			}

			var transfers []*core.Transfer

			// 提取抵押物
			if dink.IsNegative() {
				memo := core.TransferAction{
					ID:     r.FollowID,
					Source: core.ActionVatWithdraw.String(),
				}.Encode()

				transfers = append(transfers, &core.Transfer{
					TraceID:   uuid.Modify(r.TraceID, memo),
					AssetID:   c.Gem,
					Amount:    dink.Abs(),
					Memo:      memo,
					Threshold: 1,
					Opponents: []string{v.UserID},
				})
			}

			// 借出新的币
			if debt.IsPositive() {
				memo := core.TransferAction{
					ID:     r.FollowID,
					Source: core.ActionVatGenerate.String(),
				}.Encode()

				transfers = append(transfers, &core.Transfer{
					TraceID:   uuid.Modify(r.TraceID, memo),
					AssetID:   c.Dai,
					Amount:    debt.Abs(),
					Memo:      memo,
					Threshold: 1,
					Opponents: []string{v.UserID},
				})
			}

			if err := wallets.CreateTransfers(ctx, transfers); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}

			event = &core.VaultEvent{
				CreatedAt: r.Now,
				VaultID:   v.TraceID,
				Version:   r.Version,
				Action:    r.Action,
				Dink:      dink,
				Dart:      dart,
				Debt:      debt,
			}

			if err := vaults.CreateEvent(ctx, event); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.CreateEvent")
				return err
			}
		}

		// update vat
		if v.Version < r.Version {
			v.Ink = v.Ink.Add(event.Dink)
			v.Art = v.Art.Add(event.Dart)

			if err := vaults.Update(ctx, v, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		// update cat
		if c.Version < r.Version {
			c.Ink = c.Ink.Add(event.Dink)
			c.Art = c.Art.Add(event.Dart)
			c.Debt = c.Debt.Add(event.Debt)

			if err := collaterals.Update(ctx, c, r.Version); err != nil {
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
	tab := art.Mul(c.Rate)

	if err := require(ink.Mul(c.Price).GreaterThanOrEqual(tab.Mul(c.Mat)), "not-safe"); err != nil {
		return err
	}

	if err := require(tab.IsZero() || tab.GreaterThanOrEqual(c.Dust), "dust"); err != nil {
		return err
	}

	return nil
}
