package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type (
	Config struct {
		DB    db.Config `json:"db"`
		Dapp  Dapp      `json:"dapp"`
		Group Group     `json:"group,omitempty"`
	}

	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret,omitempty"`
		Pin          string `json:"pin"`
	}

	Member struct {
		ClientID string `json:"client_id,omitempty"`
		// 节点共享的用户验证签名的公钥
		VerifyKey string `json:"verify_key,omitempty"`
	}

	Group struct {
		// 节点共享的用户解密的私钥
		PrivateKey string   `json:"private_key,omitempty"`
		Members    []Member `json:"members,omitempty"`
		Threshold  uint8    `json:"threshold,omitempty"`
	}
)

// Viperion load config by viper
func Viperion(cfgFile string) (*Config, error) {
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigType("yaml")
		v.SetConfigName("config")
		v.AddConfigPath("/etc/pando/server")
		v.AddConfigPath("$HOME/.pando/server")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	b, err := jsoniter.Marshal(v.AllSettings())
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := jsoniter.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
