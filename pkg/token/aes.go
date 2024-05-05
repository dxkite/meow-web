package token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

var (
	ErrAesDataSize = errors.New("data size error")
)

func NewAesCrypto(key []byte) Crypto {
	return &aesCrypto{key: key}
}

type aesCrypto struct {
	key []byte
}

func (c *aesCrypto) Encrypt(data []byte) (encryptData []byte, err error) {
	return aesEncrypt(c.key, data)
}

func (c *aesCrypto) Decrypt(encryptData []byte) (data []byte, err error) {
	return aesDecrypt(c.key, data)
}

func aesEncrypt(key, data []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func aesDecrypt(key, data []byte) (plaintext []byte, err error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, ErrAesDataSize
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return
}
