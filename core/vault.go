package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Vault struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		TraceID   string    `sql:"size:36" json:"trace_id,omitempty"`
		Version   int64     `json:"version,omitempty"`
		UserID    string    `sql:"size:36" json:"user_id,omitempty"`
		// CollateralID represent collateral id
		CollateralID string `sql:"size:36" json:"collateral_id,omitempty"`
		// Locked Collateral
		Ink decimal.Decimal `sql:"type:decimal(64,8)" json:"ink,omitempty"`
		// Normalised Debt
		Art decimal.Decimal `sql:"type:decimal(64,16)" json:"art,omitempty"`
	}

	VaultEvent struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		VaultID   string    `sql:"size:36" json:"vault_id,omitempty"`
		Version   int64     `json:"version,omitempty"`
		Action    Action    `json:"action,omitempty"`
		// Locked Collateral change
		Dink decimal.Decimal `sql:"type:decimal(64,8)" json:"dink,omitempty"`
		// Normalised Debt change
		Dart decimal.Decimal `sql:"type:decimal(64,16)" json:"dart,omitempty"`
		// Debt change
		Debt decimal.Decimal `sql:"type:decimal(64,8)" json:"debt,omitempty"`
	}

	VaultStore interface {
		Create(ctx context.Context, vault *Vault) error
		Update(ctx context.Context, vault *Vault, version int64) error
		Find(ctx context.Context, traceID string) (*Vault, error)
		ListUser(ctx context.Context, userID string) ([]*Vault, error)
		// Events
		CreateEvent(ctx context.Context, event *VaultEvent) error
		FindEvent(ctx context.Context, vaultID string, version int64) (*VaultEvent, error)
		ListEvents(ctx context.Context) ([]*VaultEvent, error)
	}
)
