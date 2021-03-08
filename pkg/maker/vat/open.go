package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleOpen(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()

		if err := require(r.Sender != "", "anonymous"); err != nil {
			return err
		}

		c, err := cat.From(r, collaterals)
		if err != nil {
			return err
		}

		if err := require(c.Live > 0, "not-live"); err != nil {
			return maker.WithFlag(err, maker.FlagRefund)
		}

		if err := require(r.AssetID == c.Gem, "gem-not-match"); err != nil {
			return err
		}

		event, err := vaults.FindEvent(ctx, r.TraceID, r.Version)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("vaults.FindEvent")
			return err
		}

		if event.ID == 0 {
			var debt decimal.Decimal
			if err := require(r.Scan(&debt) == nil, "bad-data"); err != nil {
				return err
			}

			debt = debt.Truncate(8)
			if err := require(debt.IsPositive(), "bad-data"); err != nil {
				return err
			}

			dink := r.Amount
			dart := debt.Div(c.Rate)

			v := &core.Vault{
				CreatedAt:    r.Now,
				TraceID:      r.TraceID,
				Version:      r.Version,
				UserID:       r.Sender,
				CollateralID: c.TraceID,
			}

			if err := frob(c, v, dink, dart); err != nil {
				return maker.WithFlag(err, maker.FlagRefund)
			}

			// handle transfer
			{
				var transfers []*core.Transfer

				// 放款
				memo := core.TransferAction{
					ID:     r.FollowID,
					Source: r.Action.String(),
				}.Encode()

				transfers = append(transfers, &core.Transfer{
					TraceID:   uuid.Modify(r.TraceID, memo),
					AssetID:   c.Dai,
					Amount:    debt,
					Memo:      memo,
					Threshold: 1,
					Opponents: []string{r.Sender},
				})

				if err := wallets.CreateTransfers(ctx, transfers); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("wallets.CreateTransfers")
					return err
				}
			}

			// create vault
			{
				v.Ink, v.Art = dink, dart

				if err := vaults.Create(ctx, v); err != nil {
					logger.FromContext(ctx).WithError(err).Errorln("vaults.Create")
					return err
				}
			}

			// create vault event
			{
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
		}

		if c.Version < r.Version {
			c.Art = c.Art.Add(event.Dart)
			c.Debt = c.Debt.Add(event.Debt)
			c.Ink = c.Ink.Add(event.Dink)

			if err := collaterals.Update(ctx, c, r.Version); err != nil {
				logger.FromContext(ctx).WithError(err).Errorln("collaterals.Update")
				return err
			}
		}

		return nil
	}
}
