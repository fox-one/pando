package maker

import (
	"context"
)

type key int

const (
	versionKey key = iota
)

func WithVersion(ctx context.Context, version int64) context.Context {
	return context.WithValue(ctx, versionKey, version)
}

func VersionFrom(ctx context.Context) int64 {
	v, _ := ctx.Value(versionKey).(int64)
	return v
}
