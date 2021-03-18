package core

import (
	"context"
)

// Session define operations to parse authorization token
type Session interface {
	// Login return user mixin id
	Login(ctx context.Context, accessToken string) (*User, error)
}
