package core

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type (
	Proposal struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		PassedAt  sql.NullTime    `json:"passed_at,omitempty"`
		Version   int64           `json:"version,omitempty"`
		TraceID   string          `sql:"size:36" json:"trace_id,omitempty"`
		Creator   string          `sql:"size:36" json:"creator,omitempty"`
		AssetID   string          `sql:"size:36" json:"asset_id,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Action    Action          `json:"action,omitempty"`
		Data      string          `sql:"size:256" json:"data,omitempty"`
		Votes     pq.StringArray  `sql:"type:varchar(1024)" json:"votes,omitempty"`
	}

	ProposalStore interface {
		Create(ctx context.Context, proposal *Proposal) error
		Find(ctx context.Context, trace string) (*Proposal, error)
		Update(ctx context.Context, proposal *Proposal, version int64) error
		List(ctx context.Context, fromID int64, limit int) ([]*Proposal, error)
	}

	Parliament interface {
		Created(ctx context.Context, proposal *Proposal) error
		Approved(ctx context.Context, proposal *Proposal) error
		Passed(ctx context.Context, proposal *Proposal) error
	}
)
