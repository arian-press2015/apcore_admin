package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Get statistics",
	Run: func(cmd *cobra.Command, args []string) {
		getStatistics()
	},
}

func getStatistics() {
	url := fmt.Sprintf("%s/admin/statistics", backendURL)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching statistics: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Error: %s\n", string(body))
		return
	}

	var stats Statistics
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	fmt.Printf("Statistics: %+v\n", stats)
}

type Statistics struct {
	UserCount      int `json:"user_count"`
	CustomerCount  int `json:"customer_count"`
	TotalRevenue   int `json:"total_revenue"`
}

func init() {
	rootCmd.AddCommand(statisticsCmd)
}
