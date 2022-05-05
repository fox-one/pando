package stater

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/shopspring/decimal"
)

func cacheAssetPrice(assetz core.AssetService) core.AssetService {
	return &priceCache{
		AssetService: assetz,
		tickers:      map[string]ticker{},
	}
}

type ticker struct {
	price decimal.Decimal
	time  time.Time
}

type priceCache struct {
	core.AssetService
	tickers map[string]ticker
}

func (s *priceCache) ReadPrice(ctx context.Context, assetID string, at time.Time) (decimal.Decimal, error) {
	// don't cache a changeable price
	if time.Since(at) < 0 {
		return s.AssetService.ReadPrice(ctx, assetID, at)
	}

	key := assetID + at.Format(time.RFC3339)
	if ticker, ok := s.tickers[key]; ok && ticker.time.Equal(at) {
		return ticker.price, nil
	}

	p, err := s.AssetService.ReadPrice(ctx, assetID, at)
	if err != nil {
		return decimal.Zero, err
	}

	s.tickers[key] = ticker{
		price: p,
		time:  at,
	}

	return p, nil
}
