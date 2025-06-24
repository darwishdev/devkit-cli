package seed

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
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
		log.Err(fmt.Errorf("error: --files-path or--icons-path  one of them must be passed to this command")).Msg("")
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
		m := minify.New()

		m.AddFunc("image/svg+xml", svg.Minify)
		err = filepath.Walk(iconsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := strings.ToLower(filepath.Ext(path))
			supportedExt := ".svg"
			baseFileName := filepath.Base(path)
			if !info.IsDir() && ext == supportedExt {
				nameWithotExt := strings.ToLower(strings.TrimSuffix(baseFileName, supportedExt))
				iconName := strings.ToLower(nameWithotExt)
				iconContent, err := os.ReadFile(path) // Use os.ReadFile instead of ioutil.ReadFile
				if err != nil {
					return err
				}

				// Minify the SVG content
				var minified bytes.Buffer
				err = m.Minify("image/svg+xml", &minified, bytes.NewReader(iconContent))
				if err != nil {
					return err
				}

				if iconName != "" {
					iconQuery := "INSERT INTO icon (icon_name,icon_content) VALUES ($1, $2)"
					_, err = db.Exec(iconQuery, iconName,  minified.String()  )
					if err != nil && !strings.Contains(err.Error(), "duplicate") {
						return err
					}

				}
			} else {
				category := filepath.Base(path) // the folder name

        // list all SVGs directly in this category folder
        entries, err := os.ReadDir(path)
        if err != nil {
            return err
        }
        for _, entry := range entries {
            if entry.IsDir() || filepath.Ext(entry.Name()) != supportedExt {
                continue
            }

            filePath := filepath.Join(path, entry.Name())
            raw, err := os.ReadFile(filePath)
            if err != nil {
                return err
            }
            var minified bytes.Buffer
            if err = m.Minify("image/svg+xml", &minified, bytes.NewReader(raw)); err != nil {
                return err
            }

            name := strings.TrimSuffix(entry.Name(), supportedExt)
            iconName := strings.ToLower(name)

            // now insert with category
            _, err = db.Exec(
                "INSERT INTO icon (icon_name, icon_content, icon_category) VALUES ($1, $2, $3)",
                iconName, minified.String(), category,
            )
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
