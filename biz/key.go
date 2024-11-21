package biz

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func HashAESKey(bytestring []byte) (string, string) {
	// 키를 헥사 문자열로 인코딩
	hexKey := hex.EncodeToString(bytestring)
	fmt.Printf("헥사 키: %s\n", hexKey)

	// 키를 Base64 문자열로 인코딩
	base64Key := base64.StdEncoding.EncodeToString(bytestring)
	fmt.Printf("Base64 키: %s\n", base64Key)

	return hexKey, base64Key
}

// GenerateRandomAESKey generates a random AES key of specified byte length
func GenerateRandomAESKey(length int) ([]byte, error) {
	if length != 16 && length != 24 && length != 32 { // AES-128, AES-192, AES-256
		return nil, errors.New("invalid key length: must be 16, 24, or 32 bytes")
	}
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GenKey returns Hex and Base64 representations of the key
func GenKey(key []byte) (string, string) {
	return hex.EncodeToString(key), base64.StdEncoding.EncodeToString(key)
}
