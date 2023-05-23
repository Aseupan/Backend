package main

import (
	"fmt"
	"regexp"
)

func CheckPassword(password string) bool {
	// Check the length of the password
	if len(password) < 8 {
		return false
	}

	// Define the regular expression pattern
	pattern := "(?=[^a-z]*[a-z])(?=[^A-Z]*[A-Z])(?=[^0-9]*[0-9]).{8,}$"

	// Create a regular expression object
	regex := regexp.MustCompile(pattern)

	// Check if the password matches the pattern
	return regex.MatchString(password)
}

func main() {
	password := "MyPassword123" // Replace with the password you want to check

	// Define the regular expression pattern
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the password matches the regex pattern
	if regex.MatchString(password) {
		fmt.Println("Password meets the requirements")
	} else {
		fmt.Println("Password does not meet the requirements")
	}
}
