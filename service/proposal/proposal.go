package proposal

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

func New(
	assetz core.AssetService,
	userz core.UserService,
	cats core.CollateralStore,
) core.ProposalService {
	return &proposalService{
		assetz: assetz,
		userz:  userz,
		cats:   cats,
		cache:  cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

type proposalService struct {
	assetz core.AssetService
	userz  core.UserService
	cats   core.CollateralStore
	cache  *cache.Cache
}

func (s *proposalService) get(key string) (string, bool) {
	v, ok := s.cache.Get(key)
	if ok {
		return v.(string), true
	}

	return "", false
}

func (s *proposalService) set(key, value string) {
	s.cache.SetDefault(key, value)
}

func (s *proposalService) fetchAssetSymbol(ctx context.Context, assetID string) string {
	if v, ok := s.get(assetID); ok {
		return v
	}

	if uuid.IsNil(assetID) {
		return "ALL"
	}

	coin, err := s.assetz.Find(ctx, assetID)
	if err != nil {
		return "NULL"
	}

	s.set(assetID, coin.Symbol)
	return coin.Symbol
}

func (s *proposalService) fetchUserName(ctx context.Context, userID string) string {
	if v, ok := s.get(userID); ok {
		return v
	}

	user, err := s.userz.Find(ctx, userID)
	if err != nil {
		return "NULL"
	}

	s.set(userID, user.Name)
	return user.Name
}

func (s *proposalService) fetchCatName(ctx context.Context, id string) string {
	if v, ok := s.get(id); ok {
		return v
	}

	if uuid.IsNil(id) {
		return "ALL"
	}

	c, err := s.cats.Find(ctx, id)
	if err != nil {
		return "NULL"
	}

	s.set(id, c.Name)
	return c.Name
}

func (s *proposalService) ListItems(ctx context.Context, p *core.Proposal) ([]core.ProposalItem, error) {
	data, _ := base64.StdEncoding.DecodeString(p.Data)

	var items []core.ProposalItem

	switch p.Action {
	case core.ActionCatCreate:
		var (
			gem, dai uuid.UUID
			name     string
		)

		_, _ = mtg.Scan(data, &gem, &dai, &name)

		items = []core.ProposalItem{
			{
				Key:   "name",
				Value: name,
			},
			{
				Key:    "gem",
				Value:  gem.String(),
				Hint:   s.fetchAssetSymbol(ctx, gem.String()),
				Action: assetAction(gem.String()),
			},
			{
				Key:    "dai",
				Value:  dai.String(),
				Hint:   s.fetchAssetSymbol(ctx, dai.String()),
				Action: assetAction(dai.String()),
			},
		}
	case core.ActionCatEdit:
		var id uuid.UUID
		data, err := mtg.Scan(data, &id)

		items = []core.ProposalItem{
			{
				Key:   "cat",
				Value: id.String(),
				Hint:  s.fetchCatName(ctx, id.String()),
			},
		}

		for {
			var item core.ProposalItem
			if data, err = mtg.Scan(data, &item.Key, &item.Value); err != nil {
				break
			}

			items = append(items, item)
		}
	case core.ActionCatMove:
		var (
			from, to uuid.UUID
			amount   decimal.Decimal
		)
		_, _ = mtg.Scan(data, &from, &to, &amount)

		items = []core.ProposalItem{
			{
				Key:   "from",
				Value: from.String(),
				Hint:  s.fetchCatName(ctx, from.String()),
			},
			{
				Key:   "to",
				Value: to.String(),
				Hint:  s.fetchCatName(ctx, to.String()),
			},
			{
				Key:   "amount",
				Value: amount.String(),
				Hint:  number.Humanize(amount),
			},
		}
	case core.ActionCatGain:
		var (
			id      uuid.UUID
			amount  decimal.Decimal
			receipt uuid.UUID
		)
		_, _ = mtg.Scan(data, &id, &amount, &receipt)

		items = []core.ProposalItem{
			{
				Key:   "cat",
				Value: id.String(),
				Hint:  s.fetchCatName(ctx, id.String()),
			},
			{
				Key:   "amount",
				Value: amount.String(),
				Hint:  number.Humanize(amount),
			},
			{
				Key:    "receipt",
				Value:  receipt.String(),
				Hint:   s.fetchUserName(ctx, receipt.String()),
				Action: userAction(receipt.String()),
			},
		}
	case core.ActionOracleCreate:
		var (
			id        uuid.UUID
			price     decimal.Decimal
			hop       int64
			threshold int64
			ts        int64
		)

		_, _ = mtg.Scan(data, &id, &price, &hop, &threshold, &ts)
		items = []core.ProposalItem{
			{
				Key:    "asset",
				Value:  id.String(),
				Hint:   s.fetchAssetSymbol(ctx, id.String()),
				Action: assetAction(id.String()),
			},
			{
				Key:   "price",
				Value: price.String(),
				Hint:  number.Humanize(price),
			},
			{
				Key:   "hop",
				Value: (time.Duration(hop) * time.Second).String(),
			},
			{
				Key:   "threshold",
				Value: cast.ToString(threshold),
			},
			{
				Key:   "ts",
				Value: time.Unix(ts, 0).Format(time.RFC3339),
			},
		}
	case core.ActionOracleEdit:
		var id uuid.UUID
		data, err := mtg.Scan(data, &id)

		items = []core.ProposalItem{
			{
				Key:    "asset",
				Value:  id.String(),
				Hint:   s.fetchAssetSymbol(ctx, id.String()),
				Action: assetAction(id.String()),
			},
		}

		for {
			var item core.ProposalItem
			if data, err = mtg.Scan(data, &item.Key, &item.Value); err != nil {
				break
			}

			items = append(items, item)
		}
	case core.ActionOraclePoke:
		var (
			id    uuid.UUID
			price decimal.Decimal
			ts    int64
		)

		_, _ = mtg.Scan(data, &id, &price, &ts)
		items = []core.ProposalItem{
			{
				Key:    "asset",
				Value:  id.String(),
				Hint:   s.fetchAssetSymbol(ctx, id.String()),
				Action: assetAction(id.String()),
			},
			{
				Key:   "price",
				Value: price.String(),
				Hint:  number.Humanize(price),
			},
			{
				Key:   "ts",
				Value: time.Unix(ts, 0).Format(time.RFC3339),
			},
		}
	case core.ActionOracleRely:
		var (
			id        uuid.UUID
			publicKey types.RawMessage
		)

		_, _ = mtg.Scan(data, &id, &publicKey)
		items = []core.ProposalItem{
			{
				Key:    "id",
				Value:  id.String(),
				Hint:   s.fetchUserName(ctx, id.String()),
				Action: userAction(id.String()),
			},
			{
				Key:   "key",
				Value: base64.StdEncoding.EncodeToString(publicKey),
			},
		}
	case core.ActionOracleDeny:
		var id uuid.UUID

		_, _ = mtg.Scan(data, &id)
		items = []core.ProposalItem{
			{
				Key:    "id",
				Value:  id.String(),
				Hint:   s.fetchUserName(ctx, id.String()),
				Action: userAction(id.String()),
			},
		}
	case core.ActionSysWithdraw:
		var (
			assetID  uuid.UUID
			amount   decimal.Decimal
			opponent uuid.UUID
		)

		_, _ = mtg.Scan(data, &assetID, &amount, &opponent)
		items = []core.ProposalItem{
			{
				Key:    "asset",
				Value:  assetID.String(),
				Hint:   s.fetchAssetSymbol(ctx, assetID.String()),
				Action: assetAction(assetID.String()),
			},
			{
				Key:   "amount",
				Value: amount.String(),
				Hint:  number.Humanize(amount),
			},
			{
				Key:    "opponent",
				Value:  opponent.String(),
				Hint:   s.fetchUserName(ctx, opponent.String()),
				Action: userAction(opponent.String()),
			},
		}
	case core.ActionSysProperty:
		var key, value string
		_, _ = mtg.Scan(data, &key, &value)
		items = []core.ProposalItem{
			{
				Key:   key,
				Value: value,
			},
		}
	}

	return items, nil
}

func assetAction(id string) string {
	return fmt.Sprintf("https://mixin.one/snapshots/%s", id)
}

func userAction(id string) string {
	return mixin.URL.Users(id)
}
