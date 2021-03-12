package cat_test

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/fox-one/pando/mock"
	"github.com/fox-one/pando/pkg/maker/cat"
	"github.com/fox-one/pando/pkg/maker/makertest"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestEdit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("not gov", func(t *testing.T) {
		collaterals := mock.NewMockCollateralStore(ctrl)

		req := makertest.Next()
		err := cat.HandleEdit(collaterals)(req)
		require.NotNil(t, err)
		assert.Equal(t, "Cat/not-authorized", err.Error())
	})
}
