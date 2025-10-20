package main

import "fmt"

func formatCents(cents int) string {
	dollars := cents / 100
	remainder := cents % 100
	return fmt.Sprintf("$%d.%02d", dollars, remainder)
}
