package wallet

import (
	"encoding/json"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/store/db"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx/types"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(RawOutput{})

		if err := tx.AutoMigrate(RawOutput{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_raw_outputs_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_raw_outputs_created", "created_at").Error; err != nil {
			return err
		}

		return nil
	})
}

type RawOutput struct {
	ID        int64          `sql:"PRIMARY_KEY" json:"id,omitempty"`
	CreatedAt int64          `json:"created_at"`
	TraceID   string         `sql:"size:36" json:"trace_id"`
	Version   int64          `sql:"not null" json:"version"`
	Data      types.JSONText `sql:"type:TEXT" json:"data"`
}

func saveRawOutput(db *db.DB, output *core.Output) error {
	data, _ := json.Marshal(output)

	raw := &RawOutput{
		CreatedAt: output.CreatedAt.UnixNano(),
		TraceID:   output.TraceID,
		Version:   1,
		Data:      data,
	}

	tx := db.Update().Model(raw).
		Where("trace_id = ?", raw.TraceID).
		Updates(map[string]interface{}{
			"data":    raw.Data,
			"version": gorm.Expr("version + 1"),
		})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return db.Update().Create(raw).Error
	}

	return nil
}

func trimSuffix(raws []*RawOutput) []*RawOutput {
	var (
		r = len(raws) - 1
		l = r - 1
	)

	for l >= 0 {
		if raws[l].CreatedAt != raws[r].CreatedAt {
			break
		}

		l = l - 1
	}

	if l >= 0 {
		raws = raws[:l+1]
	}

	return raws
}

func countRawOutputs(db *db.DB) (int64, error) {
	var count int64
	if err := db.View().Model(RawOutput{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func listRawOutputs(db *db.DB, offset time.Time, limit int) ([]*core.Output, error) {
	if !offset.IsZero() {
		if err := db.Update().Where("created_at <= ?", offset.UnixNano()).Delete(RawOutput{}).Error; err != nil {
			return nil, err
		}
	}

	var raws []*RawOutput
	if err := db.View().Order("created_at").Limit(limit).Find(&raws).Error; err != nil {
		return nil, err
	}

	if len(raws) == limit {
		raws = trimSuffix(raws)
	}

	outputs := make([]*core.Output, 0, len(raws))
	for _, raw := range raws {
		var output core.Output
		if err := json.Unmarshal(raw.Data, &output); err != nil {
			return nil, err
		}

		outputs = append(outputs, &output)
	}

	core.SortOutputs(outputs)
	return outputs, nil
}
