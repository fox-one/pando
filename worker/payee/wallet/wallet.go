package wallet

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
)

func BindTransferVersion(wallets core.WalletStore) core.WalletStore {
	return &withVersion{wallets}
}

type withVersion struct {
	core.WalletStore
}

func (s *withVersion) CreateTransfers(ctx context.Context, transfers []*core.Transfer) error {
	version := maker.VersionFrom(ctx)

	for _, transfer := range transfers {
		transfer.Version = version
	}

	return s.WalletStore.CreateTransfers(ctx, transfers)
}
