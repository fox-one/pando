package api

import (
	"fmt"
)

type Error struct {
	Code int64  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}
