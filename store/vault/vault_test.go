package vault

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/store/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var noContext = context.TODO()

func TestVault(t *testing.T) {
	conn, err := dbtest.Connect()
	require.Nil(t, err)

	t.Cleanup(func() {
		dbtest.Disconnect(conn)
	})

	vaults := New(conn)

	t.Run("Create", testVaultCreate(vaults))
	t.Run("Find", testVaultFind(vaults))
	t.Run("Update", testVaultUpdate(vaults))
	t.Run("List", testVaultList(vaults))
	t.Run("CountCollateral", testVaultCountCollateral(vaults))
}

func testVaultCreate(vaults core.VaultStore) func(t *testing.T) {
	return func(t *testing.T) {
		vat := rawVat(t)
		require.Nil(t, vaults.Create(noContext, vat))
		assert.True(t, vat.ID > 0)
	}
}

func testVaultFind(vaults core.VaultStore) func(t *testing.T) {
	return func(t *testing.T) {
		want := rawVat(t)
		got, err := vaults.Find(noContext, want.TraceID)
		require.Nil(t, err)

		diff := cmp.Diff(want, got, cmpopts.IgnoreFields(core.Vault{}, "ID", "UpdatedAt"))
		assert.Empty(t, diff)
	}
}

func testVaultUpdate(vaults core.VaultStore) func(t *testing.T) {
	return func(t *testing.T) {
		vat := rawVat(t)

		before, err := vaults.Find(noContext, vat.TraceID)
		require.Nil(t, err)

		before.Art = before.Art.Add(number.Decimal("100"))
		before.Ink = before.Ink.Add(number.Decimal("200"))

		require.Nil(t, vaults.Update(noContext, before, before.Version+10))

		after, err := vaults.Find(noContext, vat.TraceID)
		require.Nil(t, err)

		assert.Equal(t, before.Ink.String(), after.Ink.String())
		assert.Equal(t, before.Art.String(), after.Art.String())
		assert.Equal(t, before.Version, after.Version)
	}
}

func testVaultList(vaults core.VaultStore) func(t *testing.T) {
	return func(t *testing.T) {
		vats, err := vaults.List(noContext, core.ListVaultRequest{
			Limit: 100,
		})

		require.Nil(t, err)
		assert.Len(t, vats, 1)
	}
}

func testVaultCountCollateral(vaults core.VaultStore) func(t *testing.T) {
	return func(t *testing.T) {
		vat := rawVat(t)
		counts, err := vaults.CountCollateral(noContext)
		require.Nil(t, err)
		assert.Equal(t, int64(1), counts[vat.CollateralID])
	}
}

func rawVat(t *testing.T) *core.Vault {
	out, err := ioutil.ReadFile("testdata/vault.json")
	require.Nil(t, err)

	var vat core.Vault
	require.Nil(t, json.Unmarshal(out, &vat))

	return &vat
}
