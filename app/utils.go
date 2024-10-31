package main

import (
	"crypto/rand"
	"fmt"
)

const alphanumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func secureRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = alphanumericChars[int(bytes[i])%len(alphanumericChars)]
	}

	return string(result), nil
}
