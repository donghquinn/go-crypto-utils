package biz

import (
	"bytes"
	"fmt"
)

// PKCS7 패딩 추가
func Pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 패딩 제거
func Pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("데이터가 비어 있습니다")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, fmt.Errorf("잘못된 패딩")
	}
	for _, v := range data[length-padding:] {
		if int(v) != padding {
			return nil, fmt.Errorf("잘못된 패딩")
		}
	}
	return data[:length-padding], nil
}
