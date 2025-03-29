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

type EndpointTemplateData struct {
	FeatureTemplateData
	ApiRequestType          string
	ApiReturnType           string
	QueryReturnType         string
	IsEmptyResponse         bool
	IsEmptyRequest          bool
	IsList                  bool
	IsNoSideEffect          bool
	RepoRequestType         string
	RepoReturnTypeActual    string
	RepoReturnTypeWithError string
	EndpointName            string
	EndpointNameLower       string
}

func (c *NewCmd) GetEndpointTemplateData(domainName string, featureName string, endpointName string, flags *pflag.FlagSet) (*EndpointTemplateData, error) {
	featureTemplateData, err := c.GetFeatureTemplateData(domainName, featureName)
	if err != nil {
		return nil, err
	}
	conf, err := c.config.GetProjectConfig()
	if err != nil {
		return nil, err
	}
	rpcService := fmt.Sprintf("%s%s", conf.ApiServiceName, conf.ApiVersion)
	isNoSideEffect, _ := flags.GetBool("get")
	isEmptyResponse, _ := flags.GetBool("empty-esponse")
	isEmptyRequest, _ := flags.GetBool("empty-request")
	isList, _ := flags.GetBool("list")

	endpointFunctionName := strcase.ToCamel(fmt.Sprintf("%s_%s", featureName, endpointName))
	apiRequestType := fmt.Sprintf("%s.%sRequest", rpcService, endpointFunctionName)
	apiReturnType := fmt.Sprintf("%s.%sResponse", rpcService, endpointFunctionName)
	queryReturnType := "one"
	repoRequestType := fmt.Sprintf("db.%sParams", endpointFunctionName)
	repoReturnTypeActual := fmt.Sprintf("*db.%sRow", endpointFunctionName)
	if isEmptyResponse {
		apiReturnType = "emptypb.Empty"
		queryReturnType = "exec"
	}
	if isEmptyRequest {
		apiRequestType = "emptypb.Empty"
	}
	if isList || strings.Contains(endpointName, "list") {
		isNoSideEffect = true
		repoReturnTypeActual = fmt.Sprintf("[]db.%sRow", endpointFunctionName)
		queryReturnType = "many"
	}

	repoReturnTypeWithError := fmt.Sprintf("(%s , error)", repoReturnTypeActual)
	return &EndpointTemplateData{
		FeatureTemplateData:     *featureTemplateData,
		ApiRequestType:          apiRequestType,
		ApiReturnType:           apiReturnType,
		QueryReturnType:         queryReturnType,
		RepoRequestType:         repoRequestType,
		RepoReturnTypeWithError: repoReturnTypeWithError,
		RepoReturnTypeActual:    repoReturnTypeActual,
		IsEmptyResponse:         isEmptyResponse,
		IsEmptyRequest:          isEmptyRequest,
		IsNoSideEffect:          isNoSideEffect || isList,
		IsList:                  isList,
		EndpointNameLower:       endpointName,
		EndpointName:            endpointFunctionName,
	}, nil

}

// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *NewCmd) NewEndpoint(args []string, flags *pflag.FlagSet) {
	endpointName := args[0]
	domainName, err := flags.GetString("domain")
	if err != nil || domainName == "" {
		log.Err(fmt.Errorf("the --domain flag is required")).Msg("")
	}
	featureName, err := flags.GetString("feature")
	if err != nil || featureName == "" {
		log.Err(fmt.Errorf("the --feature flag is required")).Msg("")
	}
	// check if the domain is found
	_, err = os.Stat(fmt.Sprintf("app/%s", domainName))
	if err != nil {
		log.Err(fmt.Errorf("this domain is not found on your app")).Msg("")
	}
	templateData, err := c.GetEndpointTemplateData(domainName, featureName, endpointName, flags)
	if err != nil {
		log.Err(err).Msg("failed to get the templates data")
		os.Exit(1)
	}
	endpointTemplates, err := c.templateUtils.LoadLayerTemplates(fmt.Sprintf("endpoint*"), templateData)
	if err != nil {
		log.Err(err).Msg("error getting the template for the endpoint")
		os.Exit(1)
	}
	endpointFiles := c.GetFeatureFiles(domainName, featureName, templateData.ApiServiceName, templateData.ApiVersion)
	injecteionFiles := map[string]string{
		"adapterinjector": fmt.Sprintf("app/%s/adapter/adapter.go", domainName),
		"usecaseinjector": fmt.Sprintf("app/%s/usecase/usecase.go", domainName),
		"protoinjector":   c.getServiceFilePath(templateData.ApiServiceName, templateData.ApiVersion),
		"repoinjector":    fmt.Sprintf("app/%s/repo/repo.go", domainName),
	}
	for key, template := range endpointTemplates {
		fileName, ok := endpointFiles[key]
		if ok {
			err := c.fileUtils.AppendToFile(fileName, template)
			if err != nil {
				log.Err(err).Str("file", fileName).Msg("error appending to file")
				os.Exit(1)
			}
			log.Info().Str("name", filepath.Base(fileName)).Msg("the endpoint code appended to this file")
			continue

		}
		injectioFile, ok := injecteionFiles[key]
		if !ok {
			continue
		}
		placeHolder := "// INJECT INTERFACE"
		if key == "protoinjector" {
			placeHolder = "// INJECT METHODS"
		}
		err = c.fileUtils.ReplaceFile(injectioFile, placeHolder, fmt.Sprintf("%s\n%s", placeHolder, template.String()))
		if err != nil {
			log.Err(err).Str("file", injectioFile).Msg("error injecting interface for file")
			os.Exit(1)
		}

		log.Info().Str("name", injectioFile).Msg("this file is injected succesfully with the base code")

	}

	log.Info().Str("name", endpointName).Msg("endpoint created")
}
