package main

import (
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/service/asset"
	"github.com/fox-one/pando/service/message"
	"github.com/fox-one/pando/service/oracle"
	"github.com/fox-one/pando/service/user"
	"github.com/fox-one/pando/service/wallet"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/google/wire"
	"golang.org/x/text/language"
)

var serviceSet = wire.NewSet(
	provideMixinClient,
	asset.New,
	message.New,
	provideUserServiceConfig,
	user.New,
	oracle.New,
	provideSystem,
	provideWalletService,
	provideLocalizer,
)

func provideMixinClient(cfg *config.Config) (*mixin.Client, error) {
	return mixin.NewFromKeystore(&cfg.Dapp.Keystore)
}

func provideWalletService(client *mixin.Client, cfg *config.Config, system *core.System) core.WalletService {
	return wallet.New(client, wallet.Config{
		Pin:       cfg.Dapp.Pin,
		Members:   system.MemberIDs(),
		Threshold: system.Threshold,
	})
}

func provideSystem(cfg *config.Config) *core.System {
	members := make([]*core.Member, 0, len(cfg.Group.Members))
	for _, m := range cfg.Group.Members {
		verifyKey, err := mtg.DecodePublicKey(m.VerifyKey)
		if err != nil {
			panic(fmt.Errorf("decode verify key for member %s failed", m.ClientID))
		}

		members = append(members, &core.Member{
			ClientID:  m.ClientID,
			VerifyKey: verifyKey,
		})
	}

	privateKey, err := mtg.DecodePrivateKey(cfg.Group.PrivateKey)
	if err != nil {
		panic(fmt.Errorf("base64 decode group private key failed: %w", err))
	}

	return &core.System{
		ClientID:     cfg.Dapp.ClientID,
		ClientSecret: cfg.Dapp.ClientSecret,
		Members:      members,
		Threshold:    cfg.Group.Threshold,
		PrivateKey:   privateKey,
		Version:      version,
	}
}

func provideUserServiceConfig(cfg *config.Config) user.Config {
	return user.Config{
		ClientSecret: cfg.Dapp.ClientSecret,
	}
}

func provideLocalizer(cfg *config.Config) (*localizer.Localizer, error) {
	files, err := localizer.FindMessageFiles(cfg.I18n.Path)
	if err != nil {
		return nil, err
	}

	lang, err := language.Parse(cfg.I18n.Language)
	if err != nil {
		return nil, err
	}

	return localizer.New(lang, files...), nil
}
