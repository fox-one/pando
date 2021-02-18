package core

import (
	"encoding/base64"
	"encoding/json"
)

// Vat Actions
const (
	_            = iota + 0
	ActionVatAdd // add collateral
	ActionVatFold
	ActionVatGov
	ActionVatOpen
	ActionVatFrob
)

// Flip Actions
const (
	_ = iota + 10
	ActionFlipGov
	ActionFlipKick
	ActionFlipTend
	ActionFlipDent
	ActionFlipDeal
)

type TransferAction struct {
	Module string `json:"m,omitempty"`
	ID     string `json:"id,omitempty"`
	Source string `json:"s,omitempty"`
}

func (action TransferAction) Encode() string {
	b, _ := json.Marshal(action)
	return base64.StdEncoding.EncodeToString(b)
}
