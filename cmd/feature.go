package cmd

import (
	// ... other imports ...
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ... other code ...

// featureCmd represents the feature command
var featureCmd = &cobra.Command{
	Use:   "feature [domain_name] [feature_name]",
	Short: "Create a new feature",
	Long:  `Create a new feature within a domain.`,
	Args:  cobra.ExactArgs(2), // Ensure exactly two arguments are provided
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]
		featureName := args[1]

		featureFiles := appFiles.CollectFeatureFiles(domainName, featureName)

		// config, err := readConfig()
		// if err != nil {
		// 	os.Exit(1)
		// }
		// featureFiles, err := collectFeatureFilePaths(domainName, featureName)
		// if err != nil {
		// 	fmt.Println("Error collecting file paths:", err)
		// 	os.Exit(1)
		// }
		for key, file := range featureFiles {
			fileName := file.(string)
			_, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error creating ", key, err)
				os.Exit(1)
			}

		}
		templates, err := appTemplates.LoadFeatureTemplates(domainName, featureName)
		if err != nil {
			fmt.Println("Error loading templates ", err)
			os.Exit(1)
		}
		for key, template := range templates {
			file, ok := featureFiles[key]
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
		config := appConfig.GetConfig()
		protoImport := fmt.Sprintf("import \"%s/%s/%s_%s.proto\"", config.ApiServiceName, config.ApiVersion, domainName, featureName)
		err = appFiles.ReplaceFile(config.ProtoServiceFilePath, "// INJECT IMPORTS", fmt.Sprintf("// INJECT IMPORTS\n%s;", protoImport))
		if err != nil {
			fmt.Println("Error replacing api", err)
			os.Exit(1)
		}

		// templateAttrs := map[string]string{
		// 	"ApiServiceName": config.ApiServiceName,
		// 	"ApiVersion":     config.ApiVersion,
		// }
		//
		// protoBuffer, err := LoadTemplate("templates/proto_base.tmpl", "base_proto", templateAttrs)
		// err = os.WriteFile(featureFiles["Proto"], protoBuffer.Bytes(), 0644)
		// if err != nil {
		// 	fmt.Println("Error writing proto:", err)
		// 	os.Exit(1)
		// }
		// adapterBuffer, err := LoadTemplate("templates/feature_adapter_base.tmpl", "feature_adapter_base", config)
		// err = os.WriteFile(featureFiles["Adapter"], adapterBuffer.Bytes(), 0644)
		// if err != nil {
		// 	fmt.Println("Error writing adapter", err)
		// 	os.Exit(1)
		// }
		// usecaseBuffer, err := LoadTemplate("templates/feature_usecase_base.tmpl", "feature_usecase_base", config)
		// err = os.WriteFile(featureFiles["Usecase"], usecaseBuffer.Bytes(), 0644)
		// if err != nil {
		// 	fmt.Println("Error writing usecase", err)
		// 	os.Exit(1)
		// }
		// repoBuffer, err := LoadTemplate("templates/feature_repo_base.tmpl", "feature_repo_base", config)
		// err = os.WriteFile(featureFiles["Repo"], repoBuffer.Bytes(), 0644)
		// if err != nil {
		// 	fmt.Println("Error writing usecase", err)
		// 	os.Exit(1)
		// }
		//
		fmt.Println("Created feature:", featureName, "in domain:", domainName)
	},
}

func init() {
	newCmd.AddCommand(featureCmd)
}
