package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/darwishdev/devkit-cli/gitclient"
	"github.com/darwishdev/devkit-cli/models"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

type ConfigInterface interface {
	GetConfig() models.ConfigModel
	GetConfigMap() models.ConfigMap
	WithFeature(domainName string, featureName string) models.ConfigMap
	Init() error
	WithDomain(domainName string) models.ConfigMap
	WithEndpoint(domainName string, featureName string, endpointName string) models.ConfigMap
	With(values models.ConfigMap) models.ConfigMap
}

type Config struct {
	App       models.ConfigModel
	AppMap    models.ConfigMap
	GitClient gitclient.GitClientInterface
}

func NewConfig() ConfigInterface {
	return &Config{
		App:    models.ConfigModel{},
		AppMap: models.ConfigMap{},
	}
}

// Helper function to find the service name from the proto directory
func findServiceName() (string, error) {
	protoPath := "proto"
	files, err := os.ReadDir(protoPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			return file.Name(), nil // Return the first directory name
		}
	}

	return "", fmt.Errorf("no subdirectories found in proto folder")
}
func (c *Config) Init() error {
	configFilePath := "devkit.yaml"
	// Check if config file already exists
	if _, err := os.Stat(configFilePath); err == nil {
		fmt.Println("Config file already exists:", configFilePath)
		os.Exit(1)
	}
	// Get API service name (can be made interactive)
	apiServiceName, err := findServiceName()
	if err != nil {
		apiServiceName = ""
	}
	// Get API version (can be made interactive)
	apiVersion := "v1"
	// Create Config struct
	exePath, err := os.Executable()
	if err != nil {
		panic(err) // Handle the error appropriately
	}

	exePathDir := strings.Replace(exePath, "/devkit-cli", "", 1)
	config := models.ConfigModel{
		BasePath:               "github.com/yourusername/example",
		ApiServiceName:         apiServiceName,
		ExePath:                exePathDir,
		ApiFilePath:            "api/api.go",
		StorePath:              "db",
		AppName:                "example",
		BaseBuf:                "buf.build/yourusername",
		DBSource:               "yourdbsource",
		Environmet:             "dev",
		GithubToken:            "yourgithubapitoken",
		DBProjectREF:           "supaprojectref",
		SupabaseServiceRoleKey: "supabase service role",
		SupabaseApiKey:         "supabase anon key",
		ProtoGenPath:           "proto_gen",
		ProtoServiceFilePath:   fmt.Sprintf("proto/%s/%s/%s_service.proto", apiServiceName, apiVersion, apiServiceName),
		ProtoPath:              fmt.Sprintf("proto/%s/%s", apiServiceName, apiVersion),
		QueryPath:              "supabase/queries",
		ApiVersion:             apiVersion,
	}
	// Marshal config data to YAML
	configData, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		os.Exit(1)
	}

	// Write config file
	err = os.WriteFile(configFilePath, configData, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		os.Exit(1)
	}
	return nil
}

func (c *Config) ReadConfig() {
	configFilePath := "devkit.yaml"
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
	}
	configMap := make(models.ConfigMap)
	var config models.ConfigModel
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(fmt.Sprintf("canot load config%s", err))
	}
	configMap["BasePath"] = config.BasePath
	configMap["ApiServiceName"] = config.ApiServiceName
	configMap["QueryPath"] = config.QueryPath
	configMap["StorePath"] = config.StorePath
	configMap["ApiFilePath"] = config.ApiFilePath
	configMap["ProtoGenPath"] = config.ProtoGenPath
	configMap["ProtoPath"] = config.ProtoPath
	configMap["ApiVersion"] = config.ApiVersion
	c.App = config
	c.AppMap = configMap
}
func (c *Config) GetConfig() models.ConfigModel {
	c.ReadConfig()
	return c.App
}
func (c *Config) GetConfigMap() models.ConfigMap {
	c.ReadConfig()
	return c.AppMap
}

func (c *Config) WithEndpoint(domainName string, featureName string, endpointName string) models.ConfigMap {
	queryReturnType := "one"
	endPointFunctionName := strcase.ToCamel(fmt.Sprintf("%s_%s", featureName, endpointName))
	repoReturnType := fmt.Sprintf("*db.%sRow", endPointFunctionName)
	if strings.Contains(endpointName, "list") {
		queryReturnType = "many"
		repoReturnType = fmt.Sprintf("[]db.%sRow", endPointFunctionName)
		endPointFunctionName = strcase.ToCamel(fmt.Sprintf("%s_%s", featureName, endpointName))
	}

	return c.With(models.ConfigMap{
		"RepoReturnType":   repoReturnType,
		"RepoRequestType":  fmt.Sprintf("%sParams", endPointFunctionName),
		"QueryReturnType":  queryReturnType,
		"EndpointName":     endPointFunctionName,
		"ApiReturnType":    fmt.Sprintf("%sResponse", endPointFunctionName),
		"ApiRequestType":   fmt.Sprintf("%sRequest", endPointFunctionName),
		"FeatureNameLower": featureName,
		"FeatureName":      strcase.ToCamel(featureName),
		"DomainNameLower":  domainName,
		"DomainName":       strcase.ToCamel(domainName),
	})
}

func (c *Config) WithFeature(domainName string, featureName string) models.ConfigMap {
	return c.With(models.ConfigMap{"FeatureNameLower": featureName, "FeatureName": strcase.ToCamel(featureName), "DomainNameLower": domainName, "DomainName": strcase.ToCamel(domainName)})
}

func (c *Config) WithDomain(domainName string) models.ConfigMap {
	return c.With(models.ConfigMap{"DomainNameLower": domainName, "DomainName": strcase.ToCamel(domainName)})
}

func (c *Config) With(values models.ConfigMap) models.ConfigMap {
	configCopy := make(map[string]interface{})
	for key, value := range c.AppMap {
		configCopy[key] = value
	}
	for k, v := range values {
		configCopy[k] = v
	}
	return configCopy
}
