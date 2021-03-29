package asset

import (
	"context"
	"strings"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
)

func New(c *mixin.Client) core.AssetService {
	return &assetService{c: c}
}

type assetService struct {
	c *mixin.Client
}

func (s *assetService) Find(ctx context.Context, id string) (*core.Asset, error) {
	asset, err := s.c.ReadAsset(ctx, id)
	if err != nil {
		if mixin.IsErrorCodes(err, 10002) {
			return &core.Asset{ID: id}, err
		}

		return nil, err
	}

	return convertAsset(asset), nil
}

func (s *assetService) List(ctx context.Context) ([]*core.Asset, error) {
	assets, err := mixin.ReadMultisigAssets(ctx)
	if err != nil {
		return nil, err
	}

	const blankLogo = "yH_I5b0GiV2zDmvrXRyr3bK5xusjfy5q7FX3lw3mM2Ryx4Dfuj6Xcw8SHNRnDKm7ZVE3_LvpKlLdcLrlFQUBhds"
	var idx int

	for _, asset := range assets {
		if asset.PriceUSD.IsZero() && strings.Contains(asset.IconURL, blankLogo) {
			continue
		}

		assets[idx] = asset
		idx += 1
	}

	return convertAssets(assets[:idx]), nil
}

func convertAsset(asset *mixin.Asset) *core.Asset {
	return &core.Asset{
		ID:      asset.AssetID,
		Name:    asset.Name,
		Symbol:  asset.Symbol,
		Logo:    asset.IconURL,
		ChainID: asset.ChainID,
		Price:   asset.PriceUSD,
	}
}

func convertAssets(assets []*mixin.Asset) []*core.Asset {
	out := make([]*core.Asset, len(assets))
	for idx, asset := range assets {
		out[idx] = convertAsset(asset)
	}

	return out
}
