package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func FixedLengthRandomString(length int) (string, error) {
	// Generate a byte slice of half the length (since hex encoding doubles the size)
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
