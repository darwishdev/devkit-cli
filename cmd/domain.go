package cmd

import (
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain [domain_name]",
	Short: "Create a new domain",
	Long:  `Create a new domain in the current directory.`,
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]

		domainDirs := appFiles.CollectDomainDirecories(domainName)
		// 1. Create directories
		for key, dir := range domainDirs {
			dirName := dir.(string)
			if _, err := os.Stat(dirName); os.IsNotExist(err) {
				err := os.MkdirAll(dirName, 0755)
				if err != nil {
					fmt.Println("Error creating:", key, err)
					os.Exit(1)
				}
			}
		}
		// 2. Create files
		domainFiles := appFiles.CollectDomainFiles(domainName)
		// 1. Create directories
		for key, file := range domainFiles {
			fileName := file.(string)
			_, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error creating:", key, err)
				os.Exit(1)
			}

		}
		templates, err := appTemplates.LoadDomainTemplates(domainName)
		if err != nil {
			fmt.Println("err loading domain templates", err)
			os.Exit(1)
		}
		for key, template := range templates {
			file, ok := domainFiles[key]
			if !ok {
				fmt.Println("not found", key)
				continue
			}
			fileName := file.(string)
			err = os.WriteFile(fileName, template.Bytes(), 0644)
			if err != nil {
				fmt.Println("Error writing ", key, err)
				os.Exit(1)
			}
		}
		conf := appConfig.GetConfig()
		camelDomainName := strcase.ToCamel(domainName)
		lowerCamelDomainName := strcase.ToLowerCamel(domainName)
		usecaseImport := fmt.Sprintf("%sUsecase \"%s/app/%s/usecase\"", lowerCamelDomainName, conf.BasePath, domainName)
		usecaseField := fmt.Sprintf("%sUsecase %sUsecase.%sUsecaseInterface", lowerCamelDomainName, lowerCamelDomainName, camelDomainName)
		usecaseInstantiation := fmt.Sprintf("%sUsecase := %sUsecase.New%sUsecase(store , config)", lowerCamelDomainName, lowerCamelDomainName, camelDomainName)
		usecaseInjection := fmt.Sprintf("%sUsecase: %sUsecase,", lowerCamelDomainName, lowerCamelDomainName)
		err = appFiles.ReplaceMultiple(conf.ApiFilePath, map[string]string{
			"// USECASE_IMPORTS":        fmt.Sprintf("// USECASE_IMPORTS\n%s", usecaseImport),
			"// USECASE_FIELDS":         fmt.Sprintf("// USECASE_FIELDS\n%s", usecaseField),
			"// USECASE_INSTANTIATIONS": fmt.Sprintf("// USECASE_INSTANTIATIONS\n%s", usecaseInstantiation),
			"// USECASE_INJECTIONS":     fmt.Sprintf("// USECASE_INJECTIONS\n%s", usecaseInjection),
		})

		if err != nil {
			fmt.Println("Error replacing api", conf.ApiFilePath, err)
			os.Exit(1)
		}

		fmt.Println("Created domain :", domainName)
	},
}

func init() {
	newCmd.AddCommand(domainCmd)
}
