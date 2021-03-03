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

		if err := tx.AddIndex("idx_transactions_cat_vat_flip", "collateral_id", "vault_id", "flip_id").Error; err != nil {
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

func (s *transactionStore) List(ctx context.Context, req *core.ListTransactionReq) ([]*core.Transaction, error) {
	tx := s.db.View()

	limit := 100
	if req != nil {
		if req.CollateralID != "" {
			tx = tx.Where("collateral_id = ?", req.CollateralID)

			if req.VaultID != "" {
				tx = tx.Where("vault_id = ?", req.VaultID)

				if req.FlipID != "" {
					tx = tx.Where("flip_id = ?", req.FlipID)
				}
			}
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
		}

		limit = req.Limit
	}

	var transactions []*core.Transaction
	if err := tx.Limit(limit).Order("id").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
