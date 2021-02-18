package vat

import (
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Dink decimal.Decimal `json:"dink,omitempty"`
	Dart decimal.Decimal `json:"dart,omitempty"`
	Debt decimal.Decimal `json:"debt,omitempty"`
}
