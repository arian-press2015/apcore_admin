package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	field, _ := reader.ReadString('\n')

	field = strings.TrimSpace(field)

	return field
}
