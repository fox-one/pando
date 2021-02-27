package vat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

func HandleWithdraw(
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
		var (
			user   uuid.UUID
			follow uuid.UUID
			id     uuid.UUID
			dink   decimal.Decimal
		)

		if err := require(r.Scan(&user, &follow, &id, &dink) == nil && dink.IsPositive(), "bad-data"); err != nil {
			return err
		}

		r = r.WithBody(user, follow, id, dink.Neg(), decimal.Zero)
		return frob(ctx, r)
	}
}
