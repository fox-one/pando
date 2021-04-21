package core

import (
	"context"
	"sort"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type (
	// Output represent Mixin Network multisig Outputs
	Output struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Version   int64     `sql:"NOT NULL" json:"version,omitempty"`
		TraceID   string    `sql:"type:char(36)" json:"trace_id,omitempty"`
		// mixin id of operator
		Sender          string          `sql:"type:char(36)" json:"sender,omitempty"`
		AssetID         string          `sql:"type:char(36)" json:"asset_id,omitempty"`
		Amount          decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Memo            string          `sql:"size:320" json:"memo,omitempty"`
		State           string          `sql:"size:24" json:"state,omitempty"`       // unspent,signed,spent
		TransactionHash string          `sql:"size:64" json:"hash,omitempty"`        // utxo.transaction_hash.hex
		OutputIndex     int             `json:"output_index,omitempty"`              // utxo.output_index
		SignedTx        string          `sql:"type:TEXT" json:"signed_tx,omitempty"` // utxo.signed_tx

		// SpentBy represent the associated transfer trace id
		SpentBy string `sql:"type:char(36);NOT NULL" json:"spent_by,omitempty"`
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
		Save(ctx context.Context, outputs []*Output, end bool) error
		// List return a list of Output by order
		List(ctx context.Context, fromID int64, limit int) ([]*Output, error)
		// ListUnspent list unspent Output
		ListUnspent(ctx context.Context, assetID string, limit int) ([]*Output, error)
		ListSpentBy(ctx context.Context, assetID string, spentBy string, limit int) ([]*Output, error)
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
		// HandleTransfer handle a transfer request
		HandleTransfer(ctx context.Context, transfer *Transfer) error
	}
)

func (a *Output) Cmp(b *Output) int {
	if dur := a.CreatedAt.Sub(b.CreatedAt); dur > 0 {
		return 1
	} else if dur < 0 {
		return -1
	}

	if a.TransactionHash > b.TransactionHash {
		return 1
	} else if a.TransactionHash < b.TransactionHash {
		return -1
	}

	if a.OutputIndex > b.OutputIndex {
		return 1
	} else if a.OutputIndex < b.OutputIndex {
		return -1
	}

	return 0
}

func SortOutputs(outputs []*Output) {
	sort.Slice(outputs, func(i, j int) bool {
		a, b := outputs[i], outputs[j]
		return a.Cmp(b) < 0
	})
}
