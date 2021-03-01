package actions

import (
	"encoding/base64"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
)

func Tx(action core.Action, args ...interface{}) (string, error) {
	values := append([]interface{}{action}, args...)
	data, err := mtg.Encode(values...)
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
