package uuid

import (
	"crypto/md5"
	"io"

	"github.com/gofrs/uuid"
)

type UUID = uuid.UUID

func New() string {
	return uuid.Must(uuid.NewV4()).String()
}

var Zero = uuid.Nil

func IsUUID(id string) bool {
	_, err := FromString(id)
	return err == nil
}

func MD5(input string) string {
	h := md5.New()
	_, _ = io.WriteString(h, input)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}

func Modify(id, modifier string) string {
	ns, err := uuid.FromString(id)
	if err != nil {
		panic(err)
	}

	return uuid.NewV5(ns, modifier).String()
}

func FromString(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

func FromBytes(b []byte) (uuid.UUID, error) {
	return uuid.FromBytes(b)
}

func IsNil(id string) bool {
	uid, err := FromString(id)
	return err != nil || uid == uuid.Nil
}
