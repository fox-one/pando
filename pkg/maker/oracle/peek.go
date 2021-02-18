package oracle

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/shopspring/decimal"
)

type Data struct {
	PeekAt  time.Time       `json:"peek_at,omitempty"`
	AssetID string          `json:"asset_id,omitempty"`
	Price   decimal.Decimal `json:"price,omitempty"`
}

type Option struct {
	Overdue time.Duration `json:"overdue,omitempty"`
}

func Peek(tx *maker.Tx, opt Option, data *Data) error {
	if tx.Now.Sub(data.PeekAt) > opt.Overdue {
		return ErrOracleOverdue
	}

	return nil
}

func ApplyPeek(tx *maker.Tx, data *Data) *core.Oracle {
	return &core.Oracle{
		CreatedAt: tx.Now,
		PeekAt:    data.PeekAt,
		TraceID:   tx.TraceID,
		AssetID:   data.AssetID,
		Price:     data.Price,
	}
}
