package main

import (
	"context"
	"os"

	"github.com/darwishdev/devkit-cli/app/new"
	"github.com/darwishdev/devkit-cli/app/seed"
	"github.com/darwishdev/devkit-cli/cmd"
	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/darwishdev/devkit-cli/pkg/db"
	"github.com/darwishdev/devkit-cli/pkg/fileutils"
	"github.com/darwishdev/devkit-cli/pkg/gitclient"
	"github.com/darwishdev/devkit-cli/pkg/templates"
	"github.com/darwishdev/sqlseeder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	appConfig, err := config.NewConfig("$HOME/.config/devkitcli", "devkit", "./", "devkit")
	if err != nil {
		panic(err)
	}
	cliConfig := appConfig.GetCliConfig()
	fileUtils := fileutils.NewFileUtils()
	templateUtils := templates.NewTemplates("tmpls")
	sqlSeeder := sqlseeder.NewSeeder(sqlseeder.SeederConfig{
		HashFunc: func(pass string) string {
			password, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
			return string(password)

		},
	})
	gitClient := gitclient.NewGitClientRepo(context.Background(), cliConfig.GithubToken)
	dbUtils := db.NewDb()
	seedCmd := seed.NewSeedCmd(appConfig, fileUtils, sqlSeeder, dbUtils)
	newCmd := new.NewNewCmd(appConfig, fileUtils, templateUtils, gitClient, dbUtils)
	command := cmd.NewCommand(appConfig, newCmd, seedCmd)
	command.Execute()
}
