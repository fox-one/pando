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
		ID           int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt    time.Time       `json:"created_at,omitempty"`
		TraceID      string          `sql:"size:36" json:"trace_id,omitempty"`
		UserID       string          `sql:"size:36" json:"user_id,omitempty"`
		FollowID     string          `sql:"size:36" json:"follow_id,omitempty"`
		AssetID      string          `sql:"size:36" json:"asset_id,omitempty"`
		Amount       decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Action       Action          `json:"action,omitempty"`
		CollateralID string          `sql:"size:36" json:"collateral_id,omitempty"`
		VaultID      string          `sql:"size:36" json:"vault_id,omitempty"`
		FlipID       string          `sql:"size:36" json:"flip_id,omitempty"`
		Status       int             `json:"status,omitempty"`
		Data         types.JSONText  `sql:"type:varchar(1024)" json:"data,omitempty"`
	}

	ListTransactionReq struct {
		CollateralID string
		VaultID      string
		FlipID       string
		Desc         bool
		FromID       int64
		Limit        int
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, trace string) (*Transaction, error)
		// FindFollow Find the last tx with the userID & followID
		FindFollow(ctx context.Context, userID, followID string) (*Transaction, error)
		// List list transactions
		List(ctx context.Context, req *ListTransactionReq) ([]*Transaction, error)
	}
)

func (tx *Transaction) WithFlip(flip *Flip) *Transaction {
	tx.FlipID = flip.TraceID
	tx.VaultID = flip.VaultID
	tx.CollateralID = flip.CollateralID
	return tx
}

func (tx *Transaction) WithVault(vault *Vault) *Transaction {
	tx.VaultID = vault.TraceID
	tx.CollateralID = vault.CollateralID
	return tx
}

func (tx *Transaction) WithCollateral(cat *Collateral) *Transaction {
	tx.CollateralID = cat.TraceID
	return tx
}

func (tx *Transaction) Write(status int, data interface{}) {
	tx.Status = status

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	tx.Data = b
}
