package token

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestAesDecrypt(t *testing.T) {
	token := strconv.FormatInt(time.Now().Add(1*time.Hour).Unix(), 10)
	key := "12345678901234567890123456789012"
	val, err := aesEncrypt([]byte(key), []byte(token))
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(base64.RawURLEncoding.EncodeToString(val))

	decode, err := aesDecrypt([]byte(key), val)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(decode)

	// 创建 token
	tok := &BinaryToken{
		Id:       uint64(time.Now().Add(1 * time.Hour).Unix()),
		ExpireAt: uint64(time.Now().Add(1 * time.Hour).Unix()),
	}

	enc := NewAesCrypto([]byte(key))

	tokStr, err := tok.Encrypt(enc)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(tokStr)

	tokGet := &BinaryToken{}

	if err := tokGet.Decrypt(tokStr, enc); err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(tokGet)
}
