package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"net/http"
	"bytes"

	"github.com/spf13/cobra"
)

var backendURL string
var tokenFile = ".auth_token"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the admin CLI",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		fmt.Print("Enter MFA code: ")
		mfaCode, _ := reader.ReadString('\n')

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)
		mfaCode = strings.TrimSpace(mfaCode)

		err := login(username, password, mfaCode)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return
		}
		fmt.Println("Login successful")
	},
}

func login(username, password, mfaCode string) error {
	url := fmt.Sprintf("%s/login", backendURL)
	data := map[string]string{
		"username": username,
		"password": password,
		"mfa_code": mfaCode,
	}
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
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

	// Save token to file
	err = ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&backendURL, "backend", "b", "http://localhost:8080", "Backend URL")
}
