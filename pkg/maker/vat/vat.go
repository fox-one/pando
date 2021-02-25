package vat

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "Vat/%s", msg)
}

func From(ctx context.Context, vaults core.VaultStore, r *maker.Request) (*core.Vault, error) {
	log := logger.FromContext(ctx)

	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	v, err := vaults.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("vaults.Find")
		return nil, err
	}

	if err := require(v.ID > 0, "not init"); err != nil {
		return nil, err
	}

	return v, nil
}
