package cfg

import (
	"github.com/spf13/viper"
)

const (
	hostKey = "api.Host"
)

func SetApiHost(host string) {
	viper.Set(hostKey, host)
}

func GetApiHost() string {
	return viper.GetString(hostKey)
}
