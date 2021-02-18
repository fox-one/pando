package maker

import (
	"errors"
	"fmt"
)

type Err struct {
	Code   int64
	Module string
	Detail string
}

func (err Err) Msg() string {
	return err.Msg()
}

func (err Err) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg())
}

func RegisterErr(code int64, module, msg string) error {
	err := Err{
		Code:   code,
		Module: module,
		Detail: msg,
	}

	msg, ok := ErrorMsg(code)
	if ok {
		panic(fmt.Errorf("duplicated error code %s", code))
	}

	_msgs[code] = err.Msg()
	return err
}

var _msgs = map[int64]string{}

func ErrorMsg(code int64) (string, bool) {
	msg, ok := _msgs[code]
	return msg, ok
}

func Unwrap(err error) (code int64, msg string, ok bool) {
	var e Err
	if !errors.As(err, &e) {
		return
	}

	return e.Code, e.Msg(), true
}
