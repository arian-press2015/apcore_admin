package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "admin-cli",
	Short: "Admin CLI tool",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(customersCmd)
	rootCmd.AddCommand(statisticsCmd)
}
