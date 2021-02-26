package message

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Message{})

		if err := tx.AutoMigrate(core.Message{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.MessageStore {
	return &messageStore{db: db}
}

type messageStore struct {
	db *db.DB
}

func (s *messageStore) Create(ctx context.Context, messages []*core.Message) error {
	return s.db.Tx(func(tx *db.DB) error {
		for _, msg := range messages {
			if err := tx.Update().Create(msg).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *messageStore) List(ctx context.Context, limit int) ([]*core.Message, error) {
	var messages []*core.Message
	if err := s.db.View().Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *messageStore) Delete(ctx context.Context, messages []*core.Message) error {
	ids := make([]int64, len(messages))
	for idx, msg := range messages {
		ids[idx] = msg.ID
	}

	if len(ids) == 0 {
		return nil
	}

	return s.db.Update().Where("id IN (?)", ids).Delete(core.Message{}).Error
}
