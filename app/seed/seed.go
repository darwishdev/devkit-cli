package seed

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *SeedCmd) NewSeed(args []string, flags *pflag.FlagSet) {
	schemaName := args[0]
	filePath, err := flags.GetString("file-path")
	if err != nil {
		log.Err(err).Msg("file path not passed")
		os.Exit(1)
	}
	conf, err := c.config.GetProjectConfig()
	if err != nil {
		log.Err(err).Msg("can't read the project config")
		os.Exit(1)
	}

	outFile, _ := flags.GetString("out-file")
	isExecute, _ := flags.GetBool("execute")
	isSkipSupabase, _ := flags.GetBool("skip-supabase")
	file, buffer, err := c.fileUtils.ReadExcelFile(filePath)
	if err != nil {
		log.Err(err).Msg("file path not passed")
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Err(err).Msg("can't close the excel file")
			os.Exit(1)
		}
	}()
	fullSqlQuery := ""
	for _, sheetName := range file.GetSheetList() {
		queryString, err := c.sqlSeeder.SeedFromExcel(*buffer, schemaName, sheetName, sheetName)
		if err != nil {
			log.Err(err).Str("Sheet", sheetName).Msg("can't get the sql query")
			os.Exit(1)
		}
		fullSqlQuery += fmt.Sprintf("\n%s\n", queryString)
		if sheetName == "users" && !isSkipSupabase {
			rows, err := file.GetRows(sheetName)
			if err != nil {
				log.Err(err).Msg("can't get users sheet rows")
				os.Exit(1)
			}
			err = c.supaClient.UsersCreateUpdate(conf, rows)
			if err != nil {
				log.Err(err).Msg("can't insert users on supabase")
				os.Exit(1)
			}
		}
	}

	if isExecute {
		db, err := c.dbUtils.Open(conf.DBSource)
		defer db.Close()

		if err != nil {
			log.Err(err).Str("source", conf.DBSource).Msg("can't connect to the database")
			os.Exit(1)
		}
		_, err = db.Exec(fullSqlQuery)
		if err != nil {
			log.Err(err).Msg("error executing the insert statement")
			os.Exit(1)
		}
	}
	if outFile != "" {
		err = os.WriteFile(outFile, []byte(fullSqlQuery), 0644)
		if err != nil {
			log.Err(err).Msg("error executing the insert statement")
		}
	}
	log.Info().Str("str", "new domain from domain").Msg("domain")
}
