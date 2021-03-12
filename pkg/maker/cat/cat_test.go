package cat_test

import (
	"context"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/mock"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/makertest"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collaterals := mock.NewMockCollateralStore(ctrl)

	t.Run("bad-data", func(t *testing.T) {
		id := "not uuid"

		req := makertest.Next().WithBody(id)
		_, err := cat.From(req, collaterals)
		require.NotNil(t, err)
		assert.Equal(t, "Cat/bad-data", err.Error())
	})

	t.Run("non-existent", func(t *testing.T) {
		id := uuid.New()

		collaterals.EXPECT().
			Find(gomock.Any(), gomock.Eq(id)).
			DoAndReturn(func(_ context.Context, id string) (*core.Collateral, error) {
				return &core.Collateral{TraceID: id}, nil
			})

		req := makertest.Next().WithBody(types.UUID(id))
		_, err := cat.From(req, collaterals)
		require.NotNil(t, err)
		assert.Equal(t, "Cat/not-init", err.Error())
	})

	t.Run("exist", func(t *testing.T) {
		id := uuid.New()

		collaterals.EXPECT().
			Find(gomock.Any(), gomock.Eq(id)).
			DoAndReturn(func(_ context.Context, id string) (*core.Collateral, error) {
				return &core.Collateral{ID: 1, TraceID: id}, nil
			})
		req := makertest.Next().WithBody(types.UUID(id))
		_, err := cat.From(req, collaterals)
		require.Nil(t, err)
	})
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collaterals := mock.NewMockCollateralStore(ctrl)

	t.Run("list with single id", func(t *testing.T) {
		id := uuid.New()

		collaterals.EXPECT().
			Find(gomock.Any(), gomock.Eq(id)).
			DoAndReturn(func(_ context.Context, id string) (*core.Collateral, error) {
				return &core.Collateral{ID: 1, TraceID: id}, nil
			})

		req := makertest.Next().WithBody(types.UUID(id))
		cats, err := cat.List(req, collaterals)
		require.Nil(t, err)
		require.Len(t, cats, 1)
	})

	t.Run("list with nil uuid", func(t *testing.T) {
		collaterals.EXPECT().List(gomock.Any())
		req := makertest.Next().WithBody(uuid.Zero)
		cats, err := cat.List(req, collaterals)
		require.Nil(t, err)
		require.Len(t, cats, 0)
	})
}
