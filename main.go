package main

import (
	"fmt"
	"os"

	"github.com/arian-press2015/apcore_admin/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
