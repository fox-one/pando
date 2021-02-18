package core

import (
	"context"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type (
	// Output represent Mixin Network multisig Outputs
	Output struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		Version   int64           `sql:"NOT NULL" json:"version,omitempty"`
		TraceID   string          `sql:"type:char(36)" json:"trace_id,omitempty"`
		AssetID   string          `sql:"type:char(36)" json:"asset_id,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Memo      string          `sql:"size:200" json:"memo,omitempty"`
		State     string          `sql:"size:24" json:"state,omitempty"`

		// SpentBy represent the associated transfer trace id
		SpentBy string `sql:"type:char(36);NOT NULL" json:"spent_by,omitempty"`

		// UTXO json Data
		Data types.JSONText `sql:"type:TEXT" json:"data,omitempty"`

		// Raw Mixin UTXO
		UTXO *mixin.MultisigUTXO `sql:"-" json:"-,omitempty"`
	}

	Transfer struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		TraceID   string          `sql:"type:char(36)" json:"trace_id,omitempty"`
		AssetID   string          `sql:"type:char(36)" json:"asset_id,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Memo      string          `sql:"size:200" json:"memo,omitempty"`
		Handled   types.BitBool   `sql:"type:bit(1)" json:"handled,omitempty"`
		Passed    types.BitBool   `sql:"type:bit(1)" json:"passed,omitempty"`
		Threshold uint8           `json:"threshold,omitempty"`
		Opponents pq.StringArray  `sql:"type:varchar(1024)" json:"opponents,omitempty"`
	}

	RawTransaction struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		TraceID   string    `sql:"type:char(36);" json:"trace_id,omitempty"`
		Data      string    `sql:"type:TEXT" json:"data,omitempty"`
	}

	WalletStore interface {
		// Save batch update multiple Output
		Save(ctx context.Context, outputs []*Output) error
		// List return a list of Output by order
		List(ctx context.Context, fromID int64, limit int) ([]*Output, error)
		// ListUnspent list unspent Output
		ListUnspent(ctx context.Context, assetID string, limit int) ([]*Output, error)
		ListSpentBy(ctx context.Context, assetID string, spentBy string) ([]*Output, error)
		// Transfers
		CreateTransfers(ctx context.Context, transfers []*Transfer) error
		UpdateTransfer(ctx context.Context, transfer *Transfer) error
		ListPendingTransfers(ctx context.Context) ([]*Transfer, error)
		ListNotPassedTransfers(ctx context.Context) ([]*Transfer, error)
		Spent(ctx context.Context, outputs []*Output, transfer *Transfer) error
		// mixin net transaction
		CreateRawTransaction(ctx context.Context, tx *RawTransaction) error
		ListPendingRawTransactions(ctx context.Context, limit int) ([]*RawTransaction, error)
		ExpireRawTransaction(ctx context.Context, tx *RawTransaction) error
	}

	WalletService interface {
		// Pull fetch NEW Output updates
		Pull(ctx context.Context, offset time.Time, limit int) ([]*Output, error)
		// Spend spend multiple Output
		Spend(ctx context.Context, outputs []*Output, transfer *Transfer) (*RawTransaction, error)
		// ReqTransfer generate payment code for multisig transfer
		ReqTransfer(ctx context.Context, transfer *Transfer) (string, error)
	}
)
