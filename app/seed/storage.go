package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *SeedCmd) StorageSeed(flags *pflag.FlagSet) {
	log.Info().Str("str", "new domain from domain").Msg("domain")
	filesPath, _ := flags.GetString("files-path")
	iconsPath, _ := flags.GetString("icons-path")
	conf, err := c.config.GetProjectConfig()
	if err != nil {
		log.Err(err).Msg("can't read the project config")
		os.Exit(1)
	}

	// Basic validation
	if filesPath == "" && iconsPath == "" {
		log.Err(fmt.Errorf("Error: --files-path or--icons-path  one of them must be passed to this command")).Msg("")
		os.Exit(1)
	}
	if filesPath != "" {
		err = c.supaClient.StorageSeed(conf, filesPath)
		if err != nil {
			log.Err(err).Msg("Error creating buckets and uploading files:")
			os.Exit(1)
		}
	}
	if iconsPath != "" {
		db, err := c.dbUtils.Open(conf.DBSource)
		if err != nil {
			log.Err(err).Msg("Error connecting to the database")
			os.Exit(1)
		}
		defer db.Close()
		err = filepath.Walk(iconsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := strings.ToLower(filepath.Ext(path))
			supportedExt := ".svg"
			baseFileName := filepath.Base(path)
			if !info.IsDir() && ext == supportedExt {
				nameWithotExt := strings.ToLower(strings.TrimSuffix(baseFileName, supportedExt))
				iconName := strcase.ToSnake(nameWithotExt)
				iconContent, err := os.ReadFile(path) // Use os.ReadFile instead of ioutil.ReadFile
				if err != nil {
					return err
				}
				if iconName != "" {
					iconQuery := fmt.Sprintf("INSERT INTO icons (icon_name, icon_content) VALUES ('%s', '%s')", iconName, string(iconContent))
					_, err = db.Exec(iconQuery)
					if err != nil && !strings.Contains(err.Error(), "duplicate") {
						return err
					}

				}
			}
			return nil
		})
		if err != nil {
			log.Err(err).Msg("error executing the icons insert")
			os.Exit(1)
		}

	}
	log.Info().Msg("your storage is seeded succesfully!")
}
