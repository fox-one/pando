package core

import (
	"crypto/ed25519"

	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
)

// System stores system information.
type System struct {
	Admins       []string
	ClientID     string
	ClientSecret string
	Members      []string
	Threshold    uint8
	VoteAsset    string
	VoteAmount   decimal.Decimal
	PrivateKey   ed25519.PrivateKey
	PublicKey    ed25519.PublicKey
	SignKey      ed25519.PrivateKey
	Version      string
}

func (s *System) IsMember(id string) bool {
	return govalidator.IsIn(id, s.Members...)
}

func (s *System) IsStaff(id string) bool {
	return govalidator.IsIn(id, s.Admins...)
}
