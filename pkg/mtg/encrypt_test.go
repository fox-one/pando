package mtg

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	_, clientPrivate, _ := ed25519.GenerateKey(rand.Reader)
	serverPublic, serverPrivate, _ := ed25519.GenerateKey(rand.Reader)

	body := make([]byte, 100)
	_, _ = io.ReadFull(rand.Reader, body)

	data, err := Encrypt(body, clientPrivate, serverPublic)
	require.Nil(t, err)

	t.Log("encrypt", len(data))

	dbody, err := Decrypt(data, serverPrivate)
	require.Nil(t, err)

	t.Log("decrypt", len(dbody))

	assert.Equal(t, body, dbody)
}
