package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
	GitUser                string `mapstructure:"GIT_USER"`
	DockerhubUser          string `mapstructure:"DOCKER_HUB_USER"`
	State                  string `mapstructure:"STATE"`
	TokenSymmetricKey      string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ClientBaseUrl          string `mapstructure:"CLIENT_BASE_URL"`
	RedisAddress           string `mapstructure:"REDIS_ADDRESS"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase          int    `mapstructure:"REDIS_DATABASE"`
	GRPCServerAddress      string `mapstructure:"GRPC_SERVER_ADDRESS"`
	BufUser                string `mapstructure:"BUF_USER"`
	GoogleClientId         string `mapstructure:"GOOGLE_CLIENT_ID"`
	ResendApiKey           string `mapstructure:"RESEND_API_KEY"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GithubToken            string `mapstructure:"GITHUB_TOKEN"`
	ApiServiceName         string `mapstructure:"API_SERVICE_NAME"`
	ApiVersion             string `mapstructure:"API_VERSION"`
	Environmet             string `mapstructure:"ENVIRONMENT"`
	DBProjectREF           string `mapstructure:"DB_PROJECT_REF"`
	SupabaseServiceRoleKey string `mapstructure:"SUPABASE_SERVICE_ROLE_KEY"`
	SupabaseApiKey         string `mapstructure:"SUPABASE_API_KEY"`
	SupabaseAPIPort        uint32 `mapstructure:"SUPABASE_API_PORT"`
	DBPort                 uint32 `mapstructure:"DB_PORT"`
	DBSource               string `mapstructure:"DB_SOURCE"`
	AppName                string `mapstructure:"APP_NAME"`
	ApiFilePath            string `mapstructure:"API_FILE_PATH"`
	ServiceFilePath        string `mapstructure:"SERVICE_FILE_PATH"`
}
type ConfigInterface interface {
	GetCliConfig() *CliConfig
	GetProjectConfig() (*ProjectConfig, error)
	InitProjectConfig() error
}
type Config struct {
	CliConfing      *CliConfig
	ProjectFileName string
	ProjectFilePath string
}

type StateConfig struct {
	State string `mapstructure:"STATE"`
}

func LoadState(path string) (config StateConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("state")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func readConf(path string, name string, target interface{}) error {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType("env")

	log.Debug().Interface("path is", path).Str("name", name).Msg("state")
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Can't read the config from  %s : %w", name, err)
	}
	err = v.Unmarshal(&target)
	if err != nil {
		return fmt.Errorf("Can't Unmarshal The Config from %s: %w", name, err)
	}
	return nil
}

func NewConfig(cliConfigPath string, cliConfigName string, projectFilePath string, projectFileName string) (ConfigInterface, error) {
	cliConfig := CliConfig{}
	err := readConf(cliConfigPath, cliConfigName, &cliConfig)
	if err != nil {
		return nil, err
	}
	return &Config{
		CliConfing:      &cliConfig,
		ProjectFilePath: projectFilePath,
		ProjectFileName: projectFileName,
	}, err
}

// InitProjectConfig creates a new project configuration file with default values
// based on the global CLI configuration.
func (c *Config) InitProjectConfig() error {
	// fullPath := fmt.Sprintf("%s/%s.env", c.ProjectFilePath, c.ProjectFileName)
	// v := viper.New()
	// currenDireName, err := os.Getwd()
	//
	//	if err != nil {
	//		return fmt.Errorf("can't get the current dir name : %w", err)
	//	}
	//
	// v.AddConfigPath(c.ProjectFilePath)
	// v.SetConfigName(c.ProjectFileName)
	// v.SetConfigType("env")
	// _, err = os.Stat(fullPath)
	//
	//	if err == nil {
	//		return fmt.Errorf("file already exists")
	//	}
	//
	// // Set the values using Viper
	// v.Set("GIT_USER", c.CliConfing.GitUser)
	// v.Set("DOCKERHUB_USER", c.CliConfing.DockerhubUser)
	// v.Set("BUF_USER", c.CliConfing.BufUser)
	// v.Set("GOOGLE_CLIENT_ID", c.CliConfing.GoogleClientId)
	// v.Set("RESEND_API_KEY", c.CliConfing.ResendApiKey)
	// v.Set("GOOGLE_CLIENT_SECRET", c.CliConfing.GoogleClientSecret)
	// v.Set("GITHUB_TOKEN", c.CliConfing.GithubToken)
	// v.Set("API_SERVICE_NAME", c.CliConfing.ApiServiceName)
	// v.Set("API_VERSION", c.CliConfing.ApiVersion)
	// v.Set("ENVIRONMENT", "dev")
	// v.Set("DB_PORT", 54321)
	// v.Set("DB_PROJECT_REF", "")
	// v.Set("SUPABASE_SERVICE_ROLE_KEY", "")
	// v.Set("SUPABASE_API_KEY", "")
	// v.Set("DB_SOURCE", "")
	// v.Set("APP_NAME", filepath.Base(currenDireName))
	// v.SafeWriteConfig()
	return nil
}
func (c *Config) GetCliConfig() *CliConfig {
	return c.CliConfing
}
func (c *Config) GetProjectConfig() (*ProjectConfig, error) {
	projectConfig := ProjectConfig{}
	state, err := LoadState("config")
	if err != nil {
		return nil, err
	}
	log.Debug().Interface("project config is from here", state.State).Msg("proejcts state")
	err = readConf("config", state.State, &projectConfig)
	if err != nil {
		return nil, err
	}
	projectConfig.State = state.State
	return &projectConfig, nil
}
