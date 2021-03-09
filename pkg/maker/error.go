package maker

const (
	// 标记成可以退款
	FlagRefund = 1 << iota // 1
	// 标记成骚扰转账，不退款
	FlagNoisy = 1 << iota // 2
)

type Error struct {
	Msg  string
	Flag int
}

func (e Error) Error() string {
	return e.Msg
}

func Require(condition bool, msg string, flags ...int) error {
	if condition {
		return nil
	}

	var flag int
	for _, v := range flags {
		flag = flag | v
	}

	return Error{
		Msg:  msg,
		Flag: flag,
	}
}

func WithFlag(err error, flag int) error {
	e, ok := err.(Error)
	if !ok {
		return err
	}

	e.Flag = e.Flag | flag
	return e
}

func ShouldRefund(flag int) bool {
	return flag&FlagRefund > 0 && flag&FlagNoisy == 0
}
