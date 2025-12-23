package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks the user for confirmation with a y/N prompt
func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N] ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// ConfirmDefault asks the user for confirmation with a Y/n prompt (default yes)
func ConfirmDefault(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [Y/n] ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return true
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" {
		return true
	}
	return response != "n" && response != "no"
}
