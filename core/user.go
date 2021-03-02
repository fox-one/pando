package core

import (
	"context"
	"time"
)

type (
	User struct {
		ID          int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt   time.Time `json:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty"`
		Version     int64     `json:"version,omitempty"`
		MixinID     string    `sql:"size:36" json:"mixin_id,omitempty"`
		Role        string    `sql:"size:24" json:"role,omitempty"`
		Lang        string    `sql:"size:36" json:"lang,omitempty"`
		Name        string    `sql:"size:64" json:"name,omitempty"`
		Avatar      string    `sql:"size:255" json:"avatar,omitempty"`
		AccessToken string    `sql:"size:512" json:"access_token,omitempty"`
	}

	UserStore interface {
		Save(ctx context.Context, user *User) error
		Find(ctx context.Context, mixinID string) (*User, error)
	}

	UserService interface {
		Find(ctx context.Context, mixinID string) (*User, error)
		Login(ctx context.Context, token string) (*User, error)
		Auth(ctx context.Context, code string) (string, error)
	}
)
