package user

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.User{})

		if err := tx.AutoMigrate(core.User{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_users_mixin_id", "mixin_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.UserStore {
	return &userStore{db: db}
}

type userStore struct {
	db *db.DB
}

func toUpdateParams(user *core.User) map[string]interface{} {
	return map[string]interface{}{
		"name":         user.Name,
		"avatar":       user.Avatar,
		"access_token": user.AccessToken,
		"version":      time.Now().UnixNano(),
	}
}

func update(db *db.DB, user *core.User) (int64, error) {
	updates := toUpdateParams(user)
	tx := db.Update().Model(user).Where("mixin_id = ?", user.MixinID).Updates(updates)
	return tx.RowsAffected, tx.Error
}

func (s *userStore) Save(_ context.Context, user *core.User) error {
	return s.db.Tx(func(tx *db.DB) error {
		rows, err := update(tx, user)
		if err != nil {
			return err
		}

		if rows == 0 {
			return tx.Update().Create(user).Error
		}

		return nil
	})
}

func (s *userStore) Find(_ context.Context, mixinID string) (*core.User, error) {
	var user core.User
	if err := s.db.View().Where("mixin_id = ?", mixinID).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
