package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var backendURL string

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
	// Implement the login logic here by making an HTTP request to the backend
	// Store the token if login is successful
	// Example:
	// token, err := authenticate(username, password, mfaCode)
	// if err != nil {
	//     return err
	// }
	// saveToken(token)
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&backendURL, "backend", "b", "http://localhost:8080", "Backend URL")
}
