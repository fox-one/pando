package core

import (
	"context"
	"time"

	"github.com/fox-one/pando/pkg/number"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type (
	// Oracle represent price information
	Oracle struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		AssetID   string    `sql:"size:36" json:"asset_id,omitempty"`
		Version   int64     `json:"version,omitempty"`
		// Current Price Value
		Current decimal.Decimal `sql:"type:decimal(24,12)" json:"current,omitempty"`
		// Next Price Value
		Next decimal.Decimal `sql:"type:decimal(24,12)" json:"next,omitempty"`
		// Time of last update
		PeekAt time.Time `json:"peek_at,omitempty"`
		// Hop time delay (seconds) between poke calls
		Hop int64 `json:"hop,omitempty"`
		// Threshold represents the number of signatures required at least;
		// don't accept any updates by set to zero
		Threshold int64 `json:"threshold,omitempty"`
		// next price providers
		Governors pq.StringArray `sql:"type:varchar(1024)" json:"governors,omitempty"`
	}

	OracleFeed struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UserID    string    `sql:"size:36" json:"user_id,omitempty"`
		PublicKey string    `sql:"size:256" json:"public_key,omitempty"`
	}

	// OracleStore defines operations for working with oracles on db.
	OracleStore interface {
		Create(ctx context.Context, oracle *Oracle) error
		Find(ctx context.Context, assetID string) (*Oracle, error)
		Update(ctx context.Context, oracle *Oracle, version int64) error
		List(ctx context.Context) ([]*Oracle, error)
		ListCurrent(ctx context.Context) (number.Values, error)
		// Rely approve a new feed
		Rely(ctx context.Context, userID, publicKey string) error
		// Deny remove an existing feed
		Deny(ctx context.Context, userID string) error
		ListFeeds(ctx context.Context) ([]*OracleFeed, error)
	}

	// OracleService define operations to parse new price from oracle service outside
	OracleService interface {
		Parse(ctx context.Context, b []byte) (*Oracle, error)
	}
)

var distantFuture = time.Now().AddDate(100, 0, 0)

func (oracle *Oracle) NextPeekAt() time.Time {
	if oracle.Threshold == 0 {
		return distantFuture
	}

	return oracle.PeekAt.Add(time.Duration(oracle.Hop) * time.Second)
}
