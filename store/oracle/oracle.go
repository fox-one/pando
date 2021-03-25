package oracle

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Oracle{})
		if err := tx.AutoMigrate(core.Oracle{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_oracles_asset", "asset_id").Error; err != nil {
			return err
		}

		return nil
	})

	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.OracleFeed{})
		if err := tx.AutoMigrate(core.OracleFeed{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_oracle_feeds_user", "user_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.OracleStore {
	return &oracleStore{db: db}
}

type oracleStore struct {
	db *db.DB
}

func (s *oracleStore) Create(ctx context.Context, oracle *core.Oracle) error {
	return s.db.Update().Where("asset_id = ?", oracle.AssetID).FirstOrCreate(oracle).Error
}

func (s *oracleStore) Update(ctx context.Context, oracle *core.Oracle, version int64) error {
	if oracle.Version >= version {
		return nil
	}

	updates := map[string]interface{}{
		"hop":       oracle.Hop,
		"current":   oracle.Current,
		"next":      oracle.Next,
		"peek_at":   oracle.PeekAt,
		"threshold": oracle.Threshold,
		"governors": oracle.Governors,
		"version":   version,
	}

	tx := s.db.Update().Model(oracle).Where("version = ?", oracle.Version).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *oracleStore) Find(ctx context.Context, assetID string) (*core.Oracle, error) {
	oracle := core.Oracle{AssetID: assetID}
	if err := s.db.View().Where("asset_id = ?", assetID).Take(&oracle).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &oracle, nil
		}

		return nil, err
	}

	return &oracle, nil
}

func (s *oracleStore) List(ctx context.Context) ([]*core.Oracle, error) {
	var oracles []*core.Oracle
	if err := s.db.View().Find(&oracles).Error; err != nil {
		return nil, err
	}

	return oracles, nil
}

func (s *oracleStore) ListCurrent(ctx context.Context) (number.Values, error) {
	var oracles []*core.Oracle
	if err := s.db.View().Select("asset_id, current").Find(&oracles).Error; err != nil {
		return nil, err
	}

	prices := make(number.Values, len(oracles))
	for _, oracle := range oracles {
		if oracle.Current.IsPositive() {
			prices.Set(oracle.AssetID, oracle.Current)
		}
	}

	return prices, nil
}

func (s *oracleStore) Rely(ctx context.Context, userID, publicKey string) error {
	feed := &core.OracleFeed{
		UserID:    userID,
		PublicKey: publicKey,
	}

	return s.db.Update().Where("user_id = ?", userID).FirstOrCreate(feed).Error
}

func (s *oracleStore) Deny(ctx context.Context, userID string) error {
	return s.db.Update().Where("user_id = ?", userID).Delete(core.OracleFeed{}).Error
}

func (s *oracleStore) ListFeeds(ctx context.Context) ([]*core.OracleFeed, error) {
	var feeds []*core.OracleFeed
	if err := s.db.View().Find(&feeds).Error; err != nil {
		return nil, err
	}

	return feeds, nil
}
