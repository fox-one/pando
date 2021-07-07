package main

import (
	"github.com/fox-one/pando/cmd/pando-worker/config"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/notifier"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	provideNotifyConfig,
	provideNotifier,
)

func provideNotifyConfig(cfg *config.Config) notifier.Config {
	return notifier.Config{
		Links: map[string]string{
			"flip_detail":  cfg.Flip.DetailPage,
			"vault_detail": cfg.Vault.DetailPage,
		},
	}
}

func provideNotifier(
	system *core.System,
	assetz core.AssetService,
	messages core.MessageStore,
	vats core.VaultStore,
	cats core.CollateralStore,
	users core.UserStore,
	flips core.FlipStore,
	localizer *localizer.Localizer,
	cfg notifier.Config,
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
			cfg,
		)
	}

	return notifier.Mute()
}
