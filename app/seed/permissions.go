package seed

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func (c *SeedCmd) SeedPermissions() {
	conf, err := c.config.GetProjectConfig()
	if err != nil {
		log.Err(err).Msg("can't load the project config")
	}
	serviceFileName := fmt.Sprintf("proto/%s/%s/%s_service.proto", conf.ApiServiceName, conf.ApiVersion, conf.ApiServiceName)
	file, err := os.Open(serviceFileName)
	// Regex to match rpc function names in proto file
	rpcRegex := regexp.MustCompile(`rpc\s+(\w+)\s*\(`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains an RPC function
		if matches := rpcRegex.FindStringSubmatch(line); len(matches) > 1 {
			functionName := matches[1]
			// Derive permission details
			permissionName := formatPermissionName(functionName)
			permissionDescription := "Permission to " + strings.ToLower(permissionName)
			// Print the SQL insert statement line
			fmt.Printf("%s\t%s\t%s\t", permissionName, functionName, permissionDescription)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

// formatPermissionName formats the function name into a more human-readable permission name.
func formatPermissionName(functionName string) string {
	// Convert camel case to spaced words and lowercase
	re := regexp.MustCompile(`[A-Z][a-z]+`)
	words := re.FindAllString(functionName, -1)
	return strings.ToLower(strings.Join(words, " "))
}
