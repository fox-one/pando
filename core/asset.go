package core

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
)

type (
	Asset struct {
		ID        string          `sql:"size:36;PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		Version   int64           `json:"version,omitempty"`
		Name      string          `sql:"size:64" json:"name,omitempty"`
		Symbol    string          `sql:"size:32" json:"symbol,omitempty" valid:"required"`
		Logo      string          `sql:"size:256" json:"logo,omitempty"`
		ChainID   string          `sql:"size:36" json:"chain_id,omitempty" valid:"uuid,required"`
		Price     decimal.Decimal `sql:"type:decimal(64,20)" json:"price,omitempty"`
	}

	// AssetStore defines operations for working with assets on db.
	AssetStore interface {
		Create(ctx context.Context, asset *Asset) error
		Update(ctx context.Context, asset *Asset) error
		Find(ctx context.Context, id string) (*Asset, error)
		List(ctx context.Context) ([]*Asset, error)
	}

	// AssetService provides access to assets information
	AssetService interface {
		Find(ctx context.Context, id string) (*Asset, error)
	}
)

func (asset *Asset) Validate() error {
	_, err := govalidator.ValidateStruct(asset)
	return err
}
