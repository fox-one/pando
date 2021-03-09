package core

import (
	"context"
	"time"

	"github.com/fox-one/pando/pkg/number"
	"github.com/shopspring/decimal"
)

type (
	Oracle struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		AssetID   string    `sql:"size:36" json:"asset_id,omitempty"`
		Version   int64     `json:"version,omitempty"`
		// Hop time delay (seconds) between poke calls
		Hop int64 `json:"hop,omitempty"`
		// Current Price Value
		Current decimal.Decimal `sql:"type:decimal(24,12)" json:"current,omitempty"`
		// Next Price Value
		Next decimal.Decimal `sql:"type:decimal(24,12)" json:"next,omitempty"`
		// Time of last update
		PeekAt time.Time `json:"peek_at,omitempty"`
	}

	OracleStore interface {
		Save(ctx context.Context, oracle *Oracle, version int64) error
		Find(ctx context.Context, assetID string) (*Oracle, error)
		List(ctx context.Context) ([]*Oracle, error)
		ListCurrent(ctx context.Context) (number.Values, error)
	}

	OracleService interface {
		Parse(b []byte) ([]byte, error)
	}
)
