package parliament

import (
	"context"

	"github.com/fox-one/pando/pkg/uuid"
)

func (s *parliament) fetchAssetSymbol(ctx context.Context, assetID string) string {
	if uuid.IsNil(assetID) {
		return "ALL"
	}

	coin, err := s.assetz.Find(ctx, assetID)
	if err != nil {
		return "NULL"
	}

	return coin.Symbol
}

func (s *parliament) fetchUserName(ctx context.Context, userID string) string {
	user, err := s.userz.Find(ctx, userID)
	if err != nil {
		return "NULL"
	}

	return user.Name
}

func (s *parliament) fetchCatName(ctx context.Context, id string) string {
	if uuid.IsNil(id) {
		return "ALL"
	}

	c, err := s.collaterals.Find(ctx, id)
	if err != nil {
		return "NULL"
	}

	return c.Name
}

func (s *parliament) fetchCatGemDai(ctx context.Context, id string) (gem, dai string) {
	c, err := s.collaterals.Find(ctx, id)
	if err != nil {
		return
	}

	gem = s.fetchAssetSymbol(ctx, c.Gem)
	dai = s.fetchAssetSymbol(ctx, c.Dai)

	return
}
