package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

func HandleWithdraw(
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
			dink decimal.Decimal
		)

		if err := require(r.Scan(&id, &dink) == nil && dink.IsPositive(), "bad-data"); err != nil {
			return err
		}

		r = r.WithBody(id, dink.Neg(), decimal.Zero)
		return frob(r)
	}
}
