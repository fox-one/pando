package vat

import (
	"context"
	"encoding/json"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/worker/payee/actions"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
)

func Frob(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
) actions.Handler {
	return &frobHandler{
		collaterals: collaterals,
		vaults:      vaults,
	}
}

type frobHandler struct {
	collaterals core.CollateralStore
	vaults      core.VaultStore
}

func (h *frobHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
	log := logger.FromContext(ctx)

	var data vat.FrobData
	_, _ = mtg.Scan(body, &data.Dink, &data.Debt)

	vault, err := h.vaults.Find(ctx, tx.TargetID)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, actions.ErrBuildAbort
		}

		log.WithError(err).Errorln("vaults.Find")
		return nil, err
	}

	cat, err := h.collaterals.Find(ctx, vault.CollateralID)
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return nil, err
	}

	if err := vat.Frob(tx, cat, vault, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *frobHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data vat.FrobData
	_ = json.Unmarshal(body, &data)

	vault, err := h.vaults.Find(ctx, tx.TargetID)
	if err != nil {
		log.WithError(err).Errorln("vaults.Find")
		return err
	}

	cat, err := h.collaterals.Find(ctx, vault.CollateralID)
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return err
	}

	vat.ApplyFrob(tx, cat, vault, data)

	if err := h.vaults.Update(ctx, vault, tx.Version); err != nil {
		log.WithError(err).Errorln("vaults.Update")
		return err
	}

	if err := h.collaterals.Update(ctx, cat, tx.Version); err != nil {
		log.WithError(err).Errorln("collaterals.Update")
		return err
	}

	return nil
}
