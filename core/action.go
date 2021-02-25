package core

import (
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/pando/pkg/mtg"
)

//go:generate stringer -type Action -trimprefix Action

type Action int

const (
	ActionSys Action = iota + 0
	ActionSysWithdraw
	ActionSysVote
)

const (
	ActionProposal Action = iota + 10
	ActionProposalInit
	ActionProposalVote
)

const (
	ActionCat Action = iota + 20
	ActionCatInit
	ActionCatSupply
	ActionCatEdit
	ActionCatFold
)

const (
	ActionVat Action = iota + 30
	ActionVatInit
	ActionVatFrob
)

const (
	ActionFlip Action = iota + 40
	ActionFlipKick
	ActionFlipTend
	ActionFlipDent
	ActionFlipDeal
	ActionFlipOpt
)

const (
	ActionOracle Action = iota + 50
	ActionOracleFeed
)

func (i Action) MarshalBinary() (data []byte, err error) {
	return mtg.BitInt(i).MarshalBinary()
}

func (i *Action) UnmarshalBinary(data []byte) error {
	var b mtg.BitInt
	if err := b.UnmarshalBinary(data); err != nil {
		return err
	}

	*i = Action(b)
	return nil
}

type TransferAction struct {
	ID     string `json:"id,omitempty"`
	Source string `json:"s,omitempty"`
}

func (action TransferAction) Encode() string {
	b, _ := json.Marshal(action)
	return base64.StdEncoding.EncodeToString(b)
}
