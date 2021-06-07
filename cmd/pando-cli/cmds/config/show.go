package config

import (
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewShowCmd() *cobra.Command {
	var (
		json bool
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "show local configures",
		RunE: func(cmd *cobra.Command, args []string) error {
			if json {
				d := jsoniter.NewEncoder(cmd.OutOrStdout())
				d.SetIndent("", "  ")
				return d.Encode(viper.AllSettings())
			}

			filename := viper.ConfigFileUsed()
			file, err := os.Open(filename)
			if err != nil {
				return err
			}

			_, err = io.Copy(cmd.OutOrStdout(), file)
			return err
		},
	}

	cmd.Flags().BoolVar(&json, "json", false, "display in json format")

	return cmd
}
