package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
)

func Tend(tx *maker.Tx, cat *core.Collateral, flip *core.Flip, opt Option, data *Data) error {
	if !flip.Tic.IsZero() && tx.Now.After(flip.Tic) {
		return ErrFlipFinishedTic
	}

	if tx.Now.After(flip.End) {
		return ErrFlipFinishedEnd
	}

	bid := tx.Amount

	// 必须支付 DAI 并且不能大于 tab
	if tx.AssetID != cat.Dai || bid.GreaterThan(flip.Tab) {
		return ErrFlipBidNotMatch
	}

	// 出价必须比上一个高
	if bid.LessThanOrEqual(flip.Bid) {
		return ErrFlipBidNotHigher
	}

	if bid.LessThan(flip.Tab) && bid.LessThan(flip.Bid.Mul(opt.Beg)) {
		return ErrFlipInsufficientIncrease
	}

	// 退款给上一个出价的人
	if flip.Guy != "" && flip.Bid.IsPositive() {
		memo := "flip: defeated by new bid"
		tx.Transfer(
			uuid.Modify(tx.TraceID, memo),
			cat.Dai,
			flip.Guy,
			memo,
			flip.Bid,
		)
	}

	data.Bid = bid
	data.Lot = flip.Lot

	return nil
}

func ApplyTend(tx *maker.Tx, flip *core.Flip, opt Option, data Data) {
	// flip
	flip.Action = core.ActionFlipTend
	flip.Bid = data.Bid
	flip.Lot = data.Lot
	flip.Guy = tx.Sender
	flip.Tic = tx.Now.Add(opt.TTL)
}
