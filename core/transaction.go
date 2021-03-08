package core

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/shopspring/decimal"
)

type TransactionStatus int

const (
	TransactionStatusPending TransactionStatus = iota
	TransactionStatusAbort
	TransactionStatusOk
)

//go:generate stringer -type TransactionStatus -trimprefix TransactionStatus

type (
	Transaction struct {
		ID         int64             `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt  time.Time         `json:"created_at,omitempty"`
		TraceID    string            `sql:"size:36" json:"trace_id,omitempty"`
		UserID     string            `sql:"size:36" json:"user_id,omitempty"`
		FollowID   string            `sql:"size:36" json:"follow_id,omitempty"`
		AssetID    string            `sql:"size:36" json:"asset_id,omitempty"`
		Amount     decimal.Decimal   `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Action     Action            `json:"action,omitempty"`
		Status     TransactionStatus `json:"status,omitempty"`
		Message    string            `sql:"size:128" json:"message,omitempty"`
		Parameters types.JSONText    `sql:"type:varchar(1024)" json:"parameters,omitempty"` // []interface{}
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, trace string) (*Transaction, error)
		// FindFollow Find the last tx with the userID & followID
		FindFollow(ctx context.Context, userID, followID string) (*Transaction, error)
		// List list transactions
		List(ctx context.Context, fromID int64, limit int) ([]*Transaction, error)
	}
)
