package core

import (
	"context"
)

// Notifier define operations to send notification to users
type Notifier interface {
	// Auth called when a user login successfully
	Auth(ctx context.Context, user *User) error
	// Transaction called when a new tx created
	Transaction(ctx context.Context, tx *Transaction) error
	// Snapshot called when a transfer confirmed by mixin main net
	Snapshot(ctx context.Context, transfer *Transfer, TxHash string) error
	// VaultUnsafe notify the vault's owner the risk
	VaultUnsafe(ctx context.Context, cat *Collateral, vault *Vault) error
}
