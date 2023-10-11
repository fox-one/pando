package main

import (
	"encoding/base64"

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
	issuers := cfg.Session.Issuers
	for _, m := range cfg.Group.Members {
		if !govalidator.IsIn(m, issuers...) {
			issuers = append(issuers, m)
		}
	}

	if !govalidator.IsIn(cfg.Dapp.ClientID, issuers...) {
		issuers = append(issuers, cfg.Dapp.ClientID)
	}

	secret, err := base64.StdEncoding.DecodeString(cfg.Session.JwtSecret)
	if err != nil {
		panic(err)
	}

	return session.Config{
		Capacity:  2048,
		Issuers:   issuers,
		JwtSecret: secret,
	}
}
