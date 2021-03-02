package notifier

import (
	"crypto/md5"
	"io"
	"math/big"

	"github.com/gofrs/uuid"
)

func mixinRawTransactionTraceId(hash string, index uint8) string {
	h := md5.New()
	_, _ = io.WriteString(h, hash)
	b := new(big.Int).SetInt64(int64(index))
	h.Write(b.Bytes())
	s := h.Sum(nil)
	s[6] = (s[6] & 0x0f) | 0x30
	s[8] = (s[8] & 0x3f) | 0x80
	sid, err := uuid.FromBytes(s)
	if err != nil {
		panic(err)
	}

	return sid.String()
}
