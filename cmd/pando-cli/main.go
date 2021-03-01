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
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/fox-one/pando/cmd/pando-cli/cmds/auth"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/cat"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/config"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/oracle"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/pay"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/proposal"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/sys"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/tx"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/use"
	"github.com/fox-one/pando/cmd/pando-cli/cmds/vat"
	"github.com/fox-one/pando/cmd/pando-cli/internal/call"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 如果指定了 keystore 文件，那么付款的时候会直接用这个账号付款
	// keystore 里面需要指定 pin
	keystoreFile string
	debug        bool
)

var rootCmd = &cobra.Command{
	Use:   "pd <command>",
	Short: "Pando's office command line tool",
}

func main() {
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
	rootCmd.AddCommand(tx.NewCmd())

	rootCmd.PersistentFlags().StringVar(&keystoreFile, "keystore", "", "keystore filename")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
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

	call.SetDebug(debug)
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
