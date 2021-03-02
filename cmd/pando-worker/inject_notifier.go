package main

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/notifier"
	"github.com/fox-one/pando/service/asset"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	provideNotifier,
)

func provideNotifier(
	system *core.System,
	assetz core.AssetService,
	messages core.MessageStore,
	vats core.VaultStore,
	cats core.CollateralStore,
	localizer *localizer.Localizer,
) core.Notifier {
	if *notify {
		return notifier.New(
			system,
			asset.Cache(assetz),
			messages,
			vats,
			cats,
			localizer,
		)
	}

	return notifier.Mute()
}
