package utils

import (
	"github.com/badoux/checkmail"
)

func EmailValidateFormat(email string) error {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}

	return nil
}

func EmailValidateDomain(email string) error {
	err := checkmail.ValidateHost(email)
	if err != nil {
		return err
	}

	return nil
}

func IsEmailValid(email string) bool {
	err := EmailValidateFormat(email)
	if err != nil {
		return false
	}

	err = EmailValidateDomain(email)
	if err != nil {
		return false
	}

	return true
}
