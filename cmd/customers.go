package cmd

import (
	"fmt"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/utils"
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
		httpParser := httpclient.NewHTTPParser(cfg)
		getCustomers(cfg, httpParser, 0)
	},
}

var createCustomerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpParser := httpclient.NewHTTPParser(cfg)
		createCustomer(cfg, httpParser)
	},
}

func init() {
	customersCmd.AddCommand(listCustomersCmd)
	customersCmd.AddCommand(createCustomerCmd)
}

func getCustomers(cfg *config.Config, parser *httpclient.HTTPParser, offset int) {
	url := fmt.Sprintf("%s/customers?offset=%d&limit=10", cfg.BackendURL, offset)

	var responseBody CustomersResponse

	err := parser.ParseRequest("GET", url, nil, &responseBody)
	if err != nil {
		fmt.Printf("%v\n", err)
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
		getCustomers(cfg, parser, offset+10)
	case "p":
		if offset-10 >= 0 {
			getCustomers(cfg, parser, offset-10)
		} else {
			getCustomers(cfg, parser, 0)
		}
	default:
		return
	}
}

func createCustomer(cfg *config.Config, parser *httpclient.HTTPParser) {
	name := utils.Prompt("Enter customer name: ")
	details := utils.Prompt("Enter customer details: ")
	phone := utils.Prompt("Enter customer phone: ")
	logo := utils.Prompt("Enter customer logo URL: ")

	customer := Customer{
		Name:    name,
		Details: details,
		Phone:   phone,
		Logo:    logo,
	}

	var responseBody CustomersResponse

	url := fmt.Sprintf("%s/customers", cfg.BackendURL)
	err := parser.ParseRequest("POST", url, customer, &responseBody)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println("Customer created successfully")
}

type CustomersResponse struct {
	Data    []Customer `json:"data"`
	Message string     `json:"message"`
	TrackID string     `json:"trackId"`
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
