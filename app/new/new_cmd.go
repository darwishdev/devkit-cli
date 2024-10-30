package new

import (
	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/darwishdev/devkit-cli/pkg/fileutils"
	"github.com/darwishdev/devkit-cli/pkg/gitclient"
	"github.com/darwishdev/devkit-cli/pkg/templates"
	"github.com/spf13/cobra"
)

type NewCmdInterface interface {
	NewApi(cmd *cobra.Command, args []string)
}
type NewCmd struct {
	config        config.ConfigInterface
	fileUtils     fileutils.FileUtilsInterface
	templateUtils templates.TemplatesInterface
	gitClient     gitclient.GitClientInterface
}

func NewNewCmd(config config.ConfigInterface, fileUtils fileutils.FileUtilsInterface, templateUtils templates.TemplatesInterface, gitClient gitclient.GitClientInterface) NewCmdInterface {
	return &NewCmd{
		config:        config,
		fileUtils:     fileUtils,
		templateUtils: templateUtils,
		gitClient:     gitClient,
	}
}
