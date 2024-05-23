package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	validateSchemaVersionShortDesc = "Validate a schema version"
	validateSchemaVersionLongDesc  = "This command validates a schema version with the given configuration.\n\n" +
		"Example:\n" +
		"validate-schema --org 'c12s' --schema_name 'schema' --version 'v1.0.0' --path '/path/to/config.yaml'"

	// Flag Constants
	flagOrganization = "org"
	flagSchemaName   = "schema_name"
	flagVersion      = "version"
	flagConfigPath   = "path"

	// Flag Shorthand Constants
	shortFlagOrganization = "o"
	shortFlagSchemaName   = "s"
	shortFlagVersion      = "v"
	shortFlagConfigPath   = "p"

	// Flag Descriptions
	descOrganization = "Organization name (required)"
	descSchemaName   = "Schema name (required)"
	descVersion      = "Schema version (required)"
	descConfigPath   = "Path to the YAML configuration file (required)"
)

var (
	organization string
	schemaName   string
	version      string
	configPath   string
)

var ValidateSchemaVersionCmd = &cobra.Command{
	Use:   "schema",
	Short: validateSchemaVersionShortDesc,
	Long:  validateSchemaVersionLongDesc,
	Run:   executeValidateSchemaVersion,
}

func executeValidateSchemaVersion(cmd *cobra.Command, args []string) {
	config := createValidateSchemaRequestConfig()

	err := utils.SendHTTPRequest(config)
	if err != nil {
		fmt.Printf("Error validating schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema validated successfully!")
}

func createValidateSchemaRequestConfig() model.HTTPRequestConfig {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	schemaDetails := model.SchemaDetails{
		Organization: organization,
		SchemaName:   schemaName,
		Version:      version,
	}

	requestBody := struct {
		SchemaDetails model.SchemaDetails `json:"schema_details"`
		Configuration string              `json:"configuration"`
	}{
		SchemaDetails: schemaDetails,
		Configuration: string(configData),
	}

	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "ValidateConfiguration")

	return model.HTTPRequestConfig{
		Method:      "GET",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: requestBody,
	}
}

func init() {
	ValidateSchemaVersionCmd.Flags().StringVarP(&organization, flagOrganization, shortFlagOrganization, "", descOrganization)
	ValidateSchemaVersionCmd.Flags().StringVarP(&schemaName, flagSchemaName, shortFlagSchemaName, "", descSchemaName)
	ValidateSchemaVersionCmd.Flags().StringVarP(&version, flagVersion, shortFlagVersion, "", descVersion)
	ValidateSchemaVersionCmd.Flags().StringVarP(&configPath, flagConfigPath, shortFlagConfigPath, "", descConfigPath)

	ValidateSchemaVersionCmd.MarkFlagRequired(flagOrganization)
	ValidateSchemaVersionCmd.MarkFlagRequired(flagSchemaName)
	ValidateSchemaVersionCmd.MarkFlagRequired(flagVersion)
	ValidateSchemaVersionCmd.MarkFlagRequired(flagConfigPath)
}