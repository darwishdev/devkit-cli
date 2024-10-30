package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getNewCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	apiCmd := c.getApiCmd()
	domainCmd := c.getDomainCmd()
	newCmd.AddCommand(apiCmd)
	newCmd.AddCommand(domainCmd)
	return newCmd

}
