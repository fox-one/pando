package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/shopspring/decimal"
)

const (
	TxStatusSuccess = 1
	TxStatusFailed  = 2
)

type (
	Transaction struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		TraceID   string          `sql:"size:36" json:"trace_id,omitempty"`
		UserID    string          `sql:"size:36" json:"user_id,omitempty"`
		FollowID  string          `sql:"size:36" json:"follow_id,omitempty"`
		TargetID  string          `sql:"size:36" json:"target_id,omitempty"`
		AssetID   string          `sql:"size:36" json:"asset_id,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Action    Action          `json:"action,omitempty"`
		Status    int             `json:"status,omitempty"`
		Data      types.JSONText  `sql:"type:varchar(1024)" json:"data,omitempty"`
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, trace string) (*Transaction, error)
		FindFollow(ctx context.Context, userID, followID string) (*Transaction, error)
		// ListTarget list transactions with given target_id (order by id desc)
		ListTarget(ctx context.Context, targetID string, from int64, limit int) ([]*Transaction, error)
	}
)

func (tx Transaction) Write(status int, data interface{}) {
	tx.Status = status

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	tx.Data = b
}
