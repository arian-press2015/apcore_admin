package cmd

import (
	"fmt"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/utils/httpclient"
	"github.com/arian-press2015/apcore_admin/utils/table"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "Get list of users",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		httpParser := httpclient.NewHTTPParser(cfg)
		getUsers(cfg, httpParser, 0)
	},
}

func init() {
	usersCmd.AddCommand(listUsersCmd)
}

func getUsers(cfg *config.Config, parser *httpclient.HTTPParser, offset int) {
	url := fmt.Sprintf("%s/users?offset=%d&limit=10", cfg.BackendURL, offset)

	var responseBody UsersReponse

	err := parser.ParseRequest("GET", url, nil, &responseBody)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if len(responseBody.Data) == 0 {
		fmt.Println("No users found.")
		return
	}

	drawUsersTable(responseBody.Data)

	fmt.Println("Press 'n' for next page, 'p' for previous page, or any other key to exit.")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "n":
		getUsers(cfg, parser, offset+10)
	case "p":
		if offset-10 >= 0 {
			getUsers(cfg, parser, offset-10)
		} else {
			getUsers(cfg, parser, 0)
		}
	default:
		return
	}
}

func drawUsersTable(data []User){
	headers := []string{"No.", "ID", "Name", "Phone", "Verified"}
	var rows [][]string
	for index, user := range data {
		verified := "✘"
		if user.Verified {
			verified = "✔"
		}
		rowNumber := fmt.Sprintf("%d", index+1)
		row := []string{rowNumber, user.ID, user.FullName, user.Phone, verified}
		rows = append(rows, row)
	}

	table.PrintTable(headers, rows)
}

type UsersReponse struct {
	Data    []User `json:"data"`
	Message string `json:"message"`
	TrackID string `json:"trackId"`
}

type User struct {
	ID            string       `json:"id"`
	CreatedAt     string       `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
	DeletedAt     *string      `json:"deleted_at"`
	FullName      string       `json:"fullName"`
	Phone         string       `json:"phone"`
	ProfileImage  string       `json:"profile_image"`
	Nid           string       `json:"nid"`
	Verified      bool         `json:"verified"`
	Roles         []string     `json:"roles"`
	Notifications []string     `json:"notifications"`
	Subscription  Subscription `json:"subscription"`
}

type Subscription struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at"`
	UserID      string  `json:"userID"`
	Method      string  `json:"method"`
	SubjectType string  `json:"subjectType"`
}
