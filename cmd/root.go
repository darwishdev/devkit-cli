package cmd

import (
	"fmt"
	"os"

	"github.com/darwishdev/devkit-cli/config"
	"github.com/darwishdev/devkit-cli/fileutils"
	"github.com/darwishdev/devkit-cli/templates"
	"github.com/darwishdev/sqlseeder"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var appConfig = config.NewConfig()
var appTemplates = templates.NewTemplates(appConfig)

var appFiles = fileutils.NewFileUtils(appConfig)

var appSeeder = sqlseeder.NewSeeder(sqlseeder.SeederConfig{
	HashFunc: hashPassword,
})

func hashPassword(pass string) string {
	password, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(password)

}

var rootCmd = &cobra.Command{
	Use:   "devkit-cli",
	Short: "Boost your development workflow with devkit-cli.", // More concise and benefit-oriented
	Long: `devkit-cli is a command-line interface (CLI) designed to 
           streamline common development tasks and improve your productivity.

           Quickly generate new projects with pre-configured settings:
               devkit-cli new my-project

           Easily add features to your existing projects:
               devkit-cli add api user

           Automate repetitive tasks like code generation and testing:
               devkit-cli generate model User

           Get started now with 'devkit-cli --help' or visit 
           [your-project-website/docs] for detailed information.`, // Added examples and a call to action
	// Run: func(cmd *cobra.Command, args []string) { },
} // Config struct to hold configuration data
// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the devkit config file",
	Long:  `Creates a config file in the user's home directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := appConfig.Init()
		if err != nil {
			fmt.Println("err", err)
			os.Exit(1)
		}
		fmt.Println("Config file created:")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(initCmd)
}
