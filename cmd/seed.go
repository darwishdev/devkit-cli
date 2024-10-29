package cmd

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"

	supaapigo "github.com/darwishdev/supaapi-go"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/supabase-community/auth-go/types"
	"github.com/xuri/excelize/v2"
)

// ... other code ...

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed [schema_name] --file-path <file_path> [--outfile <outfile>] [--execute <true|false>]",
	Short: "Seed data to a database schema",
	Long:  `Seed data to a database schema from a file.`,
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument (schema_name) is provided
	Run: func(cmd *cobra.Command, args []string) {
		schemaName := args[0]
		filePath, _ := cmd.Flags().GetString("file-path")
		outFile, _ := cmd.Flags().GetString("outfile")
		execute, _ := cmd.Flags().GetBool("execute")
		conf := appConfig.GetConfig()
		// Basic validation
		if filePath == "" {
			fmt.Println("Error: --file-path is required")
			os.Exit(1)
		}
		excelFile, err := os.ReadFile("accounts.xlsx") // Replace with your Excel file path
		if err != nil {
			log.Fatal().Err(err).Msg("cano open the file")
		}
		buffer := bytes.NewBuffer(excelFile)

		f, err := excelize.OpenReader(bytes.NewBuffer(excelFile))
		if err != nil {
			log.Fatal().Err(err).Msg("cannot open the excel")
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal().Err(err).Msg("cannot close the excel")
			}
		}()
		queries := ""
		for _, sheetName := range f.GetSheetList() {
			queryString, err := appSeeder.SeedFromExcel(*buffer, schemaName, sheetName, sheetName)
			if err != nil {
				log.Fatal().Err(err).Msg("can't get the sql query")
			}
			if sheetName == "users" {
				supaapi := supaapigo.NewSupaapi(supaapigo.SupaapiConfig{
					ProjectRef:     conf.DBProjectREF,
					Env:            supaapigo.DEV,
					Port:           54321,
					ServiceRoleKey: conf.SupabaseServiceRoleKey,
					ApiKey:         conf.SupabaseApiKey,
				})

				rows, err := f.GetRows(sheetName)
				if err != nil {
					fmt.Println("Error getting user rows:", err)
					os.Exit(1)

				}
				columns := rows[0]
				for _, row := range rows[1:] { // Start from the second row (index 1)
					supabasRequest := types.AdminUpdateUserRequest{}
					for colIndex, colCell := range row {
						if columns[colIndex] == "user_email" {
							supabasRequest.Email = colCell
						}
						if columns[colIndex] == "user_password#" {
							supabasRequest.Password = colCell
						}
					}

					_, err := supaapi.UserCreateUpdate(supabasRequest)
					if err != nil {
						fmt.Println("Error signup on supabase:", err)
						os.Exit(1)

					}
				}
			}
			queries += fmt.Sprintf("\n%s\n", queryString)
		}

		if execute {
			// Connect to the database
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

			// Run the provided query
			_, err = db.Exec(queries)
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				os.Exit(1)

			}
		}
		if outFile != "" {
			err = os.WriteFile(outFile, []byte(queries), 0644)
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				os.Exit(1)
			}
			fmt.Println("SQL written to:", outFile)
		}
		// ... (Your logic to generate SQL from the file) ...

		fmt.Println("Seed command executed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().StringP("file-path", "f", "", "Path to the data file (required)")
	seedCmd.Flags().StringP("outfile", "o", "", "Path to the output SQL file (optional)")
	seedCmd.Flags().BoolP("execute", "e", false, "Execute the generated SQL (optional)")
	seedCmd.MarkFlagRequired("file-path") // Make file-path flag required
}

// ... other code ...
