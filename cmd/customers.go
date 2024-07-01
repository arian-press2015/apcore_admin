package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/arian-press2015/apcore_admin/utils/table"
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
		httpClient := httpclient.NewHTTPClient(cfg)
		getCustomers(cfg, httpClient, 0)
	},
}

var createCustomerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient(cfg)
		createCustomer(cfg, httpClient)
	},
}

func init() {
	customersCmd.AddCommand(listCustomersCmd)
	customersCmd.AddCommand(createCustomerCmd)
}

func getCustomers(cfg *config.Config, httpClient *httpclient.HTTPClient, offset int) {
	url := fmt.Sprintf("%s/customers?offset=%d&limit=10", cfg.BackendURL, offset)
	resp, err := httpClient.MakeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()

	var responseBody struct {
		Data    []Customer `json:"data"`
		Message string     `json:"message"`
		TrackID string     `json:"trackId"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if len(responseBody.Data) == 0 {
		fmt.Println("No customers found.")
		return
	}

	headers := []string{"No.", "ID", "Name", "Phone", "Is Active", "Is Disabled"}
	var rows [][]string
	for index, customer := range responseBody.Data {
		isActive := "✘"
		if customer.IsActive {
			isActive = "✔"
		}
		isDisabled := "✘"
		if customer.IsDisabled {
			isDisabled = "✔"
		}
		rowNumber := fmt.Sprintf("%d", index+1)
		row := []string{rowNumber, customer.ID, customer.Name, customer.Phone, isActive, isDisabled}
		rows = append(rows, row)
	}

	table.PrintTable(headers, rows)

	fmt.Println("Press 'n' for next page, 'p' for previous page, or any other key to exit.")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "n":
		getCustomers(cfg, httpClient, offset+10)
	case "p":
		if offset-10 >= 0 {
			getCustomers(cfg, httpClient, offset-10)
		} else {
			getCustomers(cfg, httpClient, 0)
		}
	default:
		return
	}
}

func createCustomer(cfg *config.Config, httpClient *httpclient.HTTPClient) {
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
	resp, err := httpClient.MakeRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Customer created successfully")
}

type Customer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Details    string `json:"details"`
	Phone      string `json:"phone"`
	Logo       string `json:"logo"`
	IsActive   bool   `json:"isActive"`
	IsDisabled bool   `json:"isDisabled"`
}
