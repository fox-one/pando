package cfg

import (
	"github.com/spf13/viper"
)

const (
	clientKey  = "auth.client_id"
	tokenKey   = "auth.token"
	mixinIDKey = "auth.mixin_id"
)

func SetAuthClient(clientID string) {
	viper.Set(clientKey, clientID)
}

func GetAuthClient() string {
	return viper.GetString(clientKey)
}

func SetAuthToken(token string) {
	viper.Set(tokenKey, token)
}

func GetAuthToken() string {
	return viper.GetString(tokenKey)
}

func SetAuthMixinID(id string) {
	viper.Set(mixinIDKey, id)
}

func GetAuthMixinID() string {
	return viper.GetString(mixinIDKey)
}
