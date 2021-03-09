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

		if err := tx.AddUniqueIndex("idx_transactions_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_transactions_user_follow", "user_id", "follow_id").Error; err != nil {
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
	tx := core.Transaction{TraceID: trace}

	if err := s.db.View().Where("trace_id = ?", trace).Take(&tx).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &tx, nil
		}

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

func (s *transactionStore) List(ctx context.Context, fromID int64, limit int) ([]*core.Transaction, error) {
	var transactions []*core.Transaction
	if err := s.db.View().Where("id > ?", fromID).Limit(limit).Order("id").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
