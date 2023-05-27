package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

func SendCodeToEmail(email, code string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_EMAIL_PASSWORD")
	to := email

	// Message to be sent
	subject := "One-Time Code"
	body := fmt.Sprintf("Your one-time code is: %s", code)
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// SMTP configuration
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

func GenerateRandomCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
