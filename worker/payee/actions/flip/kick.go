package flip

import (
	"context"
	"encoding/json"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/flip"
	"github.com/fox-one/pando/worker/payee/actions"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
)

func Kick(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
) actions.Handler {
	return &kickHandler{
		collaterals: collaterals,
		vaults:      vaults,
		flips:       flips,
	}
}

type kickHandler struct {
	collaterals core.CollateralStore
	vaults      core.VaultStore
	flips       core.FlipStore
}

func (h *kickHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
	log := logger.FromContext(ctx)

	var data flip.Data

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

	opt := getOpt()
	if err := flip.Kick(tx, cat, vault, opt, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *kickHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data flip.Data
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

	opt := getOpt()
	r := flip.ApplyKick(tx, cat, vault, opt, data)

	if err := h.flips.Create(ctx, r); err != nil {
		log.WithError(err).Errorln("flips.Create")
		return err
	}

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
