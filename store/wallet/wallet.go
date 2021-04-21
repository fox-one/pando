package wallet

import (
	"context"
	"sort"
	"sync"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
	"github.com/jinzhu/gorm"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Output{})
		if err := tx.AutoMigrate(core.Output{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_outputs_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_outputs_asset_transfer", "asset_id", "spent_by").Error; err != nil {
			return err
		}

		return nil
	})

	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Transfer{})
		if err := tx.AutoMigrate(core.Transfer{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_transfers_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_transfers_handled_passed", "handled", "passed").Error; err != nil {
			return err
		}

		return nil
	})

	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.RawTransaction{})
		if err := tx.AutoMigrate(core.RawTransaction{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_raw_transactions_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.WalletStore {
	return &walletStore{db: db}
}

type walletStore struct {
	db   *db.DB
	once sync.Once
}

func save(db *db.DB, output *core.Output, ack bool) error {
	tx := db.Update().Model(output).Where("trace_id = ?", output.TraceID).Updates(map[string]interface{}{
		"state":     output.State,
		"signed_tx": output.SignedTx,
		"version":   gorm.Expr("version + 1"),
	})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		if ack {
			return db.Update().Create(output).Error
		}

		return saveRawOutput(db, output)
	}

	return nil
}

func (s *walletStore) Save(ctx context.Context, outputs []*core.Output, end bool) error {
	s.once.Do(func() {
		_ = s.runSync(ctx)
	})

	return s.db.Tx(func(tx *db.DB) error {
		for _, utxo := range outputs {
			if err := save(tx, utxo, false); err != nil {
				return err
			}
		}

		if end {
			return ackRawOutputs(tx)
		}

		return nil
	})
}

func (s *walletStore) List(_ context.Context, fromID int64, limit int) ([]*core.Output, error) {
	var outputs []*core.Output
	if err := s.db.View().
		Where("id > ?", fromID).
		Limit(limit).
		Order("id").
		Find(&outputs).Error; err != nil {
		return nil, err
	}

	return outputs, nil
}

func (s *walletStore) ListSpentBy(ctx context.Context, assetID string, spentBy string, limit int) ([]*core.Output, error) {
	var outputs []*core.Output
	if err := s.db.View().
		Where("asset_id = ? AND spent_by = ?", assetID, spentBy).
		Order("id").
		Limit(limit).
		Find(&outputs).Error; err != nil {
		return nil, err
	}

	return outputs, nil
}

func (s *walletStore) ListUnspent(_ context.Context, assetID string, limit int) ([]*core.Output, error) {
	var outputs []*core.Output
	if err := s.db.View().
		Where("asset_id = ? AND spent_by = ?", assetID, "").
		Limit(limit).
		Order("id").
		Find(&outputs).Error; err != nil {
		return nil, err
	}

	return outputs, nil
}

func afterFindTransfer(transfer *core.Transfer) {
	if transfer.Threshold == 0 {
		transfer.Threshold = uint8(len(transfer.Opponents))
	}
}

func (s *walletStore) CreateTransfers(_ context.Context, transfers []*core.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}

	sort.Slice(transfers, func(i, j int) bool {
		return transfers[i].TraceID < transfers[j].TraceID
	})

	return s.db.Tx(func(tx *db.DB) error {
		for _, transfer := range transfers {
			if err := tx.Update().Where("trace_id = ?", transfer.TraceID).FirstOrCreate(transfer).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func updateTransfer(db *db.DB, transfer *core.Transfer) error {
	return db.Update().Model(transfer).Updates(map[string]interface{}{
		"handled": transfer.Handled,
		"passed":  transfer.Passed,
	}).Error
}

func (s *walletStore) UpdateTransfer(ctx context.Context, transfer *core.Transfer) error {
	return updateTransfer(s.db, transfer)
}

func (s *walletStore) ListPendingTransfers(_ context.Context) ([]*core.Transfer, error) {
	var transfers []*core.Transfer
	if err := s.db.View().
		Where("handled = ?", 0).
		Limit(128).
		Order("id").
		Find(&transfers).Error; err != nil {
		return nil, err
	}

	// filter by asset id
	filter := make(map[string]bool)
	var idx int

	for _, t := range transfers {
		if filter[t.AssetID] {
			continue
		}

		transfers[idx] = t
		filter[t.AssetID] = true
		idx++
	}

	transfers = transfers[:idx]

	for _, t := range transfers {
		afterFindTransfer(t)
	}

	return transfers, nil
}

func (s *walletStore) ListNotPassedTransfers(ctx context.Context) ([]*core.Transfer, error) {
	var transfers []*core.Transfer

	if err := s.db.View().
		Where("handled = ? AND passed = ?", 1, 0).
		Limit(128).
		Order("id").
		Find(&transfers).Error; err != nil {
		return nil, err
	}

	for _, t := range transfers {
		afterFindTransfer(t)
	}

	return transfers, nil
}

func (s *walletStore) Spent(_ context.Context, outputs []*core.Output, transfer *core.Transfer) error {
	return s.db.Tx(func(tx *db.DB) error {
		for _, output := range outputs {
			if err := tx.Update().Model(output).Updates(map[string]interface{}{
				"spent_by": transfer.TraceID,
			}).Error; err != nil {
				return err
			}
		}

		transfer.Handled = true
		if transfer.ID > 0 {
			if err := updateTransfer(tx, transfer); err != nil {
				return err
			}
		} else {
			if err := tx.Update().Create(transfer).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *walletStore) CreateRawTransaction(_ context.Context, tx *core.RawTransaction) error {
	return s.db.Update().Where("trace_id = ?", tx.TraceID).FirstOrCreate(tx).Error
}

func (s *walletStore) ListPendingRawTransactions(_ context.Context, limit int) ([]*core.RawTransaction, error) {
	var txs []*core.RawTransaction
	if err := s.db.View().Limit(limit).Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *walletStore) ExpireRawTransaction(_ context.Context, tx *core.RawTransaction) error {
	return s.db.Update().Model(tx).Where("id = ?", tx.ID).Delete(tx).Error
}
