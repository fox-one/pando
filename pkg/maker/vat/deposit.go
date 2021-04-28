package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/property"
	"github.com/shopspring/decimal"
)

func HandleDeposit(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	wallets core.WalletStore,
	properties property.Store,
) maker.HandlerFunc {
	frob := HandleFrob(
		collaterals,
		vaults,
		wallets,
		properties,
	)

	return func(r *maker.Request) error {
		var id uuid.UUID
		if err := require(r.Scan(&id) == nil, "bad-data", maker.FlagNoisy); err != nil {
			return err
		}

		dink := r.Amount
		r = r.WithBody(id, dink, decimal.Zero)
		return frob(r)
	}
}
