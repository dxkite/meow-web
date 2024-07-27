package passwd

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

func NewHash(password string) (string, error) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := make([]byte, 16)
	if n, err := rd.Read(salt); err != nil {
		return "", err
	} else {
		salt = salt[:n]
	}
	sum := hmacSha256([]byte(password), salt)
	p1 := base64.RawURLEncoding.EncodeToString(sum)
	p2 := base64.RawURLEncoding.EncodeToString(salt)
	return p1 + "." + p2, nil
}

func VerifyHash(password, hash string) (bool, error) {
	p := strings.SplitN(hash, ".", 2)
	p1, p2 := p[0], p[1]

	sum, err := base64.RawURLEncoding.DecodeString(p1)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawURLEncoding.DecodeString(p2)
	if err != nil {
		return false, err
	}

	testSum := hmacSha256([]byte(password), salt)
	if bytes.Equal(sum, testSum) {
		return true, nil
	}

	return false, nil
}

func hmacSha256(message, key []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(message)
	return hash.Sum(nil)
}
