package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	supaapigo "github.com/darwishdev/supaapi-go"
	"github.com/iancoleman/strcase"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	storage_go "github.com/supabase-community/storage-go"
)

// ... other code ...

// seedCmd represents the seed command
var storageCmd = &cobra.Command{
	Use:   "storage --files-path <files_path> --icons-path <icons_path>",
	Short: "Seed storage with files and icons",
	Long:  `Seed storage with files from the specified paths and icons from SVG files.`,
	Run: func(cmd *cobra.Command, args []string) {
		filesPath, _ := cmd.Flags().GetString("files-path")
		iconsPath, _ := cmd.Flags().GetString("icons-path")

		// Basic validation
		if filesPath == "" || iconsPath == "" {
			fmt.Println("Error: --files-path and --icons-path are required")
			os.Exit(1)
		}
		conf := appConfig.GetConfig()
		supaapi := supaapigo.NewSupaapi(supaapigo.SupaapiConfig{
			ProjectRef:     conf.DBProjectREF,
			Env:            supaapigo.DEV,
			Port:           54321,
			ServiceRoleKey: conf.SupabaseServiceRoleKey,
			ApiKey:         conf.SupabaseApiKey,
		})

		foldersMap := map[string][]string{}
		// 1. Create buckets and upload files
		err := filepath.Walk(filesPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && path != filesPath { // Exclude the root directory
				bucketName := filepath.Base(path) // Use the folder name as the bucket name
				_, err = supaapi.StorageClient.CreateBucket(bucketName, storage_go.BucketOptions{
					Public: true,
				})
				if err != nil && !strings.Contains(err.Error(), "already exists") {
					return err
				}
				foldersMap[bucketName] = []string{}
				// ... (Your logic to create a bucket named bucketName) ...

				// Upload files in the folder to the bucket
				files, err := os.ReadDir(path) // Use os.ReadDir instead of ioutil.ReadDir
				if err != nil {
					return err
				}
				for _, file := range files {
					if !file.IsDir() {
						filePath := filepath.Join(path, file.Name())
						foldersMap[bucketName] = append(foldersMap[bucketName], filePath)
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error creating buckets and uploading files:", err)
			os.Exit(1)
		}
		for bucket, value := range foldersMap {
			for _, filePath := range value {
				fileReader, err := os.Open(filePath)
				if err != nil {
					fmt.Println("error opening file %w", err)
					os.Exit(1)
				}
				defer fileReader.Close()

				// 2. Determine ContentType (replace with your actual logic)
				contentType := "application/octet-stream" // Default content type
				if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
					contentType = "image/jpeg"
				} else if strings.HasSuffix(filePath, ".png") {
					contentType = "image/png"
				} else if strings.HasSuffix(filePath, ".webp") {
					contentType = "image/webp"
				} // Add more content type checks as needed
				fileName := filepath.Base(filePath) // Get only the file name
				// 3. Set FileOptions
				upsert := true
				fileOpts := storage_go.FileOptions{
					ContentType: &contentType,
					Upsert:      &upsert, // Make sure isUpsert is defined appropriately
				}
				// 4. Upload the file
				_, err = supaapi.StorageClient.UploadFile(bucket, fileName, fileReader, fileOpts) // Replace s.supaapi.StorageClient with your actual Supabase client
				if err != nil {
					fmt.Println("error uploading file %s to bucket %s: %w", fileName, bucket, err)
					os.Exit(1)
				}

			}
		}
		db, err := sql.Open("postgres", appConfig.GetConfig().DBSource)
		if err != nil {
			fmt.Println("Error connect to db:", err)
			os.Exit(1)

		}
		defer db.Close() // Ensure the connection is closed after execution

		// Test the connection
		err = db.Ping()
		if err != nil {
			fmt.Println("Error ping db:", err)
			os.Exit(1)

		}
		fmt.Println("Connected to the database successfully!")

		// 2. Load SVG icons and insert into database
		err = filepath.Walk(iconsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".svg" {
				iconName := strcase.ToSnake(strings.ToLower(strings.TrimSuffix(filepath.Base(path), ".svg")))
				iconContent, err := os.ReadFile(path) // Use os.ReadFile instead of ioutil.ReadFile
				if err != nil {
					return err
				}
				iconQuery := fmt.Sprintf("INSERT INTO icons (icon_name, icon_content) VALUES ('%s', '%s')", iconName, string(iconContent))
				// ... (Your logic to execute the INSERT statement) ...
				// Example:
				_, err = db.Exec(iconQuery)
				if err != nil && !strings.Contains(err.Error(), "duplicate") {
					return err
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error loading and inserting icons:", err)
			os.Exit(1)
		}

		fmt.Println("Seed storage command executed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
	storageCmd.Flags().StringP("files-path", "f", "", "Path to the files directory (required)")
	storageCmd.Flags().StringP("icons-path", "i", "", "Path to the icons directory (required)")
	storageCmd.MarkFlagRequired("files-path")
	storageCmd.MarkFlagRequired("icons-path")
}

// ... other code ...
