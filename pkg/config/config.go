package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
)

type CliConfig struct {
	GitUser            string `mapstructure:"GIT_USER"`
	DockerhubUser      string `mapstructure:"DOCKER_HUB_USER"`
	BufUser            string `mapstructure:"BUF_USER"`
	GoogleClientId     string `mapstructure:"GOOGLE_CLIENT_ID"`
	ResendApiKey       string `mapstructure:"RESEND_API_KEY"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GithubToken        string `mapstructure:"GITHUB_TOKEN"`
	ApiServiceName     string `mapstructure:"API_SERVICE_NAME"`
	ApiVersion         string `mapstructure:"API_VERSION"`
}

type ProjectConfig struct {
	GitUser                string `json:"GIT_USER"`
	DockerhubUser          string `json:"DOCKER_HUB_USER"`
	BufUser                string `json:"BUF_USER"`
	GoogleClientId         string `json:"GOOGLE_CLIENT_ID"`
	ResendApiKey           string `json:"RESEND_API_KEY"`
	GoogleClientSecret     string `json:"GOOGLE_CLIENT_SECRET"`
	GithubToken            string `json:"GITHUB_TOKEN"`
	ApiServiceName         string `json:"API_SERVICE_NAME"`
	ApiVersion             string `json:"API_VERSION"`
	Environmet             string `json:"ENVIRONMENT"`
	DBProjectREF           string `json:"DB_PROJECT_REF"`
	SupabaseServiceRoleKey string `json:"SUPABASE_SERVICE_ROLE_KEY"`
	SupabaseApiKey         string `json:"SUPABASE_API_KEY"`
	DBSource               string `json:"DB_SOURCE"`
	AppName                string `json:"APP_NAME"`
	ApiFilePath            string `json:"API_FILE_PATH"`
	ServiceFilePath        string `json:"SERVICE_FILE_PATH"`
}
type ConfigInterface interface {
	GetCliConfig() *CliConfig
	GetProjectConfig() (*ProjectConfig, error)
	InitProjectConfig() error
}
type Config struct {
	CliConfing      *CliConfig
	ProjectFilePath string
}

func readConf(path string, name string, target interface{}) error {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("cannot read the config from  %s : %w", path, err)
	}
	err = v.Unmarshal(&target)
	if err != nil {
		return fmt.Errorf("cannot parse the config from %s: %w", path, err)
	}
	return nil
}

func NewConfig(cliConfigPath string, cliConfigName string, projectFilePath string) (ConfigInterface, error) {
	cliConfig := CliConfig{}
	err := readConf(cliConfigPath, cliConfigName, &cliConfig)
	if err != nil {
		return nil, err
	}
	return &Config{
		CliConfing:      &cliConfig,
		ProjectFilePath: projectFilePath,
	}, err
}

func (c *Config) InitProjectConfig() error {
	_, err := os.Stat(c.ProjectFilePath)
	if err == nil {
		return fmt.Errorf("file already exists")
	}
	config := ProjectConfig{
		GitUser:                c.CliConfing.GitUser,
		DockerhubUser:          c.CliConfing.DockerhubUser,
		BufUser:                c.CliConfing.BufUser,
		GoogleClientId:         c.CliConfing.GoogleClientId,
		ResendApiKey:           c.CliConfing.ResendApiKey,
		GoogleClientSecret:     c.CliConfing.GoogleClientSecret,
		GithubToken:            c.CliConfing.GithubToken,
		ApiServiceName:         c.CliConfing.ApiServiceName,
		ApiVersion:             c.CliConfing.ApiVersion,
		Environmet:             "",
		DBProjectREF:           "",
		SupabaseServiceRoleKey: "",
		SupabaseApiKey:         "",
		DBSource:               "",
		AppName:                "",
		ApiFilePath:            "",
		ServiceFilePath:        "",
	}
	// Marshal config data to YAML
	configData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	// Write config file
	err = os.WriteFile(c.ProjectFilePath, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}
func (c *Config) GetCliConfig() *CliConfig {
	return c.CliConfing
}
func (c *Config) GetProjectConfig() (*ProjectConfig, error) {
	file, err := os.ReadFile(c.ProjectFilePath)
	if err != nil {

	}
	config := ProjectConfig{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
