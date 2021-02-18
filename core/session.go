package core

import (
	"context"
)

type Session interface {
	// Login return user mixin id
	Login(ctx context.Context, accessToken string) (*User, error)
}
