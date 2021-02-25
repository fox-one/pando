package maker

import (
	"fmt"
)

type Error struct {
	Msg string `json:"msg,omitempty"`
}

func (e Error) Error() string {
	return e.Msg
}

func Require(condition bool, format string, args ...interface{}) error {
	if condition {
		return nil
	}

	return Error{Msg: fmt.Sprintf(format, args...)}
}
