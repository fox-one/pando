package core

import (
	"context"
)

type Notifier interface {
	Auth(ctx context.Context, user *User) error
	Transaction(ctx context.Context, tx *Transaction) error
	Snapshot(ctx context.Context, transfer *Transfer, TxHash string) error
}
