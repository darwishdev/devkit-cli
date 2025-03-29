package cmd

import "github.com/spf13/cobra"

func (c *Command) getNewSeedPaginator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "paginate-query --query [the query string] --primary-key [string]",
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
	cmd.Flags().StringP("query", "q", "", "the query string that need to be paginated")
	cmd.Flags().StringP("primary-key", "p", "", "the primary key of the table to be able to use key set pagination")
	cmd.MarkFlagRequired("query")
	return cmd

}
