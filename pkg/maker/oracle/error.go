package oracle

import (
	"github.com/fox-one/pando/pkg/maker"
)

const module = "oracle"

var (
	ErrOracleOverdue = maker.RegisterErr(1201, module, "overdue")
)
