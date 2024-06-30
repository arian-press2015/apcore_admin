package cmd

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get list of users",
	Run: func(cmd *cobra.Command, args []string) {
		getUsers(0)
	},
}

func getUsers(offset int) {
	url := fmt.Sprintf("%s/admin/users?offset=%d&limit=10", backendURL, offset)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
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

	// Handle pagination
	fmt.Println("Press 'n' for next page, 'p' for previous page, or any other key to exit.")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "n":
		getUsers(offset + 10)
	case "p":
		if offset-10 >= 0 {
			getUsers(offset - 10)
		} else {
			getUsers(0)
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

func init() {
	rootCmd.AddCommand(usersCmd)
}
