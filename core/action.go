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
	// ActionSys System Actions
	ActionSys Action = iota + 0
	// ActionSysWithdraw Withdraw Asset from multisig wallet, Gov required
	ActionSysWithdraw
	// ActionSysProperty set custom property
	ActionSysProperty
)

const (
	// ActionProposal Proposal Actions
	ActionProposal Action = iota + 10
	// ActionProposalMake Make a new proposal
	ActionProposalMake
	// ActionProposalShout Call on other mtg members to vote for this proposal, mtg member only
	ActionProposalShout
	// ActionProposalVote Vote for this proposal, mtg member only
	ActionProposalVote
)

const (
	// ActionCat Collateral Actions
	ActionCat Action = iota + 20
	// ActionCatCreate Create a new collateral type, Gov required
	ActionCatCreate
	// ActionCatSupply Supply Dai to this collateral type
	ActionCatSupply
	// ActionCatEdit Edit Collateral parameters, Gov required
	ActionCatEdit
	// ActionCatFold Update Collateral's Rate
	ActionCatFold
	// ActionCatMove ActionCatMove move supply from ont to another
	ActionCatMove
	// ActionCatGain withdraw profits from collateral
	ActionCatGain
	// ActionCatFill make up for the loss
	ActionCatFill
)

const (
	// ActionVat Vault Actions
	ActionVat Action = iota + 30
	// ActionVatOpen Open a new Vault
	ActionVatOpen
	// ActionVatDeposit Deposit gem into Vault
	ActionVatDeposit
	// ActionVatWithdraw Withdraw gem from Vault
	ActionVatWithdraw
	// ActionVatPayback Pay back dai for Vault
	ActionVatPayback
	// ActionVatGenerate Generate dai from Vault
	ActionVatGenerate
)

const (
	ActionFlip Action = iota + 40
	// ActionFlipKick Launch an auction from an unsafe Vault by Keeper
	ActionFlipKick
	// ActionFlipBid Auction bid
	ActionFlipBid
	// ActionFlipDeal Auction done
	ActionFlipDeal
)

const (
	// ActionOracle Oracle Actions
	ActionOracle Action = iota + 50
	// ActionOracleCreate Create a new price
	ActionOracleCreate
	// ActionOracleEdit Edit can modify Next, hop & Threshold
	ActionOracleEdit
	// ActionOraclePoke Poke push the next price to be current
	ActionOraclePoke
	// ActionOracleRely Rely add a new oracle feed
	ActionOracleRely
	// ActionOracleDeny Deny remove a existed oracle feed
	ActionOracleDeny
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

const (
	TransferSourceRefund = "Refund"
)

type TransferAction struct {
	ID     string `json:"id,omitempty"`
	Source string `json:"s,omitempty"`
}

func (action TransferAction) Encode() string {
	b, _ := json.Marshal(action)
	return base64.StdEncoding.EncodeToString(b)
}
