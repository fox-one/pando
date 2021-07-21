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
	// Proposal Actions
	ActionProposal Action = iota + 10
	// Make a new proposal
	ActionProposalMake
	// Call on other mtg members to vote for this proposal, mtg member only
	ActionProposalShout
	// Vote for this proposal, mtg member only
	ActionProposalVote
)

const (
	// Collateral Actions
	ActionCat Action = iota + 20
	// Create a new collateral type, Gov required
	ActionCatCreate
	// Supply Dai to this collateral type
	ActionCatSupply
	// Edit Collateral parameters, Gov required
	ActionCatEdit
	// Update Collateral's Rate
	ActionCatFold
	// ActionCatMove move supply from ont to another
	ActionCatMove
)

const (
	// Vault Actions
	ActionVat Action = iota + 30
	// Open a new Vault
	ActionVatOpen
	// Deposit gem into Vault
	ActionVatDeposit
	// Withdraw gem from Vault
	ActionVatWithdraw
	// Pay back dai for Vault
	ActionVatPayback
	// Generate dai from Vault
	ActionVatGenerate
)

const (
	ActionFlip Action = iota + 40
	// Launch an auction from an unsafe Vault by Keeper
	ActionFlipKick
	// Auction bid
	ActionFlipBid
	// Auction done
	ActionFlipDeal
)

const (
	// Oracle Actions
	ActionOracle Action = iota + 50
	// Create a new price
	ActionOracleCreate
	// Edit can modify Next, hop & Threshold
	ActionOracleEdit
	// Poke push the next price to be current
	ActionOraclePoke
	// Rely add a new oracle feed
	ActionOracleRely
	// Deny remove a existed oracle feed
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
