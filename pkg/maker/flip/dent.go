package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
)

func Dent(tx *maker.Tx, cat *core.Collateral, vault *core.Vault, flip *core.Flip, opt Option, data *Data) error {
	if !flip.Tic.IsZero() && tx.Now.After(flip.Tic) {
		return ErrFlipFinishedTic
	}

	if tx.Now.After(flip.End) {
		return ErrFlipFinishedEnd
	}

	bid := tx.Amount
	lot := data.Lot.Truncate(8)

	// 必须支付 DAI 并且必须等于 tab
	if match := tx.AssetID == cat.Dai && bid.Equal(flip.Bid) && bid.Equal(flip.Tab); !match {
		return ErrFlipBidNotMatch
	}

	if !lot.IsPositive() || !lot.LessThan(flip.Lot) {
		return ErrFlipLotNotLower
	}

	if lot.Mul(opt.Beg).GreaterThan(flip.Lot) {
		return ErrFlipInsufficientDecrease
	}

	// 退款给上一个出价的人
	if flip.Guy != "" {
		memo := maker.EncodeMemo(module, flip.TraceID, "refund bid")
		tx.Transfer(
			uuid.Modify(tx.TraceID, memo),
			cat.Dai,
			flip.Guy,
			memo,
			flip.Bid,
		)
	}

	// 返回多余拍卖物
	{
		memo := maker.EncodeMemo(module, flip.TraceID, "refund gem")
		tx.Transfer(
			uuid.Modify(tx.TraceID, memo),
			cat.Gem,
			urn.UserID,
			memo,
			flip.Lot.Sub(lot),
		)
	}

	data.Bid = bid
	data.Lot = lot

	return nil
}

func ApplyDent(tx *maker.Tx, flip *core.Flip, opt Option, data Data) {
	// flip
	flip.Action = core.ActionFlipDent
	flip.Bid = data.Bid
	flip.Lot = data.Lot
	flip.Guy = tx.Sender
	flip.Tic = tx.Now.Add(opt.TTL)
}
