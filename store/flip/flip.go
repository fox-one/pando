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

	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.FlipEvent{})

		if err := tx.AutoMigrate(core.FlipEvent{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_flip_events_flip_version", "flip_id", "version").Error; err != nil {
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
	flip := core.Flip{TraceID: traceID}

	if err := s.db.View().Where("trace_id = ?", traceID).Take(&flip).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &flip, nil
		}

		return nil, err
	}

	return &flip, nil
}

func (s *flipStore) List(ctx context.Context, from int64, limit int) ([]*core.Flip, error) {
	tx := s.db.View()

	if from > 0 {
		tx = tx.Where("id < ?", from)
	}

	var flips []*core.Flip
	if err := tx.Limit(limit).Order("id DESC").Find(&flips).Error; err != nil {
		return nil, err
	}

	return flips, nil
}

func (s *flipStore) CreateEvent(ctx context.Context, event *core.FlipEvent) error {
	if err := s.db.Update().Where("flip_id = ? AND version = ?", event.FlipID, event.Version).FirstOrCreate(event).Error; err != nil {
		return err
	}

	return nil
}

func (s *flipStore) FindEvent(ctx context.Context, flipID string, version int64) (*core.FlipEvent, error) {
	event := core.FlipEvent{FlipID: flipID, Version: version}
	if err := s.db.View().Where("flip_id = ? AND version = ?", flipID, version).Take(&event).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &event, nil
		}

		return nil, err
	}

	return &event, nil
}

func (s *flipStore) ListEvents(ctx context.Context, flipID string) ([]*core.FlipEvent, error) {
	var events []*core.FlipEvent
	if err := s.db.View().Where("flip_id = ?", flipID).Order("version").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
