package flip

import (
	"github.com/shopspring/decimal"
)

type Data struct {
	Bid  decimal.Decimal `json:"bid,omitempty"`
	Lot  decimal.Decimal `json:"lot,omitempty"`
	Dink decimal.Decimal `json:"dink,omitempty"`
	Dart decimal.Decimal `json:"dart,omitempty"`
}
