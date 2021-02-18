package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Oracle struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		PeekAt    time.Time       `json:"peek_at,omitempty"`
		TraceID   string          `sql:"size:36" json:"trace_id,omitempty"`
		AssetID   string          `sql:"size:36" json:"asset_id,omitempty"`
		Price     decimal.Decimal `sql:"type:decimal(64,16)" json:"price,omitempty"`
	}

	OracleStore interface {
		Create(ctx context.Context, oracle *Oracle) error
		Find(ctx context.Context, assetID string, peekAt time.Time) (*Oracle, error)
		List(ctx context.Context, assetID string, dur time.Duration) ([]*Oracle, error)
	}

	OracleService interface {
		Parse(b []byte) (*Oracle, error)
	}
)
