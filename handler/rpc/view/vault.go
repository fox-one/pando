package view

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
)

func Vault(vat *core.Vault) *api.Vault {
	return &api.Vault{
		Id:           vat.TraceID,
		CreatedAt:    Time(&vat.CreatedAt),
		CollateralId: vat.CollateralID,
		Ink:          vat.Ink.String(),
		Art:          vat.Art.String(),
	}
}
