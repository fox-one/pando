package main

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/parliament"
	"github.com/fox-one/pando/service/asset"
	"github.com/google/wire"
)

var parliamentSet = wire.NewSet(
	provideParliament,
)

func provideParliament(
	messages core.MessageStore,
	userz core.UserService,
	assetz core.AssetService,
	walletz core.WalletService,
	collaterals core.CollateralStore,
	system *core.System,
) core.Parliament {
	return parliament.New(
		messages,
		userz,
		asset.Cache(assetz),
		walletz,
		collaterals,
		system,
	)
}
