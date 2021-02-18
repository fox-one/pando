package number

import (
	"math/cmplx"

	"github.com/shopspring/decimal"
)

var (
	One = decimal.NewFromInt(1)
)

type Number = decimal.Decimal

func Decimal(v string) decimal.Decimal {
	d, _ := decimal.NewFromString(v)
	return d
}

func Pow(a, b decimal.Decimal) decimal.Decimal {
	f1, _ := a.Float64()
	f2, _ := b.Float64()

	r := cmplx.Pow(
		complex(f1, 0),
		complex(f2, 0),
	)

	return decimal.NewFromFloat(real(r))
}

func SqrtN(d decimal.Decimal, n int64) decimal.Decimal {
	return Pow(d, One.Div(decimal.NewFromInt(n)))
}

func Sqrt(d decimal.Decimal) decimal.Decimal {
	return SqrtN(d, 2)
}

func Ceil(d decimal.Decimal, precision int32) decimal.Decimal {
	return d.Shift(precision).Ceil().Shift(-precision)
}
