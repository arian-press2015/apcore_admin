package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/utils"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the admin CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient(cfg)

		phone := utils.Prompt("Enter phone: ")
		password := utils.Prompt("Enter password: ")
		totp := utils.Prompt("Enter MFA code: ")

		err := login(httpClient, cfg, phone, password, totp)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return
		}
		fmt.Println("Login successful")
	},
}

func login(httpClient *httpclient.HTTPClient, cfg *config.Config, phone, password, totp string) error {
	url := fmt.Sprintf("%s/admin/auth", cfg.BackendURL)
	loginParams := LoginParams{Phone: phone, Password: password, Totp: totp}
	body, err := json.Marshal(loginParams)
	if err != nil {
		return fmt.Errorf("error marshaling login parameters: %v", err)
	}

	resp, err := httpClient.MakeLoginRequest(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var responseBody struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Message string `json:"message"`
		TrackID string `json:"trackId"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	token := responseBody.Data.Token
	if token == "" {
		return fmt.Errorf("no token found in response")
	}

	return httpClient.TokenManager.SaveToken(token)
}

type LoginParams struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Totp     string `json:"totp"`
}
