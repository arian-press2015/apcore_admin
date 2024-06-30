package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
		getCustomers(0)
	},
}

var createCustomerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	Run: func(cmd *cobra.Command, args []string) {
		createCustomer()
	},
}

func getCustomers(offset int) {
	url := fmt.Sprintf("%s/admin/customers?offset=%d&limit=10", backendURL, offset)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching customers: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
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

	// Handle pagination
	fmt.Println("Press 'n' for next page, 'p' for previous page, or any other key to exit.")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "n":
		getCustomers(offset + 10)
	case "p":
		if offset-10 >= 0 {
			getCustomers(offset - 10)
		} else {
			getCustomers(0)
		}
	default:
		return
	}
}

func createCustomer() {
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

	url := fmt.Sprintf("%s/admin/customers", backendURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error creating customer: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
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

func init() {
	rootCmd.AddCommand(customersCmd)
	customersCmd.AddCommand(listCustomersCmd)
	customersCmd.AddCommand(createCustomerCmd)
}
