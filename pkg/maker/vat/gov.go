package vat

import (
	"strconv"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/jmoiron/sqlx/types"
)

type GovData struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func Gov(tx *maker.Tx, cat *core.Collateral, data *GovData) error {
	if !tx.Gov {
		return ErrVatNotAllowed
	}

	switch data.Key {
	case "line", "dust", "dunk", "duty", "chop", "mat":
		if number.Decimal(data.Value).Truncate(8).IsPositive() {
			return nil
		}
	case "live":
		if _, err := strconv.ParseBool(data.Value); err == nil {
			return nil
		}
	}

	return ErrVatValidateFailed
}

func ApplyGov(tx *maker.Tx, cat *core.Collateral, data GovData) {
	switch data.Key {
	case "line":
		cat.Line = number.Decimal(data.Value).Truncate(8)
	case "dust":
		cat.Dust = number.Decimal(data.Value).Truncate(8)
	case "dunk":
		cat.Dunk = number.Decimal(data.Value).Truncate(8)
	case "duty":
		cat.Duty = number.Decimal(data.Value).Truncate(8)
	case "chop":
		cat.Chop = number.Decimal(data.Value).Truncate(8)
	case "mat":
		cat.Mat = number.Decimal(data.Value).Truncate(8)
	case "live":
		live, _ := strconv.ParseBool(data.Value)
		cat.Live = types.BitBool(live)
	}

}
