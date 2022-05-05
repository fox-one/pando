package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pando/store/stat"

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
	stat.New,
)

func connectDatabase(cfg db.Config, timeout time.Duration) (*db.DB, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("connect db: %w", ctx.Err())
		case <-time.After(dur):
			if conn, err := db.Open(cfg); err == nil {
				return conn, nil
			}

			dur = time.Second
		}
	}
}

func provideDatabase(cfg *config.Config) (*db.DB, error) {
	conn, err := connectDatabase(cfg.DB, 8*time.Second)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(conn); err != nil {
		return nil, err
	}

	return conn, nil
}
