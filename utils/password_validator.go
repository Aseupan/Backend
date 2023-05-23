package utils

import (
	"regexp"
)

func IsPasswordValid(password string) bool {
	// Check minimum length of eight characters
	if len(password) < 8 {
		return false
	}

	// Check if the password contains at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		return false
	}

	// Check if the password contains at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		return false
	}

	// Check if the password contains at least one digit
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(password) {
		return false
	}

	return true
}
