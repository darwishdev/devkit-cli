package new

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

type DomainTemplateData struct {
	config.ProjectConfig
	DomainName      string
	BaseDomainPath  string
	BasePath        string
	DomainNameLower string
	UsecaseName     string
}

func (c *NewCmd) GetDomainTemplateData(domainName string) (*DomainTemplateData, error) {
	projectConfig, err := c.config.GetProjectConfig()
	if err != nil {
		return nil, err
	}

	return &DomainTemplateData{
		ProjectConfig:   *projectConfig,
		DomainNameLower: domainName,
		DomainName:      strcase.ToCamel(domainName),
		UsecaseName:     fmt.Sprintf("%sUsecase", strcase.ToLowerCamel(domainName)),
		BaseDomainPath:  fmt.Sprintf("github.com/%s/%s/app/%s", projectConfig.GitUser, projectConfig.AppName, domainName),
		BasePath:        fmt.Sprintf("github.com/%s/%s", projectConfig.GitUser, projectConfig.AppName),
	}, nil

}

// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *NewCmd) NewDomain(args []string, flags *pflag.FlagSet) {
	domainName := args[0]
	templateData, err := c.GetDomainTemplateData(domainName)
	if err != nil {
		log.Err(err).Msg("failed to get the project config")
		os.Exit(1)
	}
	domainTemplates, err := c.templateUtils.LoadLayerTemplates(fmt.Sprintf("domain*"), templateData)
	if err != nil {
		log.Err(err).Msg("error getting the template for the domain")
		os.Exit(1)
	}

	adapterFolder := fmt.Sprintf("%s/%s/adapter", c.domainsFolderPath, domainName)
	repoFolder := fmt.Sprintf("%s/%s/repo", c.domainsFolderPath, domainName)
	usecaseFolder := fmt.Sprintf("%s/%s/usecase", c.domainsFolderPath, domainName)
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
		"schema":  fmt.Sprintf("%s/usecase.go", usecaseFolder),
	}

	for key, fileName := range domainFiles {
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating:", key, err)
			os.Exit(1)
		}
		log.Info().Str("name", filepath.Base(fileName)).Msg("new file created")
		template, ok := domainTemplates[key]
		if ok {
			err = c.fileUtils.AppendToFile(file.Name(), *bytes.NewBuffer(template.Bytes()))
			if err != nil {
				fmt.Println("Error adding base content for:", key, err)
				os.Exit(1)
			}
		}
		log.Info().Str("name", filepath.Base(fileName)).Msg("new file filled with base code")

	}
	usecaseImport, _ := domainTemplates["import"]
	usecaseField, _ := domainTemplates["field"]
	usecaseInstantiation, _ := domainTemplates["instantiation"]
	usecaseInjection, _ := domainTemplates["injection"]
	err = c.fileUtils.ReplaceMultiple("api/api.go", map[string]string{
		"// USECASE_IMPORTS":        fmt.Sprintf("// USECASE_IMPORTS\n%s", usecaseImport.String()),
		"// USECASE_FIELDS":         fmt.Sprintf("// USECASE_FIELDS\n%s", usecaseField.String()),
		"// USECASE_INSTANTIATIONS": fmt.Sprintf("// USECASE_INSTANTIATIONS\n%s", usecaseInstantiation.String()),
		"// USECASE_INJECTIONS":     fmt.Sprintf("// USECASE_INJECTIONS\n%s", usecaseInjection.String()),
	})
	if err != nil {
		log.Err(err).Msg("can't inject the api file")
	}
	log.Info().Msg("the api/api.go injected successfully with the new usecase")
	err = c.ExecCmd("", "supabase", "migration", "new", fmt.Sprintf("%s_schema", domainName))
	if err != nil {
		log.Err(err).Msg("supabase migration failed")
		os.Exit(1)
	}
	log.Info().Str("str", "new domain from domain").Msg("domain")
}
