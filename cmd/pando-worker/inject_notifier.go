package main

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/notifier"
	"github.com/fox-one/pando/service/asset"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	provideNotifier,
)

func provideNotifier(
	system *core.System,
	assetz core.AssetService,
	messages core.MessageStore,
) core.Notifier {
	if *notify {
		return notifier.New(
			system,
			asset.Cache(assetz),
			messages,
		)
	}

	return notifier.Mute()
}
