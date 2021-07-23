package cat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleGain(
	collaterals core.CollateralStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	return func(r *maker.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		c, err := From(r, collaterals)
		if err != nil {
			return err
		}

		var (
			amount  decimal.Decimal
			receipt uuid.UUID
		)

		if err := require(r.Scan(&amount, &receipt) == nil && amount.Truncate(8).IsPositive(), "bad-data"); err != nil {
			return err
		}

		if c.Version >= r.Version {
			return nil
		}

		amount = amount.Truncate(8)
		max := decimal.Min(
			c.Line,
			c.Art.Mul(c.Rate).Truncate(8),
		).Sub(c.Debt)

		if err := require(amount.LessThanOrEqual(max), "insufficient-profit"); err != nil {
			return err
		}

		memo := core.TransferAction{
			ID:     r.FollowID,
			Source: r.Action.String(),
		}.Encode()

		transfer := &core.Transfer{
			TraceID:   uuid.Modify(r.TraceID, memo),
			AssetID:   c.Dai,
			Amount:    amount,
			Memo:      memo,
			Threshold: 1,
			Opponents: []string{receipt.String()},
		}

		if err := wallets.CreateTransfers(ctx, []*core.Transfer{transfer}); err != nil {
			log.WithError(err).Errorln("wallets.CreateTransfers")
			return err
		}

		c.Debt = c.Debt.Add(amount)
		if err := collaterals.Update(ctx, c, r.Version); err != nil {
			log.WithError(err).Errorln("collaterals.Update")
			return err
		}

		return nil
	}
}
