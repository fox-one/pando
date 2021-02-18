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

func Open(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
) actions.Handler {
	return &openHandler{
		collaterals: collaterals,
		vaults:      vaults,
	}
}

type openHandler struct {
	collaterals core.CollateralStore
	vaults      core.VaultStore
}

func (h *openHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
	log := logger.FromContext(ctx)

	var data vat.FrobData
	_, _ = mtg.Scan(body, &data.Debt)

	cat, err := h.collaterals.Find(ctx, tx.TargetID)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, actions.ErrBuildAbort
		}

		log.WithError(err).Errorln("collaterals.Find")
		return nil, err
	}

	if err := vat.Open(tx, cat, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *openHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data vat.FrobData
	_ = json.Unmarshal(body, &data)

	cat, err := h.collaterals.Find(ctx, tx.TargetID)
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return err
	}

	vault := vat.ApplyOpen(tx, cat, data)

	if err := h.vaults.Create(ctx, vault); err != nil {
		log.WithError(err).Errorln("vaults.Create")
		return err
	}

	if err := h.collaterals.Update(ctx, cat, tx.Version); err != nil {
		log.WithError(err).Errorln("collaterals.Update")
		return err
	}

	return nil
}
