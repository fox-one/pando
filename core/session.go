package core

import (
	"net/http"
)

// Session define operations to parse authorization token
type Session interface {
	// Login return user
	Login(r *http.Request) (*User, error)
}
