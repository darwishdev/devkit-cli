package new

import (
	"fmt"
	"os"

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
	featureTemplates, err := c.templateUtils.LoadLayerTemplates(fmt.Sprintf("feature*"), templateData)
	if err != nil {
		log.Err(err).Msg("error getting the template for the domain")
		os.Exit(1)
	}

	basePath := fmt.Sprintf("app/%s", domainName)
	featureFiles := map[string]string{
		"adapter": fmt.Sprintf("%s/adapter/%s_adapter.go", basePath, featureName),
		"repo":    fmt.Sprintf("%s/repo/%s_repo.go", basePath, featureName),
		"usecase": fmt.Sprintf("%s/usecase/%s_usecase.go", basePath, featureName),
		"proto":   fmt.Sprintf("proto/%s/%s/%s_%s.proto", templateData.ApiServiceName, templateData.ApiVersion, domainName, featureName),
		"query":   fmt.Sprintf("supabase/queries/%s_%s.sql", domainName, featureName),
		"api":     fmt.Sprintf("api/%s_%s_rpc.go", domainName, featureName),
	}
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
	serviceFilePath := fmt.Sprintf("proto/%s/%s/%s_service.proto", templateData.ApiServiceName, templateData.ApiVersion, templateData.ApiServiceName)

	err = c.fileUtils.ReplaceFile(serviceFilePath, "// INJECT IMPORTS", fmt.Sprintf("// INJECT IMPORTS\n%s;", protoImport))
	if err != nil {
		fmt.Println("Error replacing api", err)
		os.Exit(1)
	}
	log.Info().Str("str", "new domain from domain").Msg("domain")
}
