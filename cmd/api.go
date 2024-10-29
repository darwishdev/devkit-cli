package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/darwishdev/devkit-cli/gitclient"
	"github.com/spf13/cobra"
)

// ... other code ...

// Helper function to copy files
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api [app-name] [base-git-url] [base-buf-url]",
	Short: "Create a new API application",
	Long:  `Create a new API application with the specified parameters.`,
	Args:  cobra.ExactArgs(3), // Ensure exactly three arguments are provided
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		baseGitURL := args[1]
		baseBufURL := args[2]
		ctx := context.Background()
		conf := appConfig.GetConfig()

		gitClient := gitclient.NewGitClientRepo(ctx, conf.GithubToken, "darwishdev", "devkit-api")
		fmt.Println("giorg", baseGitURL)
		err := gitClient.Fork(ctx, baseGitURL, appName)

		// ... (Your logic to create the API application) ...
		// 1. Run git clone

		if err != nil {
			fmt.Println("Error forking the git repo:", err)
			os.Exit(1)
		}
		gitCloneCmd := exec.Command("git", "clone", fmt.Sprintf(fmt.Sprintf("git@github.com:%s/%s.git", baseGitURL, appName)))

		gitCloneCmd.Stdout = os.Stdout
		gitCloneCmd.Stderr = os.Stderr
		err = gitCloneCmd.Run()
		if err != nil {
			fmt.Println("Error running git clone:", err)
			os.Exit(1)
		}
		err = appFiles.CopyFile(fmt.Sprintf("%s/.env.example", appName), fmt.Sprintf("%s/.env", appName))
		if err != nil {
			fmt.Println("Error copying .env.example:", err)
			return
		}
		err = appFiles.CopyFiles(fmt.Sprintf("%s/config/*.example.env", appName), func(s string) string { return strings.Replace(s, ".example", "", 1) })
		if err != nil {
			fmt.Println("Error copying example confif files:", err)
			return
		}
		err = appFiles.ReplaceInDir(appName, map[string]string{
			"module github.com/darwishdev/devkit-api": fmt.Sprintf("module github.com/%s/%s", baseGitURL, appName),
			"project_id = \"devkit-api\"":             fmt.Sprintf("project_id = \"%s\"", appName),
			"github.com/darwishdev/devkit-api":        fmt.Sprintf("github.com/%s/%s", baseGitURL, appName),
			"buf.build/ahmeddarwish/devkit-api":       filepath.Join(baseBufURL, appName),
		})
		if err != nil {
			fmt.Println("Error replacing app name and repo :", err)
			return
		}
		// appDir := filepath.Join(config.Data["base_path"].(string), "app", appName)
		// 5. Run make commands
		makeBufCmd := exec.Command("make", "buf")
		makeBufCmd.Dir = appName
		makeBufCmd.Stdout = os.Stdout
		makeBufCmd.Stderr = os.Stderr
		err = makeBufCmd.Run()
		if err != nil {
			fmt.Println("Error running make buf:", err)
			os.Exit(1)
		}

		makeSqlcCmd := exec.Command("make", "sqlc")
		makeSqlcCmd.Dir = appName
		makeSqlcCmd.Stdout = os.Stdout
		makeSqlcCmd.Stderr = os.Stderr
		err = makeSqlcCmd.Run()
		if err != nil {
			fmt.Println("Error running make sqlc:", err)
			os.Exit(1)
		}

		goModTidyCmd := exec.Command("go", "mod", "tidy")
		goModTidyCmd.Dir = appName
		goModTidyCmd.Stdout = os.Stdout
		goModTidyCmd.Stderr = os.Stderr
		err = goModTidyCmd.Run()
		if err != nil {
			fmt.Println("Error running go mod tidy:", err)
			os.Exit(1)
		}

		fmt.Println("Created API application:", appName, "with Git URL:", baseGitURL, "and Buf URL:", baseBufURL)
	},
}

func init() {
	newCmd.AddCommand(apiCmd)
}

// ... other code ...
