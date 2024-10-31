package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) geFeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature [feature_name] --domain [domain_name]",
		Short: "Create a new feature within a domain",
		Long: `The 'feature' command generates a new feature within a specified domain.
		It creates the necessary files (adapter, repo, usecase, proto, query, api) 
		and populates them with boilerplate code for the given feature name. 
		The --domain flag is required to specify the domain where the feature belongs.`,
		Args: cobra.ExactArgs(1), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			c.newCmd.NewFeature(args, cmd.Flags())
		},
	}
	cmd.Flags().StringP("domain", "d", "", "the domain name for the created feature , must be existed domain under  your app directory")
	cmd.MarkFlagRequired("domain") // Make file-path flag required
	return cmd
}
