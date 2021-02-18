package mtg

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/fox-one/pando/pkg/routes"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	var (
		typ   int8       = 1
		uid              = newUUID()
		str              = "123"
		route            = routes.Routes{1, 2, 3}
		data  RawMessage = make([]byte, 100)
	)

	_, _ = io.ReadFull(rand.Reader, data)

	body, err := Encode(typ, uid, str, route, string(data))
	require.Nil(t, err)

	var (
		dtyp   int8
		duid   uuid.UUID
		dstr   string
		droute routes.Routes
		ddata  RawMessage
	)

	remain, err := Scan(body, &dtyp)
	require.Nil(t, err)
	assert.Equal(t, dtyp, typ)

	_, err = Scan(remain, &duid, &dstr, &droute, &ddata)
	require.Nil(t, err)

	assert.Equal(t, uid.String(), duid.String())
	assert.Equal(t, str, dstr)
	assert.Equal(t, route.String(), droute.String())
	assert.Equal(t, data, ddata)
}
