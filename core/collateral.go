package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	// Collateral (CAT) represent collateral type
	Collateral struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		TraceID   string    `sql:"size:36" json:"trace_id,omitempty"`
		Version   int64     `json:"version,omitempty"`
		Name      string    `sql:"size:64" json:"name,omitempty"`
		// Gem represent deposit asset id
		Gem string `sql:"size:36" json:"gem,omitempty"`
		// Dai represent debt asset id
		Dai string `sql:"size:36" json:"dai,omitempty"`
		// Ink represent All Locked Collateral
		Ink decimal.Decimal `sql:"type:decimal(64,8)" json:"ink,omitempty"`
		// Total Normalised Debt
		Art decimal.Decimal `sql:"type:decimal(64,16)" json:"art,omitempty"`
		// Accumulated Rates
		Rate decimal.Decimal `sql:"type:decimal(64,16)" json:"rate,omitempty"`
		// Time of last drip [unix epoch time]
		Rho time.Time `json:"rho,omitempty"`
		// Debt the total quantity of dai issued
		Debt decimal.Decimal `sql:"type:decimal(64,8)" json:"debt,omitempty"`
		// the debt ceiling
		Line decimal.Decimal `sql:"type:decimal(64,8)" json:"line,omitempty"`
		// the debt floor, eg 100
		Dust decimal.Decimal `sql:"type:decimal(64,8)" json:"dust,omitempty"`
		// Price = Gem.Price / Dai.Price
		Price decimal.Decimal `sql:"type:decimal(32,12)" json:"price,omitempty"`
		// Liquidation ratio, eg 150%
		Mat decimal.Decimal `sql:"type:decimal(10,8)" json:"mat,omitempty"`
		// stability fee, eg 110%
		Duty decimal.Decimal `sql:"type:decimal(20,8)" json:"duty,omitempty"`
		// Liquidation Penalty, eg 113%
		Chop decimal.Decimal `sql:"type:decimal(10,8)" json:"chop,omitempty"`
		// Dunk, max liquidation Quantity, eg 50000
		Dunk decimal.Decimal `sql:"type:decimal(64,8)" json:"dunk,omitempty"`
		// Flip Options
		// Box, Max Dai out for liquidation
		Box decimal.Decimal `sql:"type:decimal(64,8)" json:"box,omitempty"`
		// Litter, Balance of Dai out for liquidation
		Litter decimal.Decimal `sql:"type:decimal(64,8)" json:"litter,omitempty"`
		// Beg minimum bid increase
		Beg decimal.Decimal `sql:"type:decimal(8,4)" json:"beg,omitempty"`
		// TTL bid duration in seconds
		TTL int64 `json:"ttl,omitempty"`
		// Tau flip duration in seconds
		Tau int64 `json:"tau,omitempty"`
		// Live
		Live int `json:"live,omitempty"`
	}

	// CollateralStore define operations for working with collateral on db
	CollateralStore interface {
		Create(ctx context.Context, collateral *Collateral) error
		Update(ctx context.Context, collateral *Collateral, version int64) error
		Find(ctx context.Context, traceID string) (*Collateral, error)
		List(ctx context.Context) ([]*Collateral, error)
	}
)
