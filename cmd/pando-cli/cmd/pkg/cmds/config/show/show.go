package show

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := viper.ConfigFileUsed()
			file, err := os.Open(filename)
			if err != nil {
				return err
			}

			_, err = io.Copy(cmd.OutOrStdout(), file)
			return err
		},
	}
}
