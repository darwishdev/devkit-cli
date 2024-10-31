package new

import (
	"fmt"

	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/darwishdev/devkit-cli/pkg/db"
	"github.com/darwishdev/devkit-cli/pkg/fileutils"
	"github.com/darwishdev/devkit-cli/pkg/gitclient"
	"github.com/darwishdev/devkit-cli/pkg/templates"
	"github.com/spf13/pflag"
)

type NewCmdInterface interface {
	NewApi(args []string, flags *pflag.FlagSet)
	getServiceFilePath(serviceName, version string) string
	NewFeature(args []string, flags *pflag.FlagSet)
	NewEndpoint(args []string, flags *pflag.FlagSet)
	NewDomain(args []string, flags *pflag.FlagSet)
}
type NewCmd struct {
	config        config.ConfigInterface
	fileUtils     fileutils.FileUtilsInterface
	templateUtils templates.TemplatesInterface
	gitClient     gitclient.GitClientInterface
	dbUtils       db.DbInterface
	basePath      string
}

func NewNewCmd(config config.ConfigInterface, fileUtils fileutils.FileUtilsInterface, templateUtils templates.TemplatesInterface, gitClient gitclient.GitClientInterface, dbUtils db.DbInterface) NewCmdInterface {
	return &NewCmd{
		config:        config,
		dbUtils:       dbUtils,
		fileUtils:     fileUtils,
		templateUtils: templateUtils,
		gitClient:     gitClient,
		basePath:      "app",
	}
}

func (c *NewCmd) getServiceFilePath(serviceName, version string) string {
	return fmt.Sprintf("proto/%s/%s/%s_service.proto", serviceName, version, serviceName)

}
