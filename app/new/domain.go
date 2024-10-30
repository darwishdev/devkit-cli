package new

import (
	"fmt"
	"os"

	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

type DomainTemplateData struct {
	config.ProjectConfig
	DomainName      string
	DomainNameLower string
}

func (c *NewCmd) NewDomain(args []string, flags *pflag.FlagSet) {
	domainName := args[0]
	projectConfig, err := c.config.GetProjectConfig()
	if err != nil {
		log.Err(err).Msg("failed to get the project config")
		os.Exit(1)
	}
	templateData := DomainTemplateData{
		ProjectConfig:   *projectConfig,
		DomainNameLower: domainName,
		DomainName:      strcase.ToCamel(domainName),
	}
	log.Debug().Interface("template data", templateData).Msg("temp data")
	os.Exit(1)
	adapterFolder := fmt.Sprintf("%s/adapter", c.basePath)
	repoFolder := fmt.Sprintf("%s/repo", c.basePath)
	usecaseFolder := fmt.Sprintf("%s/usecase", c.basePath)
	domainDirs := map[string]string{
		"adapter": adapterFolder,
		"repo":    repoFolder,
		"usecase": usecaseFolder,
	}

	for _, dir := range domainDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Err(err).Str("folder name : ", dir).Msg("error creating folder")
				os.Exit(1)
			}
		}
	}

	domainFiles := map[string]string{
		"adapter": fmt.Sprintf("%s/adapter.go", adapterFolder),
		"repo":    fmt.Sprintf("%s/repo.go", repoFolder),
		"usecase": fmt.Sprintf("%s/usecase.go", usecaseFolder),
	}

	for key, fileName := range domainFiles {
		_, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating:", key, err)
			os.Exit(1)
		}

	}
	log.Debug().Str("str", "new domain from domain").Msg("domain")
}
