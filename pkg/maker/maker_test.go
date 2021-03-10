package maker_test

import (
	"context"
	"testing"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/mock"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/sys"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pando/store/collateral"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newReq() *maker.Request {
	r := &maker.Request{
		Now:      time.Now(),
		Version:  1,
		TraceID:  uuid.New(),
		Sender:   uuid.New(),
		FollowID: uuid.New(),
		AssetID:  uuid.New(),
		Amount:   number.One,
	}

	return r
}

func step(r *maker.Request, action core.Action, values ...interface{}) *maker.Request {
	r = r.WithBody(values...)
	r.Version += 1
	r.Action = action
	r.Now = r.Now.Add(time.Minute)

	return r
}

func TestMaker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := newReq()

	t.Run("Sys", func(t *testing.T) {
		t.Run("withdraw", func(t *testing.T) {
			asset := uuid.New()
			amount := number.One
			opponent := uuid.New()

			wallets := mock.NewMockWalletStore(ctrl)
			wallets.EXPECT().CreateTransfers(gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, transfers []*core.Transfer) {
					require.Len(t, transfers, 1)
					transfer := transfers[0]
					assert.Equal(t, asset, transfer.AssetID)
					assert.Equal(t, amount.String(), transfer.Amount.String())
					require.Len(t, transfer.Opponents, 1)
					assert.Equal(t, opponent, transfer.Opponents[0])
					assert.True(t, transfer.Threshold > 0)
				}).Times(1)

			h := sys.HandleWithdraw(wallets)
			{
				req := step(req, core.ActionSysWithdraw, types.UUID(asset), amount, types.UUID(opponent))
				assert.NotNil(t, h(req), "gov required")
			}

			{
				req := step(req, core.ActionSysWithdraw, types.UUID(asset), amount, types.UUID(opponent))
				req.Gov = true
				assert.Nil(t, h(req))
			}
		})
	})

	collaterals := collateral.Memory()

	t.Run("Cat", func(t *testing.T) {
		assets := mock.NewMockAssetStore(ctrl)
		assetIds := make(map[string]bool)

		assets.EXPECT().Find(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, id string) (*core.Asset, error) {
				assetIds[id] = true
				return &core.Asset{ID: id, Symbol: id, ChainID: id}, nil
			}).AnyTimes()

		oracles := mock.NewMockOracleStore(ctrl)
		oracles.EXPECT().ListCurrent(gomock.Any()).
			DoAndReturn(func(_ context.Context) (number.Values, error) {
				values := number.Values{}
				for id := range assetIds {
					values.Set(id, number.One)
				}

				return values, nil
			})

		t.Run("create", func(t *testing.T) {
			gem := uuid.New()
			dai := uuid.New()
			name := "ETH-A"

			h := cat.HandleCreate(collaterals, oracles, assets, assets)

			t.Run("not gov", func(t *testing.T) {
				req := step(req, core.ActionCatCreate, types.UUID(gem), types.UUID(dai), name)
				require.NotNil(t, h(req), "gov required")
			})

			t.Run("gov", func(t *testing.T) {
				req := step(req, core.ActionCatCreate, types.UUID(gem), types.UUID(dai), name)
				req.Gov = true
				require.Nil(t, h(req), "should create a new cat")
			})

			cats, _ := collaterals.List(req.Context())
			require.Len(t, cats, 1)

			c := cats[0]
			assert.Equal(t, gem, c.Gem)
			assert.Equal(t, dai, c.Dai)
			assert.Equal(t, name, c.Name)
			assert.Equal(t, "1", c.Price.String())
		})

		t.Run("edit", func(t *testing.T) {

		})
	})
}
