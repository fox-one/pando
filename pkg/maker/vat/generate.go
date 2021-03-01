package vat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

func HandleGernerate(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	frob := HandleFrob(
		collaterals,
		vaults,
		transactions,
		wallets,
	)

	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.BindUser() == nil && r.BindFollow() == nil, "bad-data"); err != nil {
			return err
		}

		var (
			id   uuid.UUID
			debt decimal.Decimal
		)

		if err := require(r.Scan(&id, &debt) == nil && debt.IsPositive(), "bad-data"); err != nil {
			return err
		}

		r = r.WithBody(types.UUID(r.UserID), types.UUID(r.FollowID), id, decimal.Zero, debt)
		return frob(ctx, r)
	}
}
