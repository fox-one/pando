package maker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerVersion(t *testing.T) {
	versions := make(map[int]HandlerFunc)
	for i := 0; i < 10; i++ {
		v := i
		versions[v] = func(r *Request) error {
			assert.LessOrEqual(t, v, r.SysVersion)
			return nil
		}
	}

	h := HandlerVersion(versions)

	for i := 0; i < 20; i++ {
		r := &Request{SysVersion: i}
		_ = h(r)
	}
}
