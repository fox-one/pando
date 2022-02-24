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

		if err := tx.AddIndex("idx_flip_events_guy", "guy").Error; err != nil {
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

func (s *flipStore) ListParticipates(ctx context.Context, userID string) ([]string, error) {
	var ids []string
	if err := s.db.View().Model(core.FlipEvent{}).
		Select("DISTINCT(flip_id)").
		Where("guy = ?", userID).
		Pluck("flip_id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *flipStore) QueryFlips(ctx context.Context, query core.FlipQuery) ([]*core.Flip, int64, error) {
	db := s.db.View()
	tx := db.Model(core.Flip{}).Order("id DESC")

	if query.Phase > 0 {
		if query.Phase == core.FlipPhaseDeal {
			tx = tx.Where("action = ?", core.ActionFlipDeal)
		} else {
			tx = tx.Where("action != ?", core.ActionFlipDeal)

			if query.Phase == core.FlipPhaseTend {
				tx = tx.Where("bid < tab")
			} else if query.Phase == core.FlipPhaseDent {
				tx = tx.Where("bid = tab")
			}
		}
	}

	if query.VaultUserID != "" {
		sub := db.Model(core.Vault{}).Select("trace_id").Where("user_id = ?", query.VaultUserID).QueryExpr()
		tx = tx.Where("vault_id IN (?)", sub)
	}

	if query.Participator != "" {
		sub := db.Model(core.FlipEvent{}).Select("DISTINCT(flip_id)").Where("guy = ?", query.Participator).QueryExpr()
		tx = tx.Where("trace_id IN (?)", sub)
	}

	var (
		flips []*core.Flip
		total int64
	)

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total > query.Offset {
		if err := tx.Offset(query.Offset).Limit(query.Limit).Find(&flips).Error; err != nil {
			return nil, 0, err
		}
	}

	return flips, total, nil
}
