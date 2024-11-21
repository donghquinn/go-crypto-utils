package main

import (
	"encoding/json"
	"fmt"

	"org.donghyuns.com/secure/keygen/biz"
)

type testInterface struct {
	Email string `json:"email"`
}

func main() {
	// AES-256 키 생성
	newKey, _ := biz.GenerateRandomAES256Key(32)

	// 키를 헥사 문자열과 Base64 문자열로 인코딩
	hexKey, base64Key := biz.GenKey(newKey)

	text := testInterface{Email: "test@exampel.com"}

	// 구조체를 JSON으로 직렬화
	plaintextBytes, err := json.Marshal(text)
	if err != nil {
		fmt.Println("구조체 직렬화 오류:", err)
		return
	}

	encryptedString, encryptErr := biz.EncryptAES256CBC(plaintextBytes, newKey)

	if encryptErr != nil {
		fmt.Println("암호화 오류:", encryptErr)
		return
	}

	decryptedString, decryptErr := biz.DecryptAES256CBC(encryptedString, newKey)

	if decryptErr != nil {
		fmt.Println("복호화 오류:", decryptErr)
		return
	}

	fmt.Println("암호키 (Hex):", hexKey)
	fmt.Println("암호키 (Base64):", base64Key)
	fmt.Println("암호화된 문자: ", string(encryptedString))
	fmt.Println("복호화된 평문:", string(decryptedString))
}
