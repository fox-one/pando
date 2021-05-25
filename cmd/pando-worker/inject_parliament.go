package main

import (
	"github.com/fox-one/pando/cmd/pando-worker/config"
	"github.com/fox-one/pando/parliament"
	"github.com/google/wire"
)

var parliamentSet = wire.NewSet(
	provideParliamentConfig,
	parliament.New,
)

func provideParliamentConfig(cfg *config.Config) parliament.Config {
	return parliament.Config{
		Links: map[string]string{
			"flip_detail": cfg.Flip.DetailPage,
		},
	}
}
