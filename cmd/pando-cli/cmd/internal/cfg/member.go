package cfg

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/spf13/viper"
)

const (
	memberClientKey = "member.client_id"
	memberSignKey   = "member.sign_key"
)

func GetMember() (string, ed25519.PrivateKey) {
	clientID := viper.GetString(memberClientKey)
	signKey := viper.GetString(memberSignKey)
	b, _ := base64.StdEncoding.DecodeString(signKey)
	return clientID, b
}
