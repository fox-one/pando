package cfg

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

func Get(key string) string {
	v := viper.GetString(key)
	if v == "" {
		if sub := viper.Sub(key); sub != nil {
			v, _ = jsoniter.MarshalToString(sub.AllSettings())
		}
	}

	return v
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func Save() error {
	return viper.WriteConfig()
}
