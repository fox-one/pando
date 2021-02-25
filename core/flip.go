package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Flip struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Version   int64     `json:"version,omitempty"`
		TraceID   string    `sql:"size:36" json:"trace_id,omitempty"`
		VaultID   string    `sql:"size:36" json:"vault_id,omitempty"`
		Action    Action    `json:"action,omitempty"`
		// bid expiry time
		Tic time.Time `json:"tic,omitempty"`
		// auction expiry time
		End time.Time `json:"end,omitempty"`
		// pUSD paid
		Bid decimal.Decimal `sql:"type:decimal(64,8)" json:"bid,omitempty"`
		// gems in return for bid
		Lot decimal.Decimal `sql:"type:decimal(64,8)" json:"lot,omitempty"`
		// total pUSD wanted
		Tab decimal.Decimal `sql:"type:decimal(64,8)" json:"tab,omitempty"`
		// high bidder
		Guy string `sql:"size:36" json:"guy,omitempty"`
	}

	FlipStore interface {
		Create(ctx context.Context, flip *Flip) error
		Update(ctx context.Context, flip *Flip, version int64) error
		Find(ctx context.Context, traceID string) (*Flip, error)
	}
)
