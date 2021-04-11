package session

import (
	"net/http"

	"github.com/bluele/gcache"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/internal/request"
)

type cacheSession struct {
	core.Session
	tokens gcache.Cache
}

func (s *cacheSession) Login(r *http.Request) (*core.User, error) {
	accessToken := request.ExtractBearerToken(r)

	if user, err := s.tokens.Get(accessToken); err == nil {
		return user.(*core.User), nil
	}

	user, err := s.Session.Login(r)
	if err != nil {
		return nil, err
	}

	_ = s.tokens.Set(accessToken, user)
	return user, nil
}
