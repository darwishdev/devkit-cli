package new

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (c *NewCmd) ExecCmd(dir string, name string, args ...string) error {
	makeCmd := exec.Command(name, args...)
	makeCmd.Dir = dir
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	err := makeCmd.Run()
	if err != nil {
		return err
	}
	return nil

}

func (c *NewCmd) NewApi(cmd *cobra.Command, args []string) {
	cliConfig := c.config.GetCliConfig()
	appName := args[0]
	org, err := cmd.Flags().GetString("org")
	ctx := context.Background()
	if err != nil || org == "" {
		org = cliConfig.GitUser
	}
	bufUser, err := cmd.Flags().GetString("buf-user")
	if err != nil || bufUser == "" {
		bufUser = cliConfig.BufUser
	}
	err = c.gitClient.Fork(ctx, org, appName)
	if err != nil {
		log.Err(err).Msg("failed to fork the repo")
		os.Exit(1)
	}
	err = c.gitClient.Clone(ctx, org, appName)
	if err != nil {
		log.Err(err).Msg("failed to clone the repo but forked")
		os.Exit(1)
	}
	err = c.fileUtils.CopyFile(fmt.Sprintf("%s/.env.example", appName), fmt.Sprintf("%s/.env", appName))
	if err != nil {
		log.Err(err).Msg("copying .env.example")
		return
	}
	err = c.fileUtils.CopyFiles(fmt.Sprintf("%s/config/*.example.env", appName), func(s string) string { return strings.Replace(s, ".example", "", 1) })
	if err != nil {
		fmt.Println("Error copying example confif files:", err)
		return
	}
	err = c.fileUtils.ReplaceInDir(appName, map[string]string{
		"module github.com/darwishdev/devkit-api": fmt.Sprintf("module github.com/%s/%s", org, appName),
		"project_id = \"devkit-api\"":             fmt.Sprintf("project_id = \"%s\"", appName),
		"github.com/darwishdev/devkit-api":        fmt.Sprintf("github.com/%s/%s", org, appName),
		"buf.build/ahmeddarwish/devkit-api":       fmt.Sprintf("buf.build/%s/%s", bufUser, appName),
	})
	if err != nil {
		fmt.Println("Error replacing app name and repo :", err)
		return
	}
	// 5. Run make commands
	err = c.ExecCmd(appName, "make", "buf", "sqlc")
	if err != nil {
		fmt.Println("Error running make buf:", err)
		os.Exit(1)
	}
	err = c.ExecCmd(appName, "go", "mod", "tidy")
	if err != nil {
		fmt.Println("Error running go mod tidy:", err)
		os.Exit(1)
	}

	err = c.ExecCmd(appName, "devkit", "init")
	if err != nil {
		fmt.Println("error init the project:", err)
		os.Exit(1)
	}
	log.Debug().Str(bufUser, "new api from domain").Msg("domain")
}
