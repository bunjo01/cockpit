package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/cockpit/utils"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	getAppConfigShortDesc = "Retrieve and display the configuration"
	getAppConfigLongDesc  = "This command retrieves the specific configuration\n" +
		"displays it in a nicely formatted way, and saves it to both YAML and JSON files.\n\n" +
		"Example:\n" +
		"get-group-config --org 'org' --name 'app_config' --version 'v1.0.0'"

	// Flag Constants
	flagName   = "name"
	flagOutput = "output"

	// Flag Shorthand Constants
	shortFlagName   = "n"
	shortFlagOutput = "o"

	// Flag Descriptions
	descName   = "Configuration name (required)"
	descOutput = "Output format (yaml or json)"
)

var (
	name              string
	appConfigResponse model.SingleConfigGroupResponse
	outputFormat      string
)

var GetSingleConfigGroupCmd = &cobra.Command{
	Use:   "group",
	Short: getAppConfigShortDesc,
	Long:  getAppConfigLongDesc,
	Run:   executeGetAppConfig,
}

func executeGetAppConfig(cmd *cobra.Command, args []string) {
	config := createRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	render.HandleSingleConfigGroupResponse(config.Response.(*model.SingleConfigGroupResponse), outputFormat)

	err = utils.SaveAppConfigResponseToFiles(config.Response.(*model.SingleConfigGroupResponse), outputFormat)
	if err != nil {
		log.Fatalf("Failed to save response to files: %v", err)
	}
}

func createRequestConfig() model.HTTPRequestConfig {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "GetConfigGroup")

	requestBody := model.ConfigReference{
		Organization: organization,
		Name:         name,
		Version:      version,
	}

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
		Response:    &appConfigResponse,
	}
}

func init() {
	GetSingleConfigGroupCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	GetSingleConfigGroupCmd.Flags().StringVarP(&name, flagName, shortFlagName, "", descName)
	GetSingleConfigGroupCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)
	GetSingleConfigGroupCmd.Flags().StringVarP(&outputFormat, flagOutput, shortFlagOutput, "yaml", descOutput)

	GetSingleConfigGroupCmd.MarkFlagRequired(flagOrganization)
	GetSingleConfigGroupCmd.MarkFlagRequired(flagName)
	GetSingleConfigGroupCmd.MarkFlagRequired(flagVersion)
}
