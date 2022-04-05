package stat

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Stat{})
		if err := tx.AutoMigrate(core.Stat{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_stats_cat_date", "collateral_id", "date").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.StatStore {
	return &statStore{
		db: db,
	}
}

type statStore struct {
	db *db.DB
}

func (s *statStore) Save(ctx context.Context, stat *core.Stat) error {
	updates := map[string]interface{}{
		"ink":       stat.Ink,
		"debt":      stat.Debt,
		"gem_price": stat.GemPrice,
		"dai_price": stat.DaiPrice,
		"version":   stat.Version,
	}

	tx := s.db.Update().Model(stat).
		Where("collateral_id = ? AND date = ?", stat.CollateralID, stat.Date).
		Updates(updates)

	if err := tx.Error; err != nil {
		return err
	}

	if tx.RowsAffected == 0 {
		stat.Version = 1
		return tx.Create(stat).Error
	}

	return nil
}

func (s *statStore) Find(ctx context.Context, collateralID string, date time.Time) (*core.Stat, error) {
	var stats []*core.Stat

	if err := s.db.View().
		Where("collateral_id = ? AND date <= ?", collateralID, date).
		Order("date DESC").
		Limit(1).
		Find(&stats).Error; err != nil {
		return nil, err
	}

	if len(stats) == 0 {
		return &core.Stat{
			CollateralID: collateralID,
			Date:         date,
		}, nil
	}

	return stats[0], nil
}

func (s *statStore) List(ctx context.Context, collateralID string, from, to time.Time) ([]*core.Stat, error) {
	var stats []*core.Stat

	if err := s.db.View().
		Where("collateral_id = ? AND date >= ? AND date <= ?", collateralID, from, to).
		Order("date ASC").
		Find(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *statStore) Aggregate(ctx context.Context, from, to time.Time) ([]core.AggregatedStat, error) {
	rows, err := s.db.View().Model(core.Stat{}).Where("date >= ? AND date <= ?", from, to).
		Select("date, SUM(ink * gem_price) AS gem_value, SUM(debt * dai_price) AS dai_value").
		Group("date").Rows()

	if err != nil {
		return nil, err
	}

	var stats []core.AggregatedStat

	for rows.Next() {
		var stat core.AggregatedStat
		if err := rows.Scan(&stat.Date, &stat.GemValue, &stat.DaiValue); err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}

	return stats, nil
}
