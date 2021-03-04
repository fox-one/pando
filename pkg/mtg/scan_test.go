package mtg

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	var (
		typ int8 = 1
		uid      = newUUID()
		str      = "123"
	)

	body, err := Encode(typ, uid, str)
	require.Nil(t, err)

	var (
		dtyp int8
		duid uuid.UUID
		dstr string
	)

	remain, err := Scan(body, &dtyp)
	require.Nil(t, err)
	assert.Equal(t, dtyp, typ)

	_, err = Scan(remain, &duid, &dstr)
	require.Nil(t, err)

	assert.Equal(t, uid.String(), duid.String())
	assert.Equal(t, str, dstr)
}

func TestScanStruct(t *testing.T) {
	type Foo struct {
		A uuid.UUID
		B int
	}

	var (
		a = newUUID()
		b = 10
	)
	body, _ := Encode(a, b)

	var foo Foo
	require.Nil(t, ScanStructs(body, &foo))
	assert.Equal(t, a, foo.A)
	assert.Equal(t, b, foo.B)
}
