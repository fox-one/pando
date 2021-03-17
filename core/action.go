package core

import (
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/msgpack"
	"github.com/fox-one/pando/pkg/mtg/types"
)

//go:generate stringer -type Action -trimprefix Action

type Action int

const (
	ActionSys Action = iota + 0
	ActionSysWithdraw
)

const (
	ActionProposal Action = iota + 10
	ActionProposalMake
	ActionProposalShout
	ActionProposalVote
)

const (
	ActionCat Action = iota + 20
	ActionCatCreate
	ActionCatSupply
	ActionCatEdit
	ActionCatFold
)

const (
	ActionVat Action = iota + 30
	ActionVatOpen
	ActionVatDeposit
	ActionVatWithdraw
	ActionVatPayback
	ActionVatGenerate
)

const (
	ActionFlip Action = iota + 40
	ActionFlipKick
	ActionFlipBid
	ActionFlipDeal
)

const (
	ActionOracle Action = iota + 50
	ActionOraclePoke
	ActionOracleStep
)

func (i Action) MarshalBinary() (data []byte, err error) {
	return types.BitInt(i).MarshalBinary()
}

func (i *Action) UnmarshalBinary(data []byte) error {
	var b types.BitInt
	if err := b.UnmarshalBinary(data); err != nil {
		return err
	}

	*i = Action(b)
	return nil
}

type TransactionAction struct {
	FollowID []byte `msgpack:"f,omitempty"`
	Body     []byte `msgpack:"b,omitempty"`
}

func (action TransactionAction) Encode() ([]byte, error) {
	return msgpack.Marshal(action)
}

func DecodeTransactionAction(b []byte) (*TransactionAction, error) {
	var action TransactionAction
	if err := msgpack.Unmarshal(b, &action); err != nil {
		return nil, err
	}

	return &action, nil
}

type TransferAction struct {
	ID     string `json:"id,omitempty"`
	Source string `json:"s,omitempty"`
}

func (action TransferAction) Encode() string {
	b, _ := json.Marshal(action)
	return base64.StdEncoding.EncodeToString(b)
}
