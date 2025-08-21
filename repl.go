package main

import "strings"

// split user's input by spaces into slice of lowercase strings
func cleanInput(text string) []string {
	lowString := strings.ToLower(text)
	result := strings.Fields(lowString)
	return result
}
