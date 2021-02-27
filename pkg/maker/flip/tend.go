package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/shopspring/decimal"
)

func Tend(r *maker.Request, c *core.Collateral, f *core.Flip, lot decimal.Decimal, opt *Option) error {
	if err := require(r.Now().Before(f.Tic), "finished-tic"); err != nil {
		return err
	}

	if err := require(r.Now().Before(f.End), "finished-end"); err != nil {
		return err
	}

	assetID, bid := r.Payment()
	if err := require(assetID == c.Dai && bid.LessThanOrEqual(f.Tab), "bid-not-match"); err != nil {
		return err
	}

	if err := require(f.Lot.Equal(lot), "lot-not-match"); err != nil {
		return err
	}

	if err := require(bid.GreaterThan(f.Bid), "bid-not-higher"); err != nil {
		return err
	}

	if err := require(bid.Equal(f.Tab) || bid.GreaterThanOrEqual(f.Bid.Mul(opt.Beg)), "insufficient-increase"); err != nil {
		return err
	}

	return nil
}
