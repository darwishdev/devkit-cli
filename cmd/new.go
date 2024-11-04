package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getNewCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new [command]",
		Short: "Create a new resource",
		Long: `The "new" command is used to create various types of resources.
		It requires a subcommand to specify the type of resource to create. For example:

		new api       # Creates a new API
		new domain    # Creates a new domain
		new feature   # Creates a new feature
		new endpoint  # Creates a new endpoint

		Each subcommand has its own options and usage instructions.`, Run: func(cmd *cobra.Command, args []string) {
		},
	}

	apiCmd := c.getApiCmd()
	domainCmd := c.getDomainCmd()
	featureCmd := c.geFeatureCmd()
	endpointCmd := c.getEndpointCmd()
	newCmd.AddCommand(apiCmd)
	newCmd.AddCommand(domainCmd)
	newCmd.AddCommand(featureCmd)
	newCmd.AddCommand(endpointCmd)
	return newCmd

}
