package new

import (
	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/darwishdev/devkit-cli/pkg/fileutils"
	"github.com/darwishdev/devkit-cli/pkg/gitclient"
	"github.com/darwishdev/devkit-cli/pkg/templates"
	"github.com/spf13/pflag"
)

type NewCmdInterface interface {
	NewApi(args []string, flags *pflag.FlagSet)
	NewFeature(args []string, flags *pflag.FlagSet)
	NewDomain(args []string, flags *pflag.FlagSet)
}
type NewCmd struct {
	config        config.ConfigInterface
	fileUtils     fileutils.FileUtilsInterface
	templateUtils templates.TemplatesInterface
	gitClient     gitclient.GitClientInterface
	basePath      string
}

func NewNewCmd(config config.ConfigInterface, fileUtils fileutils.FileUtilsInterface, templateUtils templates.TemplatesInterface, gitClient gitclient.GitClientInterface) NewCmdInterface {
	return &NewCmd{
		config:        config,
		fileUtils:     fileUtils,
		templateUtils: templateUtils,
		gitClient:     gitClient,
		basePath:      "app",
	}
}
