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
}
