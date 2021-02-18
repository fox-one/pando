package core

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type ProposalAction int

const (
	_ ProposalAction = iota
	ProposalActionVote
	ProposalActionWithdraw
	ProposalActionAddCollateral
	ProposalActionVatGov
	ProposalActionFlipGov
)

//go:generate stringer -type ProposalAction -trimprefix ProposalAction

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
		Action    ProposalAction  `json:"action,omitempty"`
		Content   types.JSONText  `sql:"type:varchar(1024)" json:"content,omitempty"`
		Votes     pq.StringArray  `sql:"type:varchar(1024)" json:"votes,omitempty"`
	}

	ProposalStore interface {
		Create(ctx context.Context, proposal *Proposal) error
		Find(ctx context.Context, trace string) (*Proposal, error)
		Update(ctx context.Context, proposal *Proposal) error
		List(ctx context.Context, fromID int64, limit int) ([]*Proposal, error)
	}

	Parliament interface {
		ProposalCreated(ctx context.Context, proposal *Proposal, by *Member) error
		ProposalApproved(ctx context.Context, proposal *Proposal, by *Member) error
		ProposalPassed(ctx context.Context, proposal *Proposal) error
	}
)
