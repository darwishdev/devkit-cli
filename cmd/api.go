package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getApiCmd() *cobra.Command {
	apiCmd := &cobra.Command{
		Use:   "api [app-name]",
		Short: "Create a new API application",
		Long:  `Create a new API application with the specified parameters.`,
		Args:  cobra.ExactArgs(1), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			c.newCmd.NewApi(cmd, args)
		},
	}

	apiCmd.Flags().StringP("org", "o", "", "Github Org that you want to fork the app to , if not passed will use the git user from the config under ~/.config/devkitcli/devkit")

	apiCmd.Flags().StringP("buf-user", "b", "", "Buf.build username , it's used to push the api docs to this repo via buf cli later on , if not passed will use the git user from the config under ~/.config/devkitcli/devkit")

	return apiCmd
}
