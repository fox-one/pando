package dirtoracle

import (
	"fmt"

	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

type (
	PriceData struct {
		Timestamp int64           `json:"t,omitempty"`
		AssetID   string          `json:"a,omitempty"`
		Price     decimal.Decimal `json:"p,omitempty"`
		Signature *CosiSignature  `json:"s,omitempty"`
	}
)

func (p PriceData) Payload() []byte {
	return []byte(fmt.Sprintf("%d:%s:%v", p.Timestamp, p.AssetID, p.Price))
}

func (p *PriceData) MarshalBinary() (data []byte, err error) {
	asset, err := uuid.FromString(p.AssetID)
	if err != nil {
		return nil, err
	}
	return mtg.Encode(p.Timestamp, asset, p.Price, p.Signature)
}

func (p *PriceData) UnmarshalBinary(data []byte) error {
	var (
		timestamp int64
		price     decimal.Decimal
		signature CosiSignature
		asset     uuid.UUID
	)
	_, err := mtg.Scan(data, &timestamp, &asset, &price, &signature)
	if err != nil {
		return err
	}
	p.Timestamp = timestamp
	p.AssetID = asset.String()
	p.Price = price
	p.Signature = &signature
	return nil
}
