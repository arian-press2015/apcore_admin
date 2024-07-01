package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/spf13/cobra"
)

var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Get statistics",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient(cfg)
		getStatistics(cfg, httpClient)
	},
}

func getStatistics(cfg *config.Config, httpClient *httpclient.HTTPClient) {
	url := fmt.Sprintf("%s/admin/statistics", cfg.BackendURL)
	resp, err := httpClient.MakeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()

	var stats Statistics
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	fmt.Printf("Statistics: %+v\n", stats)
}

type Statistics struct {
	UserCount     int `json:"user_count"`
	CustomerCount int `json:"customer_count"`
	TotalRevenue  int `json:"total_revenue"`
}
