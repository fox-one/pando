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

func VaultEvent(event *core.VaultEvent) *api.Vault_Event {
	return &api.Vault_Event{
		VaultId:   event.VaultID,
		CreatedAt: Time(&event.CreatedAt),
		Action:    api.Action(event.Action),
		Dink:      event.Dink.String(),
		Dart:      event.Dart.String(),
		Debt:      event.Debt.String(),
	}
}
