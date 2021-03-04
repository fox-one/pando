package mtg

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func TestEncode(t *testing.T) {
	pub, pri, _ := ed25519.GenerateKey(rand.Reader)
	values := []interface{}{1, newUUID(), newUUID()}

	t.Run("encode add action", func(t *testing.T) {
		body, err := Encode(append(values, decimal.NewFromFloat(0.001))...)
		require.Nil(t, err)

		data, err := Encrypt(body, pri, pub)
		require.Nil(t, err)

		t.Log(len(data))

		memo := base64.StdEncoding.EncodeToString(data)
		t.Log(len(memo), memo)

		assert.LessOrEqual(t, len(memo), 255)
	})

	t.Run("encode swap action", func(t *testing.T) {
		body, err := Encode(append(values, newUUID(), decimal.NewFromFloat(2.123))...)
		require.Nil(t, err)

		data, err := Encrypt(body, pri, pub)
		require.Nil(t, err)

		t.Log(len(data))

		memo := base64.StdEncoding.EncodeToString(data)
		t.Log(len(memo), memo)

		assert.LessOrEqual(t, len(memo), 255)
	})

	t.Run("encode struct", func(t *testing.T) {
		type Foo struct {
			A uuid.UUID
			B decimal.Decimal
			C int
			D string
		}

		a := Foo{
			A: newUUID(),
			B: decimal.NewFromInt(10),
			C: 10,
			D: "bar",
		}

		b1, err := EncodeStruct(a)
		require.Nil(t, err)

		b2, err := EncodeStruct(&a)
		require.Nil(t, err)

		assert.Equal(t, b1, b2)

		var v Foo
		require.Nil(t, ScanAll(b1, &v.A, &v.B, &v.C, &v.D))
		assert.Equal(t, a, v)
	})
}
