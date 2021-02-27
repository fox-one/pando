package codes

import (
	"strconv"

	"github.com/twitchtv/twirp"
)

const (
	CustomCodeKey = "custom_code"

	InvalidArguments = 10002
)

func With(err error, code int) error {
	twerr, ok := err.(twirp.Error)
	if !ok {
		twerr = twirp.InternalErrorWith(err)
	}

	return twerr.WithMeta(CustomCodeKey, strconv.Itoa(code))
}

func Get(code twirp.ErrorCode) int {
	switch code {
	case twirp.InvalidArgument:
		return InvalidArguments
	default:
		return twirp.ServerHTTPStatusFromErrorCode(code)
	}
}
