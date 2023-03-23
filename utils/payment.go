package utils

import (
	"math/rand"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func IntToRupiah(value int64) string {
	printer := message.NewPrinter(language.Indonesian)
	return printer.Sprintf("%d", value)
}

func RandomOrderID() string {
	n := 10
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
