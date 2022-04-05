package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Stat struct {
		ID           int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt    time.Time       `json:"created_at,omitempty"`
		UpdatedAt    time.Time       `json:"updated_at,omitempty"`
		Version      int64           `json:"version,omitempty"`
		CollateralID string          `sql:"size:36;NOT NULL" json:"collateral_id,omitempty"`
		Date         time.Time       `sql:"type:date" json:"date,omitempty"`
		Gem          string          `sql:"size:36;NOT NULL" json:"gem,omitempty"`
		Dai          string          `sql:"size:36;NOT NULL" json:"dai,omitempty"`
		Ink          decimal.Decimal `sql:"type:decimal(32,8)" json:"ink,omitempty"`
		Debt         decimal.Decimal `sql:"type:decimal(32,8)" json:"debt,omitempty"`
		InkPrice     decimal.Decimal `sql:"type:decimal(32,8)" json:"ink_price,omitempty"`
		DebtPrice    decimal.Decimal `sql:"type:decimal(32,8)" json:"debt_price,omitempty"`
	}

	AggregatedStat struct {
		Date     time.Time
		GemValue decimal.Decimal
		DaiValue decimal.Decimal
	}

	StatStore interface {
		Save(ctx context.Context, stat *Stat) error
		Find(ctx context.Context, collateralID string, date time.Time) (*Stat, error)
		List(ctx context.Context, collateralID string, from, to time.Time) ([]Stat, error)
		Aggregate(ctx context.Context, from, to time.Time) ([]AggregatedStat, error)
	}
)
