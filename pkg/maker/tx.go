package maker

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/shopspring/decimal"
)

type Tx struct {
	Now      time.Time
	Version  int64
	TraceID  string
	FollowID string
	Sender   string
	AssetID  string
	Action   int
	TargetID string
	Gov      bool
	Amount   decimal.Decimal
	// transfer out
	Transfers []*core.Transfer
}

func (tx *Tx) Transfer(traceID, assetID, opponentID, memo string, amount decimal.Decimal) {
	tx.Transfers = append(tx.Transfers, &core.Transfer{
		TraceID:   traceID,
		AssetID:   assetID,
		Amount:    amount,
		Memo:      memo,
		Threshold: 1,
		Opponents: []string{opponentID},
	})
}

func EncodeMemo(module, id, source string) string {
	return core.TransferAction{
		Module: module,
		ID:     id,
		Source: source,
	}.Encode()
}
