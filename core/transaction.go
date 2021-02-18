package core

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx/types"
)

const (
	TransactionStatusOK int64 = 1
)

type (
	Transaction struct {
		ID        int64          `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time      `json:"created_at,omitempty"`
		TraceID   string         `sql:"size:36" json:"trace_id,omitempty"`
		FollowID  string         `sql:"size:36" json:"follow_id,omitempty"`
		UserID    string         `sql:"size:36" json:"user_id,omitempty"`
		TargetID  string         `json:"target_id,omitempty"`
		Action    int            `json:"action,omitempty"`
		Status    int64          `json:"status,omitempty"`
		Data      types.JSONText `sql:"type:varchar(2048)" json:"data,omitempty"`
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, trace string) (*Transaction, error)
		FindFollow(ctx context.Context, userID, followID string) (*Transaction, error)
		// ListTarget list transactions with given target_id (order by id desc)
		ListTarget(ctx context.Context, targetID string, from int64, limit int) ([]*Transaction, error)
	}
)
