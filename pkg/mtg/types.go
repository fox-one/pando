package mtg

import (
	"database/sql/driver"
	"errors"
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
	v := 0

	if len(data) > 0 {
		v = int(data[0])
	}

	*i = BitInt(v)
	return nil
}
