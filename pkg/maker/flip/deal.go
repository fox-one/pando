package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
)

func Deal(tx *maker.Tx, _ *core.Collateral, flip *core.Flip, data *Data) error {
	if flip.Action == core.ActionFlipDeal {
		return ErrFlipAlreadyDeal
	}

	finished := !flip.Tic.IsZero() && (tx.Now.After(flip.Tic) || tx.Now.After(flip.End))
	if !finished {
		return ErrFlipNotFinished
	}

	data.Bid = flip.Bid
	data.Lot = flip.Lot

	return nil
}

func ApplyDeal(tx *maker.Tx, cat *core.Collateral, flip *core.Flip, data Data) {
	// cat
	cat.Debt = cat.Debt.Sub(data.Bid)

	// flip
	flip.Action = core.ActionFlipDeal
}
