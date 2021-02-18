package Collateral

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Collateral{})

		if err := tx.AutoMigrate(core.Collateral{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_collaterals_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.CollateralStore {
	return &collateralStore{db: db}
}

type collateralStore struct {
	db *db.DB
}

func (s *collateralStore) Create(ctx context.Context, collateral *core.Collateral) error {
	return s.db.Update().Where("trace_id = ?", collateral.TraceID).FirstOrCreate(collateral).Error
}

func toUpdateParams(collateral *core.Collateral) map[string]interface{} {
	return map[string]interface{}{
		"art":   collateral.Art,
		"rate":  collateral.Rate,
		"rho":   collateral.Rho,
		"debt":  collateral.Debt,
		"line":  collateral.Line,
		"dust":  collateral.Dust,
		"price": collateral.Price,
		"mat":   collateral.Mat,
		"duty":  collateral.Duty,
		"chop":  collateral.Chop,
		"dunk":  collateral.Dunk,
		"live":  collateral.Live,
	}
}

func (s *collateralStore) Update(ctx context.Context, collateral *core.Collateral, version int64) error {
	if collateral.Version >= version {
		return nil
	}

	updates := toUpdateParams(collateral)
	updates["version"] = version

	tx := s.db.Update().Model(collateral).Where("version = ?", collateral.Version).Updates(updates)
	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *collateralStore) Find(ctx context.Context, traceID string) (*core.Collateral, error) {
	var collateral core.Collateral
	if err := s.db.View().Where("trace_id = ?", traceID).Take(&collateral).Error; err != nil {
		return nil, err
	}

	return &collateral, nil
}

func (s *collateralStore) List(ctx context.Context) ([]*core.Collateral, error) {
	var collaterals []*core.Collateral

	if err := s.db.View().Find(&collaterals).Error; err != nil {
		return nil, err
	}

	return collaterals, nil
}
