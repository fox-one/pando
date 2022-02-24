package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type FlipPhase int

const (
	_ FlipPhase = iota
	FlipPhaseTend
	FlipPhaseDent
	FlipPhaseBid
	FlipPhaseDeal
)

//go:generate stringer --type FlipPhase --trimprefix FlipPhase

type (
	// Flip represent auction by kicking unsafe vaults
	Flip struct {
		ID           int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
		Version      int64     `json:"version,omitempty"`
		TraceID      string    `sql:"size:36" json:"trace_id,omitempty"`
		CollateralID string    `sql:"size:36" json:"collateral_id,omitempty"`
		VaultID      string    `sql:"size:36" json:"vault_id,omitempty"`
		Action       Action    `json:"action,omitempty"`
		// bid expiry time
		Tic int64 `json:"tic,omitempty"`
		// auction expiry time
		End int64 `json:"end,omitempty"`
		// pUSD paid
		Bid decimal.Decimal `sql:"type:decimal(64,8)" json:"bid,omitempty"`
		// gems in return for bid
		Lot decimal.Decimal `sql:"type:decimal(64,8)" json:"lot,omitempty"`
		// total pUSD wanted
		Tab decimal.Decimal `sql:"type:decimal(64,8)" json:"tab,omitempty"`
		// Art
		Art decimal.Decimal `sql:"type:decimal(64,16)" json:"art,omitempty"`
		// high bidder
		Guy string `sql:"size:36" json:"guy,omitempty"`
	}

	// FlipEvent define operation history on flip
	FlipEvent struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		FlipID    string          `sql:"size:36" json:"flip_id,omitempty"`
		Version   int64           `json:"version,omitempty"`
		Action    Action          `json:"action,omitempty"`
		Bid       decimal.Decimal `sql:"type:decimal(64,8)" json:"bid,omitempty"`
		Lot       decimal.Decimal `sql:"type:decimal(64,8)" json:"lot,omitempty"`
		Guy       string          `sql:"size:36" json:"guy,omitempty"`
	}

	FlipQuery struct {
		Phase         FlipPhase
		VaultUserID   string
		Participator  string
		Offset, Limit int64
	}

	// FlipStore define operations for working on Flip & FlipEvent on db
	FlipStore interface {
		Create(ctx context.Context, flip *Flip) error
		Update(ctx context.Context, flip *Flip, version int64) error
		Find(ctx context.Context, traceID string) (*Flip, error)
		List(ctx context.Context, from int64, limit int) ([]*Flip, error)
		// Event
		CreateEvent(ctx context.Context, event *FlipEvent) error
		FindEvent(ctx context.Context, flipID string, version int64) (*FlipEvent, error)
		ListEvents(ctx context.Context, flipID string) ([]*FlipEvent, error)

		ListParticipates(ctx context.Context, userID string) ([]string, error)
		QueryFlips(ctx context.Context, query FlipQuery) ([]*Flip, int64, error)
	}
)

func (f *Flip) TicFinished(at time.Time) bool {
	return f.Tic > 0 && at.Unix() > f.Tic
}

func (f *Flip) EndFinished(at time.Time) bool {
	return at.Unix() > f.End
}
