package ticket

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	ErrAesDataSize = errors.New("data size error")
)

type AESTicketEnDecoder struct {
	Key []byte
}

func NewAESTicket(key []byte) TicketEnDecoder {
	if len(key) == 0 {
		key = make([]byte, 32)
		_, _ = io.ReadFull(rand.Reader, key)
	}
	return &AESTicketEnDecoder{Key: key}
}

func (p *AESTicketEnDecoder) Encode(session *SessionData) (string, error) {
	if session.CreateTime == 0 {
		session.CreateTime = uint64(time.Now().Unix())
	}
	if data, err := AesEncrypt(p.Key, session.Marshal()); err != nil {
		return "", err
	} else {
		return base64.RawStdEncoding.EncodeToString(data), nil
	}
}

func (p *AESTicketEnDecoder) Decode(ticket string) (*SessionData, error) {
	rawData, err := base64.RawStdEncoding.DecodeString(ticket)
	if err != nil {
		return nil, err
	}
	if data, err := AesDecrypt(p.Key, rawData); err != nil {
		return nil, err
	} else {
		t := &SessionData{}
		if err := t.Unmarshal(data); err != nil {
			return nil, err
		}
		return t, nil
	}
}

func AesEncrypt(key, data []byte) ([]byte, error) {
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

func AesDecrypt(key, data []byte) (plaintext []byte, err error) {
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
