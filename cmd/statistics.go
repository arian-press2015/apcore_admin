package cmd

import (
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
		httpParser := httpclient.NewHTTPParser(cfg)
		getStatistics(cfg, httpParser)
	},
}

func getStatistics(cfg *config.Config, parser *httpclient.HTTPParser) {
	url := fmt.Sprintf("%s/admin/statistics", cfg.BackendURL)

	var stats Statistics

	err := parser.ParseRequest("GET", url, nil, &stats)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Statistics: %+v\n", stats)
}

type Statistics struct {
	UserCount     int `json:"user_count"`
	CustomerCount int `json:"customer_count"`
	TotalRevenue  int `json:"total_revenue"`
}
