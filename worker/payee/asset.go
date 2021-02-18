package payee

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store"
)

func (w *Payee) findAsset(ctx context.Context, id string) (*core.Asset, error) {
	asset, err := w.assets.Find(ctx, id)
	if err == nil {
		return asset, nil
	}

	if !store.IsErrNotFound(err) {
		return nil, err
	}

	asset, err = w.assetz.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := w.assets.Create(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}
