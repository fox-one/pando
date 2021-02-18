package flip

import (
	"context"
	"encoding/json"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/flip"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/worker/payee/actions"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
)

func Dent(
	collaterals core.CollateralStore,
	vaults core.VaultStore,
	flips core.FlipStore,
) actions.Handler {
	return &dentHandler{
		collaterals: collaterals,
		vaults:      vaults,
		flips:       flips,
	}
}

type dentHandler struct {
	collaterals core.CollateralStore
	vaults      core.VaultStore
	flips       core.FlipStore
}

func (h *dentHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
	log := logger.FromContext(ctx)

	var data flip.Data
	_, _ = mtg.Scan(body, &data.Lot)

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
	if err := flip.Dent(tx, cat, vault, auction, opt, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *dentHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data flip.Data
	_ = json.Unmarshal(body, &data)

	auction, err := h.flips.Find(ctx, tx.TargetID)
	if err != nil {
		log.WithError(err).Errorln("flips.Find")
		return err
	}

	opt := getOpt()
	flip.ApplyDent(tx, auction, opt, data)

	if err := h.flips.Update(ctx, auction, tx.Version); err != nil {
		log.WithError(err).Errorln("flips.Update")
		return err
	}

	return nil
}
