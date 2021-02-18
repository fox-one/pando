package flip

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Flip{})

		if err := tx.AutoMigrate(core.Flip{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_flips_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.FlipStore {
	return &flipStore{
		db: db,
	}
}

type flipStore struct {
	db *db.DB
}

func (s *flipStore) Create(ctx context.Context, flip *core.Flip) error {
	if err := s.db.Update().Where("trace_id = ?", flip.TraceID).FirstOrCreate(flip).Error; err != nil {
		return err
	}

	return nil
}

func toUpdateParams(flip *core.Flip) map[string]interface{} {
	return map[string]interface{}{
		"action": flip.Action,
		"tic":    flip.Tic,
		// "end":    flip.End,
		"bid": flip.Bid,
		"lot": flip.Lot,
		"guy": flip.Guy,
	}
}

func (s *flipStore) Update(ctx context.Context, flip *core.Flip, version int64) error {
	if flip.Version >= version {
		return nil
	}

	updates := toUpdateParams(flip)
	updates["version"] = version

	tx := s.db.Update().Model(flip).Where("version = ?", flip.Version).Updates(updates)
	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *flipStore) Find(ctx context.Context, traceID string) (*core.Flip, error) {
	var flip core.Flip

	if err := s.db.View().Where("trace_id = ?", traceID).Take(&flip).Error; err != nil {
		return nil, err
	}

	return &flip, nil
}
