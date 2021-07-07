package config

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type (
	Config struct {
		DB      db.Config `json:"db"`
		Dapp    Dapp      `json:"dapp"`
		Group   Group     `json:"group"`
		Gas     Gas       `json:"gas"`
		Flip    Flip      `json:"flip"`
		Vault   Vault     `json:"vault"`
		I18n    I18n      `json:"i18n"`
		DataDog DataDog   `json:"data_dog"`
	}

	Dapp struct {
		mixin.Keystore
		Pin string `json:"pin"`
	}

	Vote struct {
		Asset  string          `json:"asset,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`
	}

	Group struct {
		// 节点管理员 mixin id
		Admins []string `json:"admins,omitempty"`
		// 节点共享的用户解密的私钥
		PrivateKey string   `json:"private_key,omitempty"`
		Members    []string `json:"members,omitempty"`
		Threshold  uint8    `json:"threshold,omitempty"`
	}

	Gas struct {
		AssetID string          `json:"asset_id"`
		Amount  decimal.Decimal `json:"amount"`
	}

	Flip struct {
		DetailPage string `json:"detail_page,omitempty"`
	}

	Vault struct {
		DetailPage string `json:"detail_page,omitempty"`
	}

	I18n struct {
		Path string `json:"path,omitempty"`
		// default language
		Language string `json:"language,omitempty"`
	}

	DataDog struct {
		ConversationID string `json:"conversation_id,omitempty"`
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

	// disable read write separation
	cfg.DB.ReadHost = ""

	defaultGas(&cfg)
	defaultI18n(&cfg)

	return &cfg, nil
}

func defaultGas(cfg *Config) {
	if cfg.Gas.AssetID == "" {
		cfg.Gas.AssetID = "965e5c6e-434c-3fa9-b780-c50f43cd955c" // cnb
	}

	if cfg.Gas.Amount.IsZero() {
		cfg.Gas.Amount = decimal.New(1, -8)
	}
}

func defaultI18n(cfg *Config) {
	if cfg.I18n.Path == "" {
		cfg.I18n.Path = "./assets/i18n"
	}

	if cfg.I18n.Language == "" {
		cfg.I18n.Language = language.English.String()
	}
}
