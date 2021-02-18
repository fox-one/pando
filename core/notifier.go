package core

import (
	"context"
)

type Notifier interface {
	Snapshot(ctx context.Context, transfer *Transfer, TxHash string) error
}
