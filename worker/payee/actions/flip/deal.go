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

func Deal(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
) actions.Handler {
	return &dealHandler{
		collaterals: collaterals,
		vaults:      vaults,
		flips:       flips,
	}
}

type dealHandler struct {
	collaterals core.CollateralStore
	vaults      core.VaultStore
	flips       core.FlipStore
}

func (h *dealHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
	log := logger.FromContext(ctx)

	var data flip.Data

	auction, err := h.flips.Find(ctx, tx.TargetID)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, actions.ErrBuildAbort
		}

		log.WithError(err).Errorln("flips.Find")
		return nil, err
	}

	vault, err := h.vaults.Find(ctx, auction.VaultID)
	if err != nil {
		log.WithError(err).Errorln("vaults.Find")
		return nil, err
	}

	cat, err := h.collaterals.Find(ctx, vault.CollateralID)
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return nil, err
	}

	if err := flip.Deal(tx, cat, auction, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *dealHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data flip.Data
	_ = json.Unmarshal(body, &data)

	auction, err := h.flips.Find(ctx, tx.TargetID)
	if err != nil {
		log.WithError(err).Errorln("flips.Find")
		return err
	}

	vault, err := h.vaults.Find(ctx, auction.VaultID)
	if err != nil {
		log.WithError(err).Errorln("vaults.Find")
		return err
	}

	cat, err := h.collaterals.Find(ctx, vault.CollateralID)
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return err
	}

	flip.ApplyDeal(tx, cat, auction, data)

	if err := h.flips.Update(ctx, auction, tx.Version); err != nil {
		log.WithError(err).Errorln("flips.Update")
		return err
	}

	if err := h.collaterals.Update(ctx, cat, tx.Version); err != nil {
		log.WithError(err).Errorln("collaterals.Update")
		return err
	}

	return nil
}
