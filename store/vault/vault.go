package vault

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
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
	if err := s.db.Update().Where("trace_id = ?", vault.TraceID).FirstOrCreate(vault).Error; err != nil {
		return err
	}

	return nil
}

func toUpdateParams(vault *core.Vault) map[string]interface{} {
	return map[string]interface{}{
		"ink": vault.Ink,
		"art": vault.Art,
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

func (s *vaultStore) ListUser(ctx context.Context, userID string) ([]*core.Vault, error) {
	var vaults []*core.Vault

	if err := s.db.View().Where("user_id = ?", userID).Find(&vaults).Error; err != nil {
		return nil, err
	}

	return vaults, nil
}
