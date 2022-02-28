package vault

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
	"github.com/shopspring/decimal"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Vault{})

		if err := tx.AutoMigrate(core.Vault{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_vaults_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_vaults_user", "user_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_vaults_collateral_rate", "collateral_id", "rate").Error; err != nil {
			return err
		}

		return nil
	})

	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.VaultEvent{})

		if err := tx.AutoMigrate(core.VaultEvent{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_vault_events_vault_version", "vault_id", "version").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.VaultStore {
	return &vaultStore{db: db}
}

type vaultStore struct {
	db *db.DB
}

func (s *vaultStore) Create(ctx context.Context, vault *core.Vault) error {
	if vault.Ink.IsPositive() {
		vault.Rate = vault.Art.Div(vault.Ink)
	}

	if err := s.db.Update().Where("trace_id = ?", vault.TraceID).FirstOrCreate(vault).Error; err != nil {
		return err
	}

	return nil
}

func toUpdateParams(vault *core.Vault) map[string]interface{} {
	rate := decimal.Zero
	if vault.Ink.IsPositive() {
		rate = vault.Art.Div(vault.Ink)
	}

	return map[string]interface{}{
		"ink":  vault.Ink,
		"art":  vault.Art,
		"rate": rate,
	}
}

func (s *vaultStore) Update(ctx context.Context, vault *core.Vault, version int64) error {
	if vault.Version >= version {
		return nil
	}

	updates := toUpdateParams(vault)
	updates["version"] = version

	tx := s.db.Update().Model(vault).Where("version = ?", vault.Version).Updates(updates)
	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *vaultStore) Find(ctx context.Context, traceID string) (*core.Vault, error) {
	vault := core.Vault{TraceID: traceID}

	if err := s.db.View().Where("trace_id = ?", traceID).Take(&vault).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &vault, nil
		}

		return nil, err
	}

	return &vault, nil
}

func (s *vaultStore) List(ctx context.Context, req core.ListVaultRequest) ([]*core.Vault, error) {
	var vaults []*core.Vault

	tx := s.db.View()
	if req.CollateralID != "" {
		tx = tx.Where("collateral_id = ?", req.CollateralID)

		if req.Rate.IsPositive() {
			tx = tx.Where("rate > ?", req.Rate)
		}
	}

	if req.UserID != "" {
		tx = tx.Where("user_id = ?", req.UserID)
	}

	if req.Desc {
		if req.FromID > 0 {
			tx = tx.Where("id < ?", req.FromID)
		}

		tx = tx.Order("id DESC")
	} else {
		if req.FromID > 0 {
			tx = tx.Where("id > ?", req.FromID)
		}

		tx = tx.Order("id")
	}

	if req.Limit > 0 {
		tx = tx.Limit(req.Limit)
	}

	if err := tx.Find(&vaults).Error; err != nil {
		return nil, err
	}

	return vaults, nil
}

func (s *vaultStore) PluckUser(ctx context.Context, userID string) ([]string, error) {
	var ids []string
	if err := s.db.View().Model(core.Vault{}).
		Select("trace_id").
		Where("user_id = ?", userID).
		Pluck("trace_id", ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *vaultStore) CountCollateral(ctx context.Context) (map[string]int64, error) {
	tx := s.db.View().Model(core.Vault{}).
		Select("collateral_id,Count(id) as count").
		Group("collateral_id")

	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		results = make(map[string]int64)
		row     struct {
			CollateralID string
			Count        int64
		}
	)

	for rows.Next() {
		if err := tx.ScanRows(rows, &row); err != nil {
			return nil, err
		}

		results[row.CollateralID] = row.Count
	}

	return results, nil
}

func (s *vaultStore) CreateEvent(ctx context.Context, event *core.VaultEvent) error {
	if err := s.db.Update().Where("vault_id = ? AND version = ?", event.VaultID, event.Version).FirstOrCreate(event).Error; err != nil {
		return err
	}

	return nil
}

func (s *vaultStore) FindEvent(ctx context.Context, vaultID string, version int64) (*core.VaultEvent, error) {
	event := core.VaultEvent{VaultID: vaultID, Version: version}
	if err := s.db.View().Where("vault_id = ? AND version = ?", vaultID, version).Take(&event).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &event, nil
		}

		return nil, err
	}

	return &event, nil
}

func (s *vaultStore) ListEvents(ctx context.Context, vaultID string) ([]*core.VaultEvent, error) {
	var events []*core.VaultEvent
	if err := s.db.View().Where("vault_id = ?", vaultID).Order("version DESC").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
