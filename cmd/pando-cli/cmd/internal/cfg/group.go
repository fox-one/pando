package cfg

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/spf13/viper"
)

const (
	membersKey   = "group.members"
	thresholdKey = "group.threshold"
	verifyKey    = "group.verify"
)

func SetGroupMembers(members []string) {
	viper.Set(membersKey, members)
}

func GetGroupMembers() []string {
	return viper.GetStringSlice(membersKey)
}

func SetGroupThreshold(threshold int) {
	viper.Set(thresholdKey, threshold)
}

func GetGroupThreshold() int {
	return viper.GetInt(thresholdKey)
}

func SetGroupVerify(verify ed25519.PublicKey) {
	viper.Set(verifyKey, base64.StdEncoding.EncodeToString(verify))
}

func GetGroupVerify() ed25519.PublicKey {
	key := viper.GetString(verifyKey)
	b, _ := base64.StdEncoding.DecodeString(key)
	return b
}
