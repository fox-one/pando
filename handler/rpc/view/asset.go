package view

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
)

func Asset(asset *core.Asset, chain *core.Asset) *api.Asset {
	view := &api.Asset{
		Id:      asset.ID,
		Name:    asset.Name,
		Symbol:  asset.Symbol,
		Logo:    asset.Logo,
		ChainId: asset.ChainID,
		Price:   asset.Price.String(),
	}

	if chain != nil && asset.ChainID != "" {
		view.Chain = Asset(chain, nil)
	}

	return view
}
