package main

import (
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/store/asset"
	"github.com/fox-one/pando/store/collateral"
	"github.com/fox-one/pando/store/flip"
	"github.com/fox-one/pando/store/message"
	"github.com/fox-one/pando/store/oracle"
	"github.com/fox-one/pando/store/proposal"
	"github.com/fox-one/pando/store/transaction"
	"github.com/fox-one/pando/store/user"
	"github.com/fox-one/pando/store/vault"
	"github.com/fox-one/pando/store/wallet"
	"github.com/fox-one/pkg/store/db"
	"github.com/google/wire"
)

var storeSet = wire.NewSet(
	provideDatabase,
	asset.New,
	collateral.New,
	flip.New,
	proposal.New,
	transaction.New,
	user.New,
	vault.New,
	oracle.New,
	wallet.New,
	message.New,
)

func provideDatabase(cfg *config.Config) (*db.DB, error) {
	database, err := db.Open(cfg.DB)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(database); err != nil {
		return nil, err
	}

	return database, nil
}
