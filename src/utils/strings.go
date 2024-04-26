package utils

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/dlclark/regexp2"
)

func StringReplaceAll(reg, input, replacement string) (string, error) {
	r, err := regexp2.Compile(reg, 0)
	if err != nil {
		return "", err
	}
	v, err := r.Replace(input, replacement, -1, -1)
	if err != nil {
		return "", err
	}
	return v, nil
}

func GenerateRequestId() string {
	b := make([]byte, 16)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, err := rd.Read(b)
	if err != nil {
		return ""
	}
	return "req_" + base64.RawURLEncoding.EncodeToString(b)
}

func InStringSlice(s string, arr []string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
