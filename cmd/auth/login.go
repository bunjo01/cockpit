package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	shortLoginDescription = "Login into application"
	longLoginDescription  = "Input your username after that you will be prompted to input your password.\n" +
		"Your token will be saved in the token.txt file, which will be sent with all of your request headers.\n\n" +
		"Example:\n" +
		"login --username \"username\""
	tokenPath = "token.txt"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: shortLoginDescription,
	Long:  longLoginDescription,
	Run: func(cmd *cobra.Command, args []string) {
		password, err := utils.PromptForPassword()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		err = login(username, password)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Login successful!")
		fmt.Println()
	},
}

func init() {
	LoginCmd.Flags().StringVarP(&username, flagUsername, shortUsername, "", "Username for login")
	LoginCmd.MarkFlagRequired(flagUsername)
}

func login(username, password string) error {
	credentials := model.Credentials{
		Username: username,
		Password: password,
	}

	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loginURL := clients.BuildURL("core", "v1", "LoginUser")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, bytes.NewBuffer(credentialsJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode == http.StatusOK {
		var tokenResponse model.TokenResponse
		if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}

		if err := saveTokenToFile(tokenResponse.Token); err != nil {
			return fmt.Errorf("failed to save token: %v", err)
		}

		return nil
	}

	return fmt.Errorf("login failed: %s", bodyString)
}

func saveTokenToFile(token string) error {
	tokenFilePath := tokenPath
	return ioutil.WriteFile(tokenFilePath, []byte(token), 0600)
}
