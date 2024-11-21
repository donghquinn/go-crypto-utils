package test

import (
	"encoding/json"
	"log"
	"testing"

	"org.donghyuns.com/secure/keygen/biz"
)

type testInterface struct {
	Email string `json:"email"`
}

func TestAes(t *testing.T) {
	newKey, _ := biz.GenerateRandomAES256Key(32)
	hexKey, base64Key := biz.GenKey(newKey)

	log.Printf("newKey: %v", string(base64Key))

	text := testInterface{Email: "test@exampel.com"}

	// 구조체를 JSON으로 직렬화
	plaintextBytes, err := json.Marshal(text)
	if err != nil {
		t.Fatal("구조체 직렬화 오류:", err)
		return
	}

	encryptedString, encryptErr := biz.EncryptAES256CBC(plaintextBytes, newKey)

	if encryptErr != nil {
		t.Fatal("암호화 오류:", encryptErr)
		return
	}

	decryptedString, decryptErr := biz.DecryptAES256CBC(encryptedString, newKey)

	if decryptErr != nil {
		t.Fatal("복호화 오류:", decryptErr)
		return
	}

	t.Log("암호키 (Hex):", hexKey)
	t.Log("암호키 (Base64):", base64Key)
	t.Log("암호화된 문자: ", string(encryptedString))
	t.Log("복호화된 평문:", string(decryptedString))

	if string(plaintextBytes) != string(decryptedString) {
		t.Errorf("Decrypted string is not equal to original string")
	}
}
