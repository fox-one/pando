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

func Tend(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
) actions.Handler {
	return &tendHandler{
		collaterals: collaterals,
		vaults:      vaults,
		flips:       flips,
	}
}

type tendHandler struct {
	collaterals core.CollateralStore
	vaults      core.VaultStore
	flips       core.FlipStore
}

func (h *tendHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
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

	opt := getOpt()
	if err := flip.Tend(tx, cat, auction, opt, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *tendHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data flip.Data
	_ = json.Unmarshal(body, &data)

	auction, err := h.flips.Find(ctx, tx.TargetID)
	if err != nil {
		log.WithError(err).Errorln("flips.Find")
		return err
	}

	opt := getOpt()
	flip.ApplyTend(tx, auction, opt, data)

	if err := h.flips.Update(ctx, auction, tx.Version); err != nil {
		log.WithError(err).Errorln("flips.Update")
		return err
	}

	return nil
}
