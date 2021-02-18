package flip

import (
	"time"

	"github.com/fox-one/pando/pkg/maker/flip"
	"github.com/fox-one/pando/pkg/number"
)

func getOpt() flip.Option {
	return flip.Option{
		Beg: number.Decimal("1.05"),
		TTL: time.Minute * 30,
		Tau: 3 * time.Hour,
	}
}
