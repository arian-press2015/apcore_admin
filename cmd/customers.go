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

var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Manage customers",
}

var listCustomersCmd = &cobra.Command{
	Use:   "list",
	Short: "Get list of customers",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient()
		tokenManager := token.NewTokenManager(cfg)
		getCustomers(cfg, httpClient, tokenManager, 0)
	},
}

var createCustomerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient()
		tokenManager := token.NewTokenManager(cfg)
		createCustomer(cfg, httpClient, tokenManager)
	},
}

func init() {
	customersCmd.AddCommand(listCustomersCmd)
	customersCmd.AddCommand(createCustomerCmd)
}

func getCustomers(cfg *config.Config, httpClient *http.Client, tokenManager *token.TokenManager, offset int) {
	url := fmt.Sprintf("%s/customers?offset=%d&limit=10", cfg.BackendURL, offset)
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
		fmt.Printf("Error fetching customers: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error: %s\n", string(body))
		return
	}

	var customers []Customer
	err = json.NewDecoder(resp.Body).Decode(&customers)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	for _, customer := range customers {
		fmt.Printf("ID: %s, Name: %s, Phone: %s\n", customer.ID, customer.Name, customer.Phone)
	}

	fmt.Println("Press 'n' for next page, 'p' for previous page, or any other key to exit.")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "n":
		getCustomers(cfg, httpClient, tokenManager, offset+10)
	case "p":
		if offset-10 >= 0 {
			getCustomers(cfg, httpClient, tokenManager, offset-10)
		} else {
			getCustomers(cfg, httpClient, tokenManager, 0)
		}
	default:
		return
	}
}

func createCustomer(cfg *config.Config, httpClient *http.Client, tokenManager *token.TokenManager) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter customer name: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("Enter customer details: ")
	details, _ := reader.ReadString('\n')
	fmt.Print("Enter customer phone: ")
	phone, _ := reader.ReadString('\n')
	fmt.Print("Enter customer logo URL: ")
	logo, _ := reader.ReadString('\n')

	customer := Customer{
		Name:    strings.TrimSpace(name),
		Details: strings.TrimSpace(details),
		Phone:   strings.TrimSpace(phone),
		Logo:    strings.TrimSpace(logo),
	}

	data, err := json.Marshal(customer)
	if err != nil {
		fmt.Printf("Error marshaling customer: %v\n", err)
		return
	}

	url := fmt.Sprintf("%s/customers", cfg.BackendURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	err = tokenManager.AuthenticateRequest(req)
	if err != nil {
		fmt.Printf("Authentication error: %v\n", err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error creating customer: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error: %s\n", string(body))
		return
	}

	fmt.Println("Customer created successfully")
}

type Customer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Details string `json:"details"`
	Phone   string `json:"phone"`
	Logo    string `json:"logo"`
}
