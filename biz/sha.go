package biz

import "crypto/sha256"

// HashSHA256 hashes data using SHA-256
func HashSHA256(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return hash[:], nil
}

// bytesRepeat is a helper function to repeat a byte
func bytesRepeat(b byte, count int) []byte {
	bytes := make([]byte, count)
	for i := range bytes {
		bytes[i] = b
	}
	return bytes
}
