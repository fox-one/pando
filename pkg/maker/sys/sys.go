package sys

import (
	"github.com/fox-one/pando/pkg/maker"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "Sys/"+msg)
}
