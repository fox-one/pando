package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string, flags ...int) error {
	return maker.Require(condition, "Vat/"+msg, flags...)
}

func From(r *maker.Request, vaults core.VaultStore) (*core.Vault, error) {
	ctx := r.Context()
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
