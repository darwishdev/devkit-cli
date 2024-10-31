package cmd

import (
	"github.com/spf13/cobra"
)

func (c *Command) getSstorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage --files-path [files_path] --icons-path [icons_path]",
		Short: "Seed Supabase storage with files and icons.",
		Long: `The 'storage' command seeds Supabase storage with files from the specified 
		'files-path' and icons from SVG files in the 'icons-path'. It creates buckets based on 
		the subfolder structure in 'files-path' and uploads the files to their respective buckets. 
		It also inserts the SVG content from 'icons-path' into the 'icons' table in your database.`, Run: func(cmd *cobra.Command, args []string) {
			c.seedCmd.StorageSeed(cmd.Flags())
		},
	}
	cmd.Flags().StringP("files-path", "f", "", "Path to the folder that holds your images and files (required)")
	cmd.Flags().StringP("icons-path", "i", "", "Path to the folder that holds you svg icons")
	return cmd
}
