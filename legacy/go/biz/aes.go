package biz

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// AES-256-CBC 암호화
func EncryptAES256CBC(plaintext, key []byte) (string, error) {
	// 블록 암호 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 패딩 적용
	paddedPlaintext := Pkcs7Padding(plaintext, block.BlockSize())

	// IV 생성
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// CBC 모드 인코더 생성
	mode := cipher.NewCBCEncrypter(block, iv)

	// 암호화
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// IV + 암호문을 결합하고 Base64로 인코딩
	finalCiphertext := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(finalCiphertext), nil
}

// AES-256-CBC 복호화 (Base64 문자열 입력)
func DecryptAES256CBC(ciphertextBase64 string, key []byte) ([]byte, error) {
	// Base64 디코딩
	ciphertextWithIV, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, err
	}

	// 블록 암호 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(ciphertextWithIV) < blockSize {
		return nil, fmt.Errorf("암호문이 너무 짧습니다")
	}

	// IV 추출
	iv := ciphertextWithIV[:blockSize]
	ciphertext := ciphertextWithIV[blockSize:]

	// 암호문 길이가 블록 크기의 배수인지 확인
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("암호문의 길이가 올바르지 않습니다")
	}

	// CBC 모드 디코더 생성
	mode := cipher.NewCBCDecrypter(block, iv)

	// 복호화
	plaintextPadded := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintextPadded, ciphertext)

	// PKCS7 패딩 제거
	plaintext, err := Pkcs7Unpadding(plaintextPadded)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptAES256GCM encrypts plaintext using AES-256 in GCM mode
func EncryptAES256GCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES256GCM decrypts ciphertext using AES-256 in GCM mode
func DecryptAES256GCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
