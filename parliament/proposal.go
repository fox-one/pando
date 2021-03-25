package parliament

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

func (s *parliament) renderProposalItems(ctx context.Context, action core.Action, data []byte) (items []Item) {
	switch action {
	case core.ActionCatCreate:
		var (
			gem, dai uuid.UUID
			name     string
		)

		_, _ = mtg.Scan(data, &gem, &dai, &name)

		items = []Item{
			{
				Key:   "name",
				Value: name,
			},
			{
				Key:    "gem",
				Value:  s.fetchAssetSymbol(ctx, gem.String()),
				Action: assetAction(gem.String()),
			},
			{
				Key:    "dai",
				Value:  s.fetchAssetSymbol(ctx, dai.String()),
				Action: assetAction(dai.String()),
			},
		}
	case core.ActionCatEdit:
		var id uuid.UUID
		data, err := mtg.Scan(data, &id)

		items = []Item{
			{
				Key:   "cat",
				Value: s.fetchCatName(ctx, id.String()),
			},
		}

		for {
			var item Item
			if data, err = mtg.Scan(data, &item.Key, &item.Value); err != nil {
				break
			}

			items = append(items, item)
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
		items = []Item{
			{
				Key:    "asset",
				Value:  s.fetchAssetSymbol(ctx, id.String()),
				Action: assetAction(id.String()),
			},
			{
				Key:   "price",
				Value: number.Humanize(price),
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

		items = []Item{
			{
				Key:   "asset",
				Value: s.fetchAssetSymbol(ctx, id.String()),
			},
		}

		for {
			var item Item
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
		items = []Item{
			{
				Key:    "asset",
				Value:  s.fetchAssetSymbol(ctx, id.String()),
				Action: assetAction(id.String()),
			},
			{
				Key:   "price",
				Value: number.Humanize(price),
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
		items = []Item{
			{
				Key:   "id",
				Value: id.String(),
			},
			{
				Key:   "key",
				Value: base64.StdEncoding.EncodeToString(publicKey),
			},
		}
	case core.ActionOracleDeny:
		var id uuid.UUID

		_, _ = mtg.Scan(data, &id)
		items = []Item{
			{
				Key:   "id",
				Value: id.String(),
			},
		}
	case core.ActionSysWithdraw:
		var (
			assetID  uuid.UUID
			amount   decimal.Decimal
			opponent uuid.UUID
		)

		_, _ = mtg.Scan(data, &assetID, &amount, &opponent)
		items = []Item{
			{
				Key:    "asset",
				Value:  fmt.Sprintf("%s %s", number.Humanize(amount), s.fetchAssetSymbol(ctx, assetID.String())),
				Action: assetAction(assetID.String()),
			},
			{
				Key:    "opponent",
				Value:  s.fetchUserName(ctx, opponent.String()),
				Action: userAction(opponent.String()),
			},
		}
	}

	return
}
