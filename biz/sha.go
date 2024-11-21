package biz

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
)

// HashData hashes data using the specified method
func HashData(data []byte, method string) ([]byte, error) {
	switch method {
	case "SHA-256":
		hash := sha256.Sum256(data)
		return hash[:], nil
	case "SHA-512":
		hash := sha512.Sum512(data)
		return hash[:], nil
	default:
		return nil, errors.New("unsupported hash method")
	}
}

// bytesRepeat is a helper function to repeat a byte
func BytesRepeat(b byte, count int) []byte {
	bytes := make([]byte, count)
	for i := range bytes {
		bytes[i] = b
	}
	return bytes
}
