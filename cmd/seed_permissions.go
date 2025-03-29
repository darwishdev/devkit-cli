package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getPermissionsSeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions [schema_name] --file-path [excel_file_path]",
		Short: "Generate SQL insert statements from an Excel file.",
		Long: `The 'seed' command automates the process of seeding your database tables with data 
		from an Excel file. It generates SQL insert statements based on the sheet names and column 
		headers in the Excel file. It can also create users in Supabase Auth if the Excel file 
		contains a "users" sheet. The command requires the schema name and the file path to the 
		Excel file.`,
		Args: cobra.ExactArgs(0), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			c.seedCmd.SeedPermissions()
		},
	}
	return cmd
}
