package cmd

import (
	"github.com/darwishdev/devkit-cli/app/new"
	"github.com/darwishdev/devkit-cli/app/seed"
	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

type Command struct {
	config config.ConfigInterface
	newCmd new.NewCmdInterface

	seedCmd seed.SeedCmdInterface
}

func NewCommand(config config.ConfigInterface, newCmd new.NewCmdInterface, seedCmd seed.SeedCmdInterface) *Command {
	return &Command{
		config:  config,
		newCmd:  newCmd,
		seedCmd: seedCmd,
	}
}
func (c *Command) getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "devkit-cli",
		Short: "Boost your development workflow with devkit-cli.",
		Long:  `devkit-cli is a command-line interface designed to streamline the creation of new Go projects. It provides a solid foundation with essential backend functionalities out-of-the-box, allowing you to focus on building core features instead of writing repetitive boilerplate code.`,
	}
	rootCmd.AddCommand(c.getInitCmd())
	rootCmd.AddCommand(c.getSeedCmd())
	rootCmd.AddCommand(c.getNewCmd())
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return rootCmd
}
func (c *Command) getInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the devkit config file",
		Long:  `Creates a config file in the user's home directory.`,
		Run: func(cmd *cobra.Command, args []string) {
			// 	err := c.config.InitProjectConfig()
			// 	if err != nil {
			// 		fmt.Println("err initing the project config file:", err)
			// 		os.Exit(1)
			// 	}
		},
	}
}
func (c *Command) Execute() {
	rootCmd := c.getRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
