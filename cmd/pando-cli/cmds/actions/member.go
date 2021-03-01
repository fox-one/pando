package actions

import (
	"encoding/base64"

	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
)

func Member(args ...interface{}) (string, error) {
	clientID, signKey := cfg.GetMember()

	values := []interface{}{types.UUID(clientID)}
	values = append(values, args...)

	data, err := mtg.Encode(values...)
	if err != nil {
		return "", err
	}

	sig := mtg.Sign(data, signKey)
	data = mtg.Pack(data, sig)

	return base64.StdEncoding.EncodeToString(data), nil
}

func MakeProposal(action core.Action, args ...interface{}) (string, error) {
	values := []interface{}{
		core.ActionProposalMake,
		types.UUID(uuid.New()),
		action,
	}

	values = append(values, args...)
	return Member(values...)
}
