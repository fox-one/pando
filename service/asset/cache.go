package asset

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

func Cache(assetz core.AssetService) core.AssetService {
	return &cacheAssets{
		AssetService: assetz,
		assets:       cache.New(cache.NoExpiration, cache.NoExpiration),
		sf:           &singleflight.Group{},
	}
}

type cacheAssets struct {
	core.AssetService
	assets *cache.Cache
	sf     *singleflight.Group
}

func (c *cacheAssets) Find(ctx context.Context, id string) (*core.Asset, error) {
	asset, err, _ := c.sf.Do(id, func() (interface{}, error) {
		if v, ok := c.assets.Get(id); ok {
			return v, nil
		}

		asset, err := c.AssetService.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		c.assets.SetDefault(id, asset)
		return asset, nil
	})

	if err != nil {
		return nil, err
	}

	return asset.(*core.Asset), nil
}
