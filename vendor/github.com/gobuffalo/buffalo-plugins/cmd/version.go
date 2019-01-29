package cmd

import (
	"fmt"

	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of <%= opts.ShortName %>",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("buffalo-plugins", plugins.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
