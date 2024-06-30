package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/token"
	"github.com/arian-press2015/apcore_admin/utils"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the admin CLI",
	Run: func(cmd *cobra.Command, args []string) {
		phone := utils.Prompt("Enter phone: ")
		password := utils.Prompt("Enter password: ")
		mfaCode := utils.Prompt("Enter MFA code: ")

		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient()
		tokenManager := token.NewTokenManager(cfg)

		err := login(cfg, httpClient, tokenManager, phone, password, mfaCode)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return
		}
		fmt.Println("Login successful")
	},
}

func login(cfg *config.Config, httpClient *http.Client, tokenManager *token.TokenManager, phone, password, totp string) error {
	url := fmt.Sprintf("%s/admin/auth", cfg.BackendURL)
	data := map[string]string{
		"phone":    phone,
		"password": password,
		"totp":     totp,
	}
	jsonData, _ := json.Marshal(data)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: %s", string(body))
	}

	var responseBody struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Message string `json:"message"`
		TrackID string `json:"trackId"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	token := responseBody.Data.Token
	if token == "" {
		return fmt.Errorf("no token found in response")
	}

	return tokenManager.SaveToken(token)
}
