package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/shopspring/decimal"
)

const year int64 = 60 * 60 * 24 * 365

type FoldData struct {
	Rate     decimal.Decimal `json:"rate,omitempty"`
	GemPrice decimal.Decimal `json:"gem_price,omitempty"`
	DaiPrice decimal.Decimal `json:"dai_price,omitempty"`
}

// Fold modify the debt multiplier, creating / destroying corresponding debt
func Fold(tx *maker.Tx, cat *core.Collateral, data *FoldData) error {
	if !cat.Live {
		return ErrVatNotLive
	}

	n := tx.Now.Unix() - cat.Rho.Unix()
	if n > 0 {
		q := decimal.NewFromInt(n).Div(decimal.NewFromInt(year))
		f := number.Pow(cat.Duty, q)
		data.Rate = cat.Rate.Mul(f).Sub(cat.Rate).Truncate(16)
	} else {
		data.Rate = decimal.Zero
	}

	return nil
}

func ApplyFold(tx *maker.Tx, cat *core.Collateral, data FoldData) {
	cat.Rate = cat.Rate.Add(data.Rate)
	cat.Rho = tx.Now

	if data.GemPrice.IsPositive() && data.DaiPrice.IsPositive() {
		cat.Price = data.GemPrice.Div(data.DaiPrice).Truncate(8)
	}
}
