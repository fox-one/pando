package proposal

import (
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/number"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Withdraw struct {
	Opponent string `json:"opponent,omitempty"`
	Asset    string `json:"asset,omitempty"`
	Amount   string `json:"amount,omitempty"`
}

func (w Withdraw) MarshalBinary() (data []byte, err error) {
	opponent, err := uuid.FromString(w.Opponent)
	if err != nil {
		return nil, err
	}

	asset, err := uuid.FromString(w.Asset)
	if err != nil {
		return nil, err
	}

	amount := number.Decimal(w.Amount)
	return mtg.Encode(opponent, asset, amount)
}

func (w *Withdraw) UnmarshalBinary(data []byte) error {
	var (
		opponent, asset uuid.UUID
		amount          decimal.Decimal
	)

	if _, err := mtg.Scan(data, &opponent, &asset, &amount); err != nil {
		return err
	}

	w.Opponent = opponent.String()
	w.Asset = asset.String()
	w.Amount = amount.String()
	return nil
}
