package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

func HandleGenerated(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	wallets core.WalletStore,
) maker.HandlerFunc {
	frob := HandleFrob(
		collaterals,
		vaults,
		wallets,
	)

	return func(r *maker.Request) error {
		var (
			id   uuid.UUID
			debt decimal.Decimal
		)

		if err := require(r.Scan(&id, &debt) == nil && debt.IsPositive(), "bad-data"); err != nil {
			return err
		}

		r = r.WithBody(id, decimal.Zero, debt)
		return frob(r)
	}
}
