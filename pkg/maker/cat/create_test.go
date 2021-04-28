package cat_test

import (
	"context"
	"testing"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/mock"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/makertest"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gem := uuid.New()
	dai := uuid.New()
	name := "XIN"

	t.Run("not gov", func(t *testing.T) {
		collaterals := mock.NewMockCollateralStore(ctrl)
		oracles := mock.NewMockOracleStore(ctrl)

		req := makertest.Next().WithBody(types.UUID(gem), types.UUID(dai), name)
		err := cat.HandleCreate(collaterals, oracles)(req)
		require.NotNil(t, err)
		assert.Equal(t, "Cat/not-authorized", err.Error())
	})

	t.Run("gov", func(t *testing.T) {
		req := makertest.Next().WithBody(types.UUID(gem), types.UUID(dai), name)

		oracles := mock.NewMockOracleStore(ctrl)
		oracles.EXPECT().ListCurrent(gomock.Any()).DoAndReturn(func(_ context.Context) (number.Values, error) {
			return number.Values{
				gem: number.Decimal("100"),
				dai: number.Decimal("10"),
			}, nil
		})

		collaterals := mock.NewMockCollateralStore(ctrl)
		collaterals.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, c *core.Collateral) {
				assert.Equal(t, gem, c.Gem)
				assert.Equal(t, dai, c.Dai)
				assert.Equal(t, name, c.Name)
				assert.Equal(t, req.TraceID, c.TraceID)
				assert.Equal(t, req.Now, c.CreatedAt)
				assert.Equal(t, req.Version, c.Version)
				assert.Equal(t, "10", c.Price.String())
				assert.True(t, c.Debt.IsZero())
				assert.True(t, c.Ink.IsZero())
				assert.True(t, c.Rate.Equal(number.One))
				assert.Equal(t, req.Now, c.Rho)
				assert.True(t, c.Line.IsZero())
				assert.True(t, c.Duty.IsPositive())
				assert.True(t, c.Mat.IsPositive())
				assert.True(t, c.Chop.IsPositive())
				assert.True(t, c.Dust.IsPositive())
				assert.True(t, c.Dunk.IsPositive())
				assert.True(t, c.Beg.IsPositive())
				assert.True(t, c.TTL > 0)
				assert.True(t, c.Tau > c.TTL)
				assert.True(t, c.Live == 0)
			})

		req.Governors = []string{"a", "b"}
		err := cat.HandleCreate(collaterals, oracles)(req)
		require.Nil(t, err)
	})
}
