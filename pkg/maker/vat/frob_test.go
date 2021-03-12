package vat_test

import (
	"context"
	"testing"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/mock"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/makertest"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := makertest.Next()

	wallets := mock.NewMockWalletStore(ctrl)
	wallets.EXPECT().CreateTransfers(gomock.Any(), gomock.Any()).AnyTimes()

	v := &core.Vault{
		ID:           1,
		TraceID:      uuid.New(),
		Version:      1,
		UserID:       uuid.New(),
		CollateralID: uuid.New(),
	}

	c := &core.Collateral{
		ID:      1,
		TraceID: v.CollateralID,
		Version: 1,
		Gem:     uuid.New(),
		Dai:     uuid.New(),
		Rate:    number.One,
		Rho:     r.Now,
		Line:    number.Decimal("1000"),
		Dust:    number.Decimal("100"),
		Price:   number.Decimal("15"),
		Mat:     number.Decimal("1.2"),
		Duty:    number.Decimal("1.06"),
		Live:    1,
	}

	vaults := mock.NewMockVaultStore(ctrl)
	vaults.EXPECT().Find(gomock.Any(), gomock.Eq(v.TraceID)).
		DoAndReturn(func(_ context.Context, id string) (*core.Vault, error) {
			return v, nil
		}).AnyTimes()

	vaults.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, v *core.Vault, version int64) {
			v.Version = version
		}).AnyTimes()

	vaults.EXPECT().FindEvent(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id string, version int64) (*core.VaultEvent, error) {
			return &core.VaultEvent{}, nil
		}).AnyTimes()

	vaults.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, event *core.VaultEvent) {
			t.Log("create event", event.Dink, event.Dart, event.Debt)
			if noDebt := event.Debt.IsZero() && event.Dart.IsZero(); !noDebt {
				assert.True(t, event.Dart.Mul(event.Debt).IsPositive())
			}
		}).AnyTimes()

	collaterals := mock.NewMockCollateralStore(ctrl)
	collaterals.EXPECT().Find(gomock.Any(), gomock.Eq(c.TraceID)).
		DoAndReturn(func(_ context.Context, id string) (*core.Collateral, error) {
			return c, nil
		}).AnyTimes()

	collaterals.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, c *core.Collateral, version int64) {
			c.Version = version
		}).AnyTimes()

	// fold
	{
		r := makertest.Next().WithBody(types.UUID(c.TraceID))
		r.Now = r.Now.Add(time.Hour)
		require.Nil(t, cat.HandleFold(collaterals)(r))
		assert.Equal(t, r.Version, c.Version)
		assert.True(t, c.Rate.GreaterThan(number.One))
		assert.Equal(t, r.Now, c.Rho)
	}

	t.Run("deposit", func(t *testing.T) {
		req := makertest.Next().WithBody(types.UUID(v.TraceID))
		req.AssetID = c.Gem
		req.Amount = number.Decimal("100")
		err := vat.HandleDeposit(collaterals, vaults, wallets)(req)
		require.Nil(t, err)
	})

	t.Run("generate less than dust", func(t *testing.T) {
		debt := c.Dust.Div(decimal.NewFromInt(2))
		req := makertest.Next().WithBody(types.UUID(v.TraceID), debt)
		req.Sender = v.UserID
		err := vat.HandleGenerated(collaterals, vaults, wallets)(req)
		require.NotNil(t, err)
		assert.Equal(t, "Vat/dust", err.Error())
	})

	t.Run("generate", func(t *testing.T) {
		debt := c.Dust
		req := makertest.Next().WithBody(types.UUID(v.TraceID), debt)
		req.Sender = v.UserID
		err := vat.HandleGenerated(collaterals, vaults, wallets)(req)
		require.Nil(t, err)
	})

	t.Run("payback", func(t *testing.T) {
		debt := v.Art.Mul(c.Rate)
		debt = debt.Truncate(8)
		req := makertest.Next().WithBody(types.UUID(v.TraceID))
		req.AssetID = c.Dai
		req.Amount = debt
		err := vat.HandlePayback(collaterals, vaults, wallets)(req)
		require.Nil(t, err)
	})
}
