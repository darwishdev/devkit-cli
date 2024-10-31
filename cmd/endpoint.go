package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getEndpointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint [endpoint_name] --domain [domain_name] --feature [domain_name]",
		Short: "Create a new feature within a domain",
		Long: `The 'feature' command generates a new feature within a specified domain.
		It creates the necessary files (adapter, repo, usecase, proto, query, api) 
		and populates them with boilerplate code for the given feature name. 
		The --domain flag is required to specify the domain where the feature belongs.`,
		Args: cobra.ExactArgs(1), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			c.newCmd.NewEndpoint(args, cmd.Flags())
		},
	}
	cmd.Flags().StringP("domain", "d", "", "the domain name for the created feature , must be existed domain under  your app directory")
	cmd.Flags().StringP("feature", "f", "", "the feature name for the created feature , must be existed domain under  your app directory")
	cmd.Flags().BoolP("allow-get", "g", false, "if enabled this endpoint will be marked as no side effect and the rest endpoint will be of type GET(optional)")
	cmd.Flags().BoolP("empty-response", "v", false, "if enabled this endpoint will return empty responnse (optional)")
	cmd.Flags().BoolP("empty-request", "e", false, "if enabled this endpoint will accept empty request (optional)")
	cmd.Flags().BoolP("list", "l", false, "if enabled then this endpoint will return set of records if disable will use single record except the endpoints with term 'list' in their name (optional)")
	return cmd
}
