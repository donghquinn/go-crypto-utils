package biz

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateCustomRandomString generates a random string with letters (a-z, A-Z) and digits (0-9).
func GenerateCustomRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	// Define the character set (a-z, A-Z, 0-9)
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result string
	for i := 0; i < length; i++ {
		// Randomly select a character from the charSet
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", fmt.Errorf("error generating random index: %v", err)
		}
		result += string(charSet[index.Int64()])
	}

	return result, nil
}
