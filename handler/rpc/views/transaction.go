package views

import (
	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
)

func Transaction(tx *core.Transaction) *api.Transaction {
	return &api.Transaction{
		Id:         tx.TraceID,
		CreatedAt:  Time(&tx.CreatedAt),
		AssetId:    tx.AssetID,
		Amount:     tx.Amount.String(),
		Action:     api.Action(tx.Action),
		Status:     api.Transaction_Status(tx.Status),
		Msg:        tx.Message,
		Parameters: string(tx.Parameters),
	}
}
