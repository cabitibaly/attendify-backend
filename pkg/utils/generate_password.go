package utils

import (
	"crypto/rand"
)

var passchars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%&*!"

func GeneratePassword(length int) (string, error) {
	pass := make([]byte, length)

	_, err := rand.Read(pass)

	if err != nil {
		return "", err
	}

	for i, b := range pass {
		pass[i] = passchars[b%byte(len(passchars))]
	}

	return string(pass), nil
}
