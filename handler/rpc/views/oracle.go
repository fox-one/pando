package views

import (
	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
)

func Oracle(oracle *core.Oracle) *api.Oracle {
	return &api.Oracle{
		AssetId:   oracle.AssetID,
		Hop:       int32(oracle.Hop),
		Current:   oracle.Current.String(),
		Next:      oracle.Next.String(),
		PeekAt:    Time(&oracle.PeekAt),
		Threshold: int32(oracle.Threshold),
		Governors: oracle.Governors,
	}
}
