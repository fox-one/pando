package vat

import (
	"context"
	"encoding/json"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/worker/payee/actions"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
)

func Fold(
	collaterals core.CollateralStore,
	oracles core.OracleStore,
) actions.Handler {
	return &foldHandler{
		collaterals: collaterals,
		oracles:     oracles,
	}
}

type foldHandler struct {
	collaterals core.CollateralStore
	oracles     core.OracleStore
}

func (h *foldHandler) Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error) {
	log := logger.FromContext(ctx)

	var data vat.FoldData

	cat, err := h.collaterals.Find(ctx, tx.TargetID)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, actions.ErrBuildAbort
		}

		return nil, err
	}

	at := tx.Now
	g, err := h.oracles.Find(ctx, cat.Gem, at)
	if err != nil {
		log.WithError(err).Errorf("oracles.Find(%q)", cat.Gem)
		return nil, err
	}

	d, err := h.oracles.Find(ctx, cat.Dai, at)
	if err != nil {
		log.WithError(err).Errorf("oracles.Find(%q)", cat.Gem)
		return nil, err
	}

	data.GemPrice = g.Price
	data.DaiPrice = d.Price

	if err := vat.Fold(tx, cat, &data); err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (h *foldHandler) Apply(ctx context.Context, tx *maker.Tx, body []byte) error {
	log := logger.FromContext(ctx)

	var data vat.FoldData
	_ = json.Unmarshal(body, &data)

	cat, err := h.collaterals.Find(ctx, tx.TargetID)
	if err != nil {
		log.WithError(err).Errorln("collaterals.Find")
		return err
	}

	vat.ApplyFold(tx, cat, data)

	if err := h.collaterals.Update(ctx, cat, tx.Version); err != nil {
		log.WithError(err).Errorln("collaterals.Update")
		return err
	}

	return nil
}
