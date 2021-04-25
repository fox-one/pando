package keeper

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/uuid"
)

func (w *Keeper) foldCats(ctx context.Context) error {
	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-time.After(dur):
			trace := uuid.MD5(t.Truncate(time.Minute).Format(time.RFC3339Nano))

			if err := w.handleTransfer(ctx, trace, core.ActionCatFold, uuid.Zero); err != nil {
				dur = time.Second
			} else {
				dur = time.Minute
			}
		}
	}
}
