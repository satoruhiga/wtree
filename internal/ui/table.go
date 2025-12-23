package ui

import (
	"fmt"
	"strings"
)

// PrintTable prints a table with the given headers and rows
func PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				// Strip ANSI codes for width calculation
				plainCell := stripAnsi(cell)
				if len(plainCell) > widths[i] {
					widths[i] = len(plainCell)
				}
			}
		}
	}

	// Print header
	for i, h := range headers {
		if i > 0 {
			fmt.Print("  ")
		}
		fmt.Printf("%-*s", widths[i], h)
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Print("  ")
			}
			// Pad considering ANSI codes
			plainLen := len(stripAnsi(cell))
			padding := widths[i] - plainLen
			if padding < 0 {
				padding = 0
			}
			fmt.Print(cell + strings.Repeat(" ", padding))
		}
		fmt.Println()
	}
}

// stripAnsi removes ANSI escape codes from a string
func stripAnsi(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}
