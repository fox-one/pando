package cfg

import (
	"github.com/spf13/viper"
)

func Get(key string) string {
	return viper.GetString(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func Save() error {
	return viper.WriteConfig()
}
