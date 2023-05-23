package utils

import (
	"github.com/asaskevich/govalidator"
)

func IsEmailValid(email string) bool {
	return govalidator.IsEmail(email)
}
