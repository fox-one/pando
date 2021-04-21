package main

import (
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/service/asset"
	"github.com/fox-one/pando/service/message"
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
	provideSystem,
	provideWalletServiceConfig,
	wallet.New,
	provideLocalizer,
)

func provideMixinClient(cfg *config.Config) (*mixin.Client, error) {
	return mixin.NewFromKeystore(&cfg.Dapp.Keystore)
}

func provideWalletServiceConfig(cfg *config.Config) wallet.Config {
	return wallet.Config{
		Pin:       cfg.Dapp.Pin,
		Members:   cfg.Group.Members,
		Threshold: cfg.Group.Threshold,
	}
}

func provideSystem(cfg *config.Config) *core.System {
	publicKey, err := mtg.DecodePublicKey(cfg.Group.PublicKey)
	if err != nil {
		panic(fmt.Errorf("base64 decode group private key failed: %w", err))
	}

	return &core.System{
		ClientID:     cfg.Dapp.ClientID,
		ClientSecret: cfg.Dapp.ClientSecret,
		Members:      cfg.Group.Members,
		Threshold:    cfg.Group.Threshold,
		PublicKey:    publicKey,
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
