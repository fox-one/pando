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
		// Action the proposal applied
		Action Action `json:"action,omitempty"`
		// Data action parameters
		Data string `sql:"size:256" json:"data,omitempty"`
		// Votes mtg member voted for this proposal
		Votes pq.StringArray `sql:"type:varchar(1024)" json:"votes,omitempty"`
	}

	// ProposalStore define operations for working with proposals on db.
	ProposalStore interface {
		Create(ctx context.Context, proposal *Proposal) error
		Find(ctx context.Context, trace string) (*Proposal, error)
		Update(ctx context.Context, proposal *Proposal, version int64) error
		List(ctx context.Context, fromID int64, limit int) ([]*Proposal, error)
	}

	// Parliament is a proposal version notifier to mtg member admins
	Parliament interface {
		// Created Called when a new proposal created
		Created(ctx context.Context, proposal *Proposal) error
		// Approved called when a proposal has a new vote
		Approved(ctx context.Context, proposal *Proposal) error
		// Passed called when a proposal is passed
		Passed(ctx context.Context, proposal *Proposal) error
	}
)
