package vat

import (
	"github.com/shopspring/decimal"
)

type Data struct {
	Dink decimal.Decimal `json:"dink,omitempty"`
	Debt decimal.Decimal `json:"debt,omitempty"`
	Dart decimal.Decimal `json:"dart,omitempty"`
}
