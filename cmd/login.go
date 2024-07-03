package cmd

import (
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
		httpParser := httpclient.NewHTTPParser(cfg)

		phone := utils.Prompt("Enter phone: ")
		password := utils.Prompt("Enter password: ")
		totp := utils.Prompt("Enter MFA code: ")

		err := login(httpParser, cfg, phone, password, totp)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return
		}
		fmt.Println("Login successful")
	},
}

func login(parser *httpclient.HTTPParser, cfg *config.Config, phone, password, totp string) error {
	url := fmt.Sprintf("%s/admin/auth", cfg.BackendURL)
	loginParams := LoginParams{Phone: phone, Password: password, Totp: totp}

	var responseBody LoginResponse

	err := parser.ParseUnauthenticatedRequest("POST", url, &loginParams, &responseBody)
	if err != nil {
		return err
	}

	token := responseBody.Data.Token
	if token == "" {
		return fmt.Errorf("no token found in response")
	}

	return parser.Client.TokenManager.SaveToken(token)
}

type LoginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
	TrackID string `json:"trackId"`
}

type LoginParams struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Totp     string `json:"totp"`
}
