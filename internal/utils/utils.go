package utils

import (
	"crypto/rand"
	"encoding/hex"
	"reflect"
)

func FixedLengthRandomString(length int) (string, error) {
	// Generate a byte slice of half the length (since hex encoding doubles the size)
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// TODO: Not done yet
// ---------------------------------------
func SetDefaultValue[T any](value, defaultValue T) T {
	var empty T

	if reflect.DeepEqual(empty, value) {
		return defaultValue
	}

	return value
}

// ---------------------------------------
