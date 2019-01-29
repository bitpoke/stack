package cmd

import (
	"github.com/gobuffalo/buffalo-plugins/plugins/plugcmds"
	"github.com/spf13/cobra"
)

var Available = plugcmds.NewAvailable()

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "tools for working with buffalo plugins",
}

func init() {
	pluginsCmd.AddCommand(addCmd)
	pluginsCmd.AddCommand(listCmd)
	pluginsCmd.AddCommand(generateCmd)
	pluginsCmd.AddCommand(removeCmd)
	pluginsCmd.AddCommand(installCmd)

	Available.Add("generate", generateCmd)
	Available.Add("root", pluginsCmd)
	Available.Listen(Listen)
	Available.Mount(rootCmd)
}
