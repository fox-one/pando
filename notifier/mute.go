package notifier

import (
	"context"

	"github.com/fox-one/pando/core"
)

func Mute() core.Notifier {
	return &dumb{}
}

type dumb struct{}

func (d *dumb) Snapshot(ctx context.Context, transfer *core.Transfer, signedTx string) error {
	return nil
}

func (d *dumb) Auth(ctx context.Context, user *core.User) error {
	return nil
}

func (d *dumb) Transaction(ctx context.Context, tx *core.Transaction) error {
	return nil
}
