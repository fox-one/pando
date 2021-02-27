package request

import (
	"context"

	"github.com/fox-one/pando/core"
)

type key int

const (
	userKey key = iota
	clientIPKey
)

func WithUser(ctx context.Context, user *core.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFrom(ctx context.Context) (*core.User, bool) {
	user, ok := ctx.Value(userKey).(*core.User)
	return user, ok
}

func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

func ClientIPFrom(ctx context.Context) string {
	return ctx.Value(clientIPKey).(string)
}
