package actions

import (
	"encoding/base64"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/spf13/cobra"
)

func Build(cmd *cobra.Command, values ...interface{}) (string, error) {
	body, err := mtg.Encode(values...)
	if err != nil {
		return "", err
	}

	user, _ := uuid.FromString(cfg.GetAuthMixinID())
	follow, _ := uuid.FromString(uuid.New())

	cmd.Println("tx follow id:", follow)

	action := core.TransactionAction{
		UserID:   user.Bytes(),
		FollowID: follow.Bytes(),
		Body:     body,
	}

	data, err := action.Encode()
	if err != nil {
		return "", err
	}

	key := mixin.GenerateEd25519Key()
	pub := cfg.GetGroupVerify()
	data, err = mtg.Encrypt(data, key, pub)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
