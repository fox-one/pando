package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"

	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

type RawMessage []byte

func (r RawMessage) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(r))
	copy(b, r)
	return b, nil
}

func (r *RawMessage) UnmarshalBinary(data []byte) error {
	*r = append((*r)[0:0], data...)
	return nil
}

func (r RawMessage) Value() (driver.Value, error) {
	return []byte(r), nil
}

func (r *RawMessage) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("incompatible type for RawMessage")
	}

	*r = append((*r)[0:0], source...)
	return nil
}

type BitInt int

func (i BitInt) MarshalBinary() ([]byte, error) {
	b := byte(i)
	return []byte{b}, nil
}

func (i *BitInt) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return errors.New("BitInt must be 1 byte")
	}

	v := int(data[0])
	*i = BitInt(v)
	return nil
}

const (
	TypeString  = "str"     // string
	TypeBit     = "bit"     // BitInt
	TypeUUID    = "uuid"    // UUID
	TypeDecimal = "decimal" // Decimal
	TypeInt     = "int"     // Int
)

func UUID(v string) Value {
	return Value{
		typ: TypeUUID,
		raw: v,
	}
}

func Decimal(v string) Value {
	return Value{
		typ: TypeDecimal,
		raw: v,
	}
}

type Value struct {
	typ string
	raw string
}

func (v Value) MarshalBinary() ([]byte, error) {
	b, err := castValue(v.typ, v.raw)
	if err != nil {
		return nil, err
	}

	data, err := mtg.Encode(b)
	if err != nil {
		return nil, err
	}

	return data[1:], nil
}

func castValue(typ, value string) (interface{}, error) {
	var v interface{}

	switch typ {
	case TypeString:
		v = value
	case TypeBit:
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		v = BitInt(i)
	case TypeUUID:
		u, err := uuid.FromString(value)
		if err != nil {
			return nil, err
		}

		v = u
	case TypeDecimal:
		d, err := decimal.NewFromString(value)
		if err != nil {
			return nil, err
		}

		v = d
	case TypeInt:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}

		v = i
	default:
		return nil, fmt.Errorf("unknown value type %q", typ)
	}

	return v, nil
}

func EncodeWithTypes(typeValues ...string) ([]byte, error) {
	var values []interface{}

	for idx := 0; idx < len(typeValues)-1; idx += 2 {
		typ, value := typeValues[idx], typeValues[idx+1]
		v, err := castValue(typ, value)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return mtg.Encode(values...)
}
