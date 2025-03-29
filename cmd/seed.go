package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getSeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed [schema_name] --file-path [excel_file_path]",
		Short: "Generate SQL insert statements from an Excel file.",
		Long: `The 'seed' command automates the process of seeding your database tables with data 
		from an Excel file. It generates SQL insert statements based on the sheet names and column 
		headers in the Excel file. It can also create users in Supabase Auth if the Excel file 
		contains a "users" sheet. The command requires the schema name and the file path to the 
		Excel file.`,
		Args: cobra.ExactArgs(1), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			c.seedCmd.NewSeed(args, cmd.Flags())
		},
	}
	cmd.Flags().StringP("file-path", "f", "", "Path to the data file (required)")
	cmd.Flags().StringP("out-file", "o", "", "Path to the output SQL file (optional)")
	cmd.Flags().BoolP("skip-supabase", "s", false, "Execute the generated SQL (optional)")
	cmd.Flags().BoolP("execute", "e", false, "Execute the generated SQL (optional)")
	cmd.MarkFlagRequired("file-path")

	permissionsSeed := c.getPermissionsSeedCmd()
	storageCmd := c.getSstorageCmd()
	superUserCmd := c.getSeedSuperUserCmd()
	paginatoreCmd := c.getNewSeedPaginator()
	cmd.AddCommand(storageCmd)
	cmd.AddCommand(permissionsSeed)
	cmd.AddCommand(superUserCmd)
	cmd.AddCommand(paginatoreCmd)
	return cmd
}
