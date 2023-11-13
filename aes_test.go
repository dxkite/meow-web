package suda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAesDecrypt(t *testing.T) {
	token := &Token{
		ExpireAt: time.Now().Add(1 * time.Hour).Unix(),
		Value:    "dxkite",
	}
	vv, _ := json.Marshal(token)
	key := "12345678901234567890123456789012"
	val, err := AesEncrypt([]byte(key), vv)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(base64.RawURLEncoding.EncodeToString(val))
}
