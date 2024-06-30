package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/token"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get list of users",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient()
		tokenManager := token.NewTokenManager(cfg)
		getUsers(cfg, httpClient, tokenManager, 0)
	},
}

func getUsers(cfg *config.Config, httpClient *http.Client, tokenManager *token.TokenManager, offset int) {
	url := fmt.Sprintf("%s/admin/users?offset=%d&limit=10", cfg.BackendURL, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	err = tokenManager.AuthenticateRequest(req)
	if err != nil {
		fmt.Printf("Authentication error: %v\n", err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error: %s\n", string(body))
		return
	}

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}

	fmt.Println("Press 'n' for next page, 'p' for previous page, or any other key to exit.")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "n":
		getUsers(cfg, httpClient, tokenManager, offset+10)
	case "p":
		if offset-10 >= 0 {
			getUsers(cfg, httpClient, tokenManager, offset-10)
		} else {
			getUsers(cfg, httpClient, tokenManager, 0)
		}
	default:
		return
	}
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
