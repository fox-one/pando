/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/auth"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/cat"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/config"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/oracle"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/proposal"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/sys"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/use"
	"github.com/fox-one/pando/cmd/pando-cli/cmd/pkg/cmds/vat"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 如果指定了 keystore 文件，那么付款的时候会直接用这个账号付款
	// keystore 里面需要指定 pin
	keystoreFile string
)

var rootCmd = &cobra.Command{
	Use:   "pd <command>",
	Short: "Pando's office command line tool",
}

func Execute() {
	version := os.Getenv("PANDO_VERSION")
	commit := os.Getenv("PANDO_COMMIT")
	rootCmd.Version = fmt.Sprintf("%s(%s)", version, commit)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initKeystore)

	rootCmd.AddCommand(use.NewCmd())
	rootCmd.AddCommand(config.NewCmd())
	rootCmd.AddCommand(auth.NewCmd())
	rootCmd.AddCommand(proposal.NewCmd())
	rootCmd.AddCommand(sys.NewCmd())
	rootCmd.AddCommand(cat.NewCmd())
	rootCmd.AddCommand(vat.NewCmd())
	rootCmd.AddCommand(oracle.NewCmd())

	rootCmd.PersistentFlags().StringVar(&keystoreFile, "keystore", "", "keystore filename")
}

func initConfig() {
	dir, _ := homedir.Expand("~/.pando")
	_ = os.MkdirAll(dir, os.ModePerm)

	filename := path.Join(dir, "conf.yaml")
	viper.SetConfigFile(filename)

	_ = viper.SafeWriteConfigAs(filename)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("read config file failed: %v", err)
		}
	}
}

func initKeystore() {
	if keystoreFile == "" {
		return
	}

	data, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		log.Fatalf("read keystore file failed: %v", err)
	}

	var store pay.Keystore
	if err := json.Unmarshal(data, &store); err != nil {
		log.Fatalf("decode keystore failed: %v", err)
	}

	pay.UseKeystore(&store)
}
