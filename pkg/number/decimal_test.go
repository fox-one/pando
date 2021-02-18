package number

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/shopspring/decimal"
)

func TestCeil(t *testing.T) {
	data := map[string]string{
		"0.10304":     "0.11",
		"0.100000001": "0.11",
		"0.108":       "0.11",
	}

	for k, v := range data {
		t.Run(k, func(t *testing.T) {
			_k := Decimal(k)
			c := Ceil(Decimal(k), 2)
			t.Log(k, c, _k.Round(2))
			assert.Equal(t, v, c.String(), "should be ceil")
		})
	}
}

func TestSqrtN(t *testing.T) {
	type args struct {
		d decimal.Decimal
		n int64
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "sqrt1",
			args: args{
				d: Decimal("10"),
				n: 1,
			},
			want: Decimal("10"),
		},
		{
			name: "sqrt2",
			args: args{
				d: Decimal("4"),
				n: 2,
			},
			want: Decimal("2"),
		},
		{
			name: "sqrt4",
			args: args{
				d: Decimal("16"),
				n: 4,
			},
			want: Decimal("2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SqrtN(tt.args.d, tt.args.n); !got.Equal(tt.want) {
				t.Errorf("SqrtN() = %v, want %v", got, tt.want)
			}
		})
	}
}
