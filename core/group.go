package core

import (
	"crypto/ed25519"
	"errors"

	"github.com/fox-one/pando/pkg/mtg"
	"github.com/gofrs/uuid"
)

type Member struct {
	ClientID  string
	Name      string
	VerifyKey ed25519.PublicKey
}

func DecodeMemberAction(message []byte, members []*Member) (*Member, []byte, error) {
	body, sig, err := mtg.Unpack(message)
	if err != nil {
		return nil, nil, err
	}

	var id uuid.UUID
	content, err := mtg.Scan(body, &id)
	if err != nil {
		return nil, nil, err
	}

	for _, member := range members {
		if member.ClientID != id.String() {
			continue
		}

		if !mtg.Verify(body, sig, member.VerifyKey) {
			return nil, nil, errors.New("verify sig failed")
		}

		return member, content, nil
	}

	return nil, nil, errors.New("member not found")
}
