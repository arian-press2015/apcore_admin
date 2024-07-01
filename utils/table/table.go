package table

import (
    "github.com/olekukonko/tablewriter"
    "os"
)

func PrintTable(headers []string, rows [][]string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader(headers)
    table.AppendBulk(rows)
    table.Render()
}