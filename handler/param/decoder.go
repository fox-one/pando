package param

import (
	"reflect"

	"github.com/gorilla/schema"
	"github.com/shopspring/decimal"
)

var globalDecoder = queryDecoder()

func queryDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter(decimal.Decimal{}, convertDecimal)

	return decoder
}

func convertDecimal(s string) reflect.Value {
	if d, err := decimal.NewFromString(s); err == nil {
		return reflect.ValueOf(d)
	}

	return reflect.Value{}
}
