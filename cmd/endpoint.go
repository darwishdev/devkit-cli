/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/spf13/cobra"
)

// ... other code ...

func getEndpointCmd(config config.ConfigInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "endpoint [domain] [feature] [endpoint]",
		Short: "Create a new endpoint",
		Long:  `Create a new endpoint within a feature and domain.`,
		Args:  cobra.ExactArgs(3), // Ensure exactly three arguments are provided
		Run: func(cmd *cobra.Command, args []string) {
			// 	domainName := args[0]
			// 	featureName := args[1]
			// 	endpointName := args[2]
			//
			// 	endpointFiles := appFiles.CollectFeatureFiles(domainName, featureName)
			// 	data := appConfig.WithEndpoint(domainName, featureName, endpointName)
			// 	templates, err := appTemplates.LoadEndpointTemplates(domainName, featureName, endpointName, data)
			//
			// 	if err != nil {
			// 		fmt.Println("err loading endpoint templates", err)
			// 		os.Exit(1)
			// 	}
			//
			// 	injecteionFiles := map[string]string{
			// 		"adapterinjector": fmt.Sprintf("app/%s/adapter/adapter.go", featureName),
			// 		"usecaseinjector": fmt.Sprintf("app/%s/usecase/usecase.go", featureName),
			// 		"repoinjector":    fmt.Sprintf("app/%s/repo/repo.go", featureName),
			// 	}
			// 	for key, template := range templates {
			// 		file, ok := endpointFiles[key]
			// 		if ok {
			// 			fileName := file.(string)
			// 			err := appFiles.AppendToFile(fileName, template)
			// 			if err != nil {
			// 				fmt.Println("Error writing ", key, err)
			// 				os.Exit(1)
			// 			}
			// 			continue
			//
			// 		}
			// 		injectioFile, ok := injecteionFiles[key]
			// 		if !ok {
			// 			continue
			// 		}
			// 		err = appFiles.ReplaceFile(injectioFile, "// INJECT INTERFACE", fmt.Sprintf("// INJECT INTERFACE\n%s", template.String()))
			// 		if err != nil {
			// 			fmt.Println("Error replacing api", err)
			// 			os.Exit(1)
			// 		}
			//
			// 	}
			// 	config := appConfig.GetConfig()
			// 	endpointFunction := strcase.ToCamel(fmt.Sprintf("%s_%s", featureName, endpointName))
			// 	grpcMehodName := fmt.Sprintf("rpc %s(%sRequest) returns (%sResponse) {} ", endpointFunction, endpointFunction, endpointFunction)
			// 	err = appFiles.ReplaceFile(config.ProtoServiceFilePath, "// INJECT METHODS", fmt.Sprintf("// INJECT METHODS\n%s", grpcMehodName))
			// 	if err != nil {
			// 		fmt.Println("Error replacing api", err)
			// 		os.Exit(1)
			// 	}
			// 	fmt.Println("Created endpoints:", data["EndpointName"], "in feature:", featureName, "in domain:", domainName)
		},
	}

}
