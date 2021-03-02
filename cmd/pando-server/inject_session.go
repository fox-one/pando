package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/session"
	"github.com/google/wire"
)

var sessionSet = wire.NewSet(
	provideSessionConfig,
	session.New,
)

func provideSessionConfig(cfg *config.Config) session.Config {
	var issuers []string
	for _, m := range cfg.Group.Members {
		issuers = append(issuers, m)
	}

	if !govalidator.IsIn(cfg.Dapp.ClientID, issuers...) {
		issuers = append(issuers, cfg.Dapp.ClientID)
	}

	return session.Config{
		Capacity: 2048,
		Issuers:  issuers,
	}
}
