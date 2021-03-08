package asset

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Asset{})

		if err := tx.AutoMigrate(core.Asset{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.AssetStore {
	return &assetStore{db: db}
}

type assetStore struct {
	db *db.DB
}

func (s *assetStore) Create(ctx context.Context, asset *core.Asset) error {
	return s.db.Update().Where("id = ?", asset.ID).FirstOrCreate(asset).Error
}

func toUpdateParams(asset *core.Asset) map[string]interface{} {
	return map[string]interface{}{
		"logo":  asset.Logo,
		"price": asset.Price,
	}
}

func (s *assetStore) Update(ctx context.Context, asset *core.Asset) error {
	updates := toUpdateParams(asset)
	updates["version"] = asset.Version + 1

	tx := s.db.Update().Model(asset).Where("version = ?", asset.Version).Updates(updates)
	if err := tx.Error; err != nil {
		return err
	}

	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *assetStore) Find(ctx context.Context, id string) (*core.Asset, error) {
	var asset core.Asset
	if err := s.db.View().Where("id = ?", id).Take(&asset).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *assetStore) List(ctx context.Context) ([]*core.Asset, error) {
	var assets []*core.Asset
	if err := s.db.View().Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}
