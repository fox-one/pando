package number

import (
	"database/sql/driver"
	"fmt"
	"net/url"

	"github.com/shopspring/decimal"
)

type Values map[string]decimal.Decimal

func (v Values) Get(key string) decimal.Decimal {
	return v[key]
}

func (v Values) Set(key string, value decimal.Decimal) {
	v[key] = value
}

func (v Values) Add(key string, value decimal.Decimal) {
	v[key] = v[key].Add(value)
}

func (v Values) Del(key string) {
	delete(v, key)
}

func (v Values) Merge(other Values) {
	for key, value := range other {
		v.Add(key, value)
	}
}

func (v Values) Encode() string {
	query := make(url.Values)
	for key, value := range v {
		query.Set(key, value.String())
	}

	return query.Encode()
}

func (v Values) String() string {
	return v.Encode()
}

func (v Values) Sum() decimal.Decimal {
	sum := decimal.Zero
	for _, value := range v {
		sum = sum.Add(value)
	}

	return sum
}

func (v Values) Mul(m decimal.Decimal) Values {
	values := Values{}
	for key, value := range v {
		values[key] = value.Mul(m)
	}

	return values
}

func (v Values) Div(d decimal.Decimal) Values {
	values := Values{}
	for key, value := range v {
		values[key] = value.Div(d)
	}

	return values
}

func (v Values) Truncate(p int32) Values {
	values := Values{}
	for key, value := range v {
		values[key] = value.Truncate(p)
	}

	return values
}

func (v Values) Keys() []string {
	var keys []string
	for key := range v {
		keys = append(keys, key)
	}

	return keys
}

func ParseValues(input string) (Values, error) {
	v := make(Values)

	query, err := url.ParseQuery(input)
	if err != nil {
		return v, err
	}

	for key, values := range query {
		for _, value := range values {
			if d := Decimal(value); !d.IsZero() {
				v.Add(key, d)
			}
		}
	}

	return v, nil
}

func (v Values) Value() (driver.Value, error) {
	return v.Encode(), nil
}

func (v *Values) Scan(value interface{}) error {
	var input string

	if value != nil {
		switch t := value.(type) {
		case string:
			input = t
		case []byte:
			input = string(t)
		default:
			input = fmt.Sprintf("%v", t)
		}
	}

	var err error
	*v, err = ParseValues(input)
	return err
}
