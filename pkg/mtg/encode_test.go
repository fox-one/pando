package mtg

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"

	"github.com/fox-one/pando/pkg/routes"
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

	var proposal RawMessage = make([]byte, 100)
	_, _ = io.ReadFull(rand.Reader, proposal)
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
		body, err := Encode(append(values, newUUID(), routes.Routes{1, 2, 3}, decimal.NewFromFloat(2.123))...)
		require.Nil(t, err)

		data, err := Encrypt(body, pri, pub)
		require.Nil(t, err)

		t.Log(len(data))

		memo := base64.StdEncoding.EncodeToString(data)
		t.Log(len(memo), memo)

		assert.LessOrEqual(t, len(memo), 255)
	})

	t.Run("encode proposal", func(t *testing.T) {
		body, err := Encode(append(values, string(proposal))...)
		require.Nil(t, err)

		data, err := Encrypt(body, pri, pub)
		require.Nil(t, err)

		t.Log(len(data))

		memo := base64.StdEncoding.EncodeToString(data)
		t.Log(len(memo), memo)

		assert.LessOrEqual(t, len(memo), 255)
	})
}
