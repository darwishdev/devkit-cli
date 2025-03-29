package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

type FeatureTemplateData struct {
	DomainTemplateData
	FeatureName      string
	FeatureNameLower string
	UsecaseName      string
}

func (c *NewCmd) GetFeatureFiles(domainName string, featureName string, serviceName string, serviceVersion string) map[string]string {
	basePath := fmt.Sprintf("%s/%s", c.domainsFolderPath, domainName)
	return map[string]string{
		"adapter": fmt.Sprintf("%s/adapter/%s_adapter.go", basePath, featureName),
		"repo":    fmt.Sprintf("%s/repo/%s_repo.go", basePath, featureName),
		"usecase": fmt.Sprintf("%s/usecase/%s_usecase.go", basePath, featureName),
		"proto":   fmt.Sprintf("proto/%s/%s/%s_%s.proto", serviceName, serviceVersion, domainName, featureName),
		"query":   fmt.Sprintf("%s/%s_%s.sql", c.queriesFolderPath, domainName, featureName),
		"api":     fmt.Sprintf("api/%s_%s_rpc.go", domainName, featureName),
	}
}
func (c *NewCmd) GetFeatureTemplateData(domainName string, featureName string) (*FeatureTemplateData, error) {
	domainTemplateData, err := c.GetDomainTemplateData(domainName)
	if err != nil {
		return nil, err
	}

	return &FeatureTemplateData{
		DomainTemplateData: *domainTemplateData,
		FeatureNameLower:   featureName,
		FeatureName:        strcase.ToCamel(featureName),
	}, nil

}

func (c *NewCmd) InheritFiles(pattern string, inheritFrom string, featureName string) {
	// patter is app/{domain}/*/*_{feature}.go
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Err(err)
		os.Exit(1)
	}

	log.Info().Interface("from", matches).Str("to", pattern).Msg("file copied and renamed")
	for _, filePath := range matches {
		// Extract the directory and base file name
		dir := filepath.Dir(filePath)
		base := filepath.Base(filePath)

		// Replace the inheritFrom part with featureName in the base file name
		newFileName := strings.Replace(base, inheritFrom, featureName, 1)

		// Construct the new file path
		newFilePath := filepath.Join(dir, newFileName)

		// Copy the file content
		input, err := os.ReadFile(filePath)
		if err != nil {
			log.Err(err).Str("file", filePath).Msg("error reading file")
			os.Exit(1)
		}
		newInput := strings.ReplaceAll(string(input), inheritFrom, featureName)
		newInput = strings.ReplaceAll(newInput, strcase.ToCamel(inheritFrom), strcase.ToCamel(featureName))

		err = os.WriteFile(newFilePath, []byte(newInput), 0644)
		if err != nil {
			log.Err(err).Str("file", newFilePath).Msg("error writing file")
			os.Exit(1)
		}

		log.Info().Str("from", filePath).Str("to", newFilePath).Msg("file copied and renamed")
	}

}

// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *NewCmd) NewFeature(args []string, flags *pflag.FlagSet) {
	featureName := args[0]
	domainName, err := flags.GetString("domain")
	if err != nil || domainName == "" {
		log.Err(fmt.Errorf("the --domain flag is required")).Msg("")
	}

	// check if the domain is found
	_, err = os.Stat(fmt.Sprintf("app/%s", domainName))
	if err != nil {
		log.Err(fmt.Errorf("this domain is not found on your app")).Msg("")
	}
	templateData, err := c.GetFeatureTemplateData(domainName, featureName)
	if err != nil {
		log.Err(err).Msg("failed to get the project config")
		os.Exit(1)
	}
	inheritFrom, err := flags.GetString("inherit")
	if err == nil && len(inheritFrom) > 0 {
		// patter is app/{domain}/*/*_{feature}.go
		appPattern := fmt.Sprintf("%s/%s/*/%s_*.go", c.domainsFolderPath, domainName, inheritFrom)
		apiPattern := fmt.Sprintf("api/%s_%s_rpc.go", domainName, inheritFrom)
		dbPattern := fmt.Sprintf("supabase/queries/%s_%s.sql", domainName, inheritFrom)
		protoPattern := fmt.Sprintf("proto/%s/%s/%s_%s.proto", templateData.ApiServiceName, templateData.ApiVersion, domainName, inheritFrom)
		c.InheritFiles(appPattern, inheritFrom, featureName)
		c.InheritFiles(dbPattern, inheritFrom, featureName)
		c.InheritFiles(protoPattern, inheritFrom, featureName)
		c.InheritFiles(apiPattern, inheritFrom, featureName)
		os.Exit(0)
	}

	featureTemplates, err := c.templateUtils.LoadLayerTemplates(fmt.Sprintf("feature*"), templateData)
	if err != nil {
		log.Err(err).Msg("error getting the template for the domain")
		os.Exit(1)
	}

	featureFiles := c.GetFeatureFiles(domainName, featureName, templateData.ApiServiceName, templateData.ApiVersion)
	log.Debug().Interface("test", featureTemplates).Msg("test")
	for key, fileName := range featureFiles {
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating:", key, err)
			os.Exit(1)
		}
		template, ok := featureTemplates[key]
		if ok {

			_, err = file.Write(template.Bytes())
			if err != nil {
				fmt.Println("Error adding base content for:", key, err)
				os.Exit(1)
			}

		}

	}
	protoImport := fmt.Sprintf("import \"%s/%s/%s_%s.proto\"", templateData.ApiServiceName, templateData.ApiVersion, domainName, featureName)
	serviceFilePath := c.getServiceFilePath(templateData.ApiServiceName, templateData.ApiVersion)
	err = c.fileUtils.ReplaceFile(serviceFilePath, "// INJECT IMPORTS", fmt.Sprintf("// INJECT IMPORTS\n%s;", protoImport))
	if err != nil {
		fmt.Println("Error replacing api", err)
		os.Exit(1)
	}
	log.Info().Str("name", featureName).Msg("feature created successfully")
}
