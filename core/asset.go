package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	// Asset represent mixin asset
	Asset struct {
		ID        string          `sql:"size:36;PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		Version   int64           `json:"version,omitempty"`
		Name      string          `sql:"size:64" json:"name,omitempty"`
		Symbol    string          `sql:"size:32" json:"symbol,omitempty"`
		Logo      string          `sql:"size:256" json:"logo,omitempty"`
		ChainID   string          `sql:"size:36" json:"chain_id,omitempty"`
		Price     decimal.Decimal `sql:"type:decimal(24,12)" json:"price,omitempty"`
	}

	// AssetStore defines operations for working with assets on db.
	AssetStore interface {
		Save(ctx context.Context, asset *Asset) error
		Find(ctx context.Context, id string) (*Asset, error)
		List(ctx context.Context) ([]*Asset, error)
	}

	// AssetService provides access to remote mixin assets information
	AssetService interface {
		Find(ctx context.Context, id string) (*Asset, error)
		List(ctx context.Context) ([]*Asset, error)
		ReadPrice(ctx context.Context, assetID string, at time.Time) (decimal.Decimal, error)
	}
)
