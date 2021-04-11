package config

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type (
	Config struct {
		DB    db.Config `json:"db"`
		Dapp  Dapp      `json:"dapp"`
		Group Group     `json:"group,omitempty"`
		I18n  I18n      `json:"i18n,omitempty"`
	}

	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret,omitempty"`
		Pin          string `json:"pin"`
	}

	Group struct {
		// 节点共享的用户解密的私钥
		PublicKey string   `json:"public_key,omitempty"`
		Members   []string `json:"members,omitempty"`
		Threshold uint8    `json:"threshold,omitempty"`
	}

	I18n struct {
		Path string `json:"path,omitempty"`
		// default language
		Language string `json:"language,omitempty"`
	}
)

func viperion(cfgFile, embed string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else if embed != "" {
		b, err := base64.StdEncoding.DecodeString(embed)
		if err != nil {
			return nil, fmt.Errorf("decode embed config failed: %w", err)
		}

		return v, v.ReadConfig(bytes.NewReader(b))
	}

	return v, v.ReadInConfig()
}

// Viperion load config by viper
func Viperion(cfgFile, embed string) (*Config, error) {
	v, err := viperion(cfgFile, embed)
	if err != nil {
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

	defaultI18n(&cfg)
	return &cfg, nil
}

func defaultI18n(cfg *Config) {
	if cfg.I18n.Path == "" {
		cfg.I18n.Path = "./assets/i18n"
	}

	if cfg.I18n.Language == "" {
		cfg.I18n.Language = language.English.String()
	}
}
