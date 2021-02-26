package oracle

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Oracle{})
		if err := tx.AutoMigrate(core.Oracle{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_oracles_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_oracles_asset_peek", "asset_id", "peek_at").Error; err != nil {
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
	return s.db.Update().Where("trace_id = ?", oracle.TraceID).FirstOrCreate(oracle).Error
}

func (s *oracleStore) Find(ctx context.Context, assetID string, peekAt time.Time) (*core.Oracle, error) {
	oracle := core.Oracle{AssetID: assetID}
	if err := s.db.View().Where("asset_id = ? AND peek_at <= ?", assetID, peekAt).Last(&oracle).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &oracle, nil
		}

		return nil, err
	}

	return &oracle, nil
}

func (s *oracleStore) List(ctx context.Context, assetID string, dur time.Duration) ([]*core.Oracle, error) {
	var oracles []*core.Oracle
	if err := s.db.View().Where("asset_id = ? AND peed_at >= ?", assetID, time.Now().Add(-dur)).Find(&oracles).Error; err != nil {
		return nil, err
	}

	return oracles, nil
}
