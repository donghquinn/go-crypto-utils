package biz

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenKey(bytestring []byte) (string, string, string) {
	// 키를 헥사 문자열로 인코딩
	hexKey := hex.EncodeToString(bytestring)
	fmt.Printf("헥사 키: %s\n", hexKey)

	// 키를 Base64 문자열로 인코딩
	base64Key := base64.StdEncoding.EncodeToString(bytestring)
	fmt.Printf("Base64 키: %s\n", base64Key)

	return string(bytestring), hexKey, base64Key
}

// GenerateRandomAES256Key generates a secure random 32-byte key for AES-256
func GenerateRandomAES256Key(keyLength int64) ([]byte, error) {
	key := make([]byte, keyLength) // 32 bytes * 8 = 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
