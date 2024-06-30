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

var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Get statistics",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpClient := httpclient.NewHTTPClient()
		tokenManager := token.NewTokenManager(cfg)
		getStatistics(cfg, httpClient, tokenManager)
	},
}

func getStatistics(cfg *config.Config, httpClient *http.Client, tokenManager *token.TokenManager) {
	url := fmt.Sprintf("%s/admin/statistics", cfg.BackendURL)
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
		fmt.Printf("Error fetching statistics: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
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
	UserCount     int `json:"user_count"`
	CustomerCount int `json:"customer_count"`
	TotalRevenue  int `json:"total_revenue"`
}
