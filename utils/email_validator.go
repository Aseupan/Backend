package utils

import (
	"github.com/asaskevich/govalidator"
)

func ValidateEmail(email string) bool {
	return govalidator.IsEmail(email)
}
