package flip

import (
	"github.com/fox-one/pando/pkg/maker"
)

const module = "flip"

var (
	ErrFlipNotUnsafe            = maker.RegisterErr(1101, module, "not-unsafe")
	ErrFlipNullAuction          = maker.RegisterErr(1102, module, "null-auction")
	ErrFlipBidNotMatch          = maker.RegisterErr(1103, module, "bid-not-match")
	ErrFlipBidNotHigher         = maker.RegisterErr(1104, module, "bid-not-higher")
	ErrFlipLotNotLower          = maker.RegisterErr(1105, module, "lot-not-lower")
	ErrFlipFinishedTic          = maker.RegisterErr(1106, module, "finished-tic")
	ErrFlipFinishedEnd          = maker.RegisterErr(1107, module, "finished-end")
	ErrFlipInsufficientIncrease = maker.RegisterErr(1108, module, "insufficient-increase")
	ErrFlipInsufficientDecrease = maker.RegisterErr(1109, module, "insufficient-decrease")
	ErrFlipNotLive              = maker.RegisterErr(1110, module, "not-live")
	ErrFlipNotFinished          = maker.RegisterErr(1111, module, "not-finished")
)
