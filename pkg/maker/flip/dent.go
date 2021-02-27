package flip

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/shopspring/decimal"
)

func Dent(r *maker.Request, c *core.Collateral, f *core.Flip, lot decimal.Decimal, opt *Option) error {
	if err := require(r.Now().Before(f.Tic), "finished-tic"); err != nil {
		return err
	}

	if err := require(r.Now().Before(f.End), "finished-end"); err != nil {
		return err
	}

	assetID, bid := r.Payment()
	if err := require(assetID == c.Dai && bid.Equal(f.Bid) && bid.Equal(f.Tab), "bid-not-match"); err != nil {
		return err
	}

	if err := require(lot.IsPositive() && lot.LessThan(f.Lot), "lot-not-lower"); err != nil {
		return err
	}

	if err := require(lot.Mul(opt.Beg).LessThanOrEqual(f.Lot), "insufficient-decrease"); err != nil {
		return err
	}

	return nil
}
