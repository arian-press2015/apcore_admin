package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/token"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the admin CLI",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter phone: ")
		phone, _ := reader.ReadString('\n')
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		fmt.Print("Enter MFA code: ")
		mfaCode, _ := reader.ReadString('\n')

		phone = strings.TrimSpace(phone)
		password = strings.TrimSpace(password)
		mfaCode = strings.TrimSpace(mfaCode)

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

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	token, ok := result["token"]
	if !ok {
		return fmt.Errorf("no token found in response")
	}

	return tokenManager.SaveToken(token)
}
