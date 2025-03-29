package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getSeedSuperUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "super-user --email [user email] --phone [user phone] --password [user-password] --name [user email]",
		Short: "create the super user with super admin priviladged for inital setup",
		Long: `The 'seed' command automates the process of seeding your database tables with data 
		from an Excel file. It generates SQL insert statements based on the sheet names and column 
		headers in the Excel file. It can also create users in Supabase Auth if the Excel file 
		contains a "users" sheet. The command requires the schema name and the file path to the 
		Excel file.`,
		Args: cobra.ExactArgs(0), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			c.seedCmd.SeedSuperUser(cmd.Flags())
		},
	}
	cmd.Flags().StringP("email", "e", "", "Path to the data file (required)")
	cmd.Flags().StringP("phone", "t", "", "Path to the output SQL file (optional)")
	cmd.Flags().StringP("password", "p", "123456", "Execute the generated SQL (optional)")
	cmd.Flags().StringP("name", "n", "", "Execute the generated SQL (optional)")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("name")
	return cmd
}
