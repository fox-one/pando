package transaction

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Transaction{})

		if err := tx.AutoMigrate(core.Transaction{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_transactions_trace").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_transactions_user_follow", "user_id", "follow_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_transactions_target", "target_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.TransactionStore {
	return &transactionStore{db: db}
}

type transactionStore struct {
	db *db.DB
}

func (s *transactionStore) Create(ctx context.Context, tx *core.Transaction) error {
	return s.db.Update().Where("trace_id = ?", tx.TraceID).FirstOrCreate(tx).Error
}

func (s *transactionStore) Find(ctx context.Context, trace string) (*core.Transaction, error) {
	var tx core.Transaction

	if err := s.db.View().Where("trace_id = ?", trace).Take(&tx).Error; err != nil {
		return nil, err
	}

	return &tx, nil
}

func (s *transactionStore) FindFollow(ctx context.Context, userID, followID string) (*core.Transaction, error) {
	var tx core.Transaction

	if err := s.db.View().Where("user_id = ? AND follow_id = ?", userID, followID).Last(&tx).Error; err != nil {
		return nil, err
	}

	return &tx, nil
}

func (s *transactionStore) ListTarget(ctx context.Context, targetID string, from int64, limit int) ([]*core.Transaction, error) {
	tx := s.db.View().Where("target_id = ?", targetID)

	if from > 0 {
		tx = tx.Where("id > ?", from)
	}

	var transactions []*core.Transaction
	if err := tx.Limit(limit).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
