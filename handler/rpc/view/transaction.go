package view

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
)

func Transaction(tx *core.Transaction) *api.Transaction {
	return &api.Transaction{
		Id:        tx.TraceID,
		CreatedAt: Time(&tx.CreatedAt),
		TargetId:  tx.TargetID,
		AssetId:   tx.AssetID,
		Amount:    tx.Amount.String(),
		Action:    int32(tx.Action),
		Status:    int32(tx.Status),
		Data:      string(tx.Data),
	}
}
