package main

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/notifier"
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
	users core.UserStore,
	flips core.FlipStore,
	localizer *localizer.Localizer,
) core.Notifier {
	if _flag.notify {
		return notifier.New(
			system,
			assetz,
			messages,
			vats,
			cats,
			users,
			flips,
			localizer,
		)
	}

	return notifier.Mute()
}
