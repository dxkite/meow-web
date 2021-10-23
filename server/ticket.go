package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type TicketProvider interface {
	EncodeTicket(uin uint64) (string, error)
	DecodeTicket(ticket string) (*Ticket, error)
}

type Ticket struct {
	Uin        uint64
	CreateTime uint64
}

func (t *Ticket) Marshal() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf, t.Uin)
	binary.BigEndian.PutUint64(buf[8:], t.CreateTime)
	return buf
}

func (t *Ticket) Unmarshal(buf []byte) error {
	if len(buf) != 16 {
		return errors.New("error size")
	}
	t.Uin = binary.BigEndian.Uint64(buf[:8])
	t.CreateTime = binary.BigEndian.Uint64(buf[8:])
	return nil
}

type AESTicketProvider struct {
	Key []byte
}

func NewAESTicketProvider() *AESTicketProvider {
	key := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, key)
	return &AESTicketProvider{Key: key}
}

func (p *AESTicketProvider) EncodeTicket(uin uint64) (string, error) {
	t := &Ticket{
		Uin:        uin,
		CreateTime: uint64(time.Now().Unix()),
	}
	if data, err := Encrypt(p.Key, t.Marshal()); err != nil {
		return "", err
	} else {
		return EncodeBase64WebString(base64.RawStdEncoding.EncodeToString(data)), nil
	}
}

func (p *AESTicketProvider) DecodeTicket(ticket string) (*Ticket, error) {
	rawData, err := base64.RawStdEncoding.DecodeString(DecodeBase64WebString(ticket))
	if err != nil {
		return nil, err
	}
	if data, err := Decrypt(p.Key, rawData); err != nil {
		return nil, err
	} else {
		t := &Ticket{}
		if err := t.Unmarshal(data); err != nil {
			return nil, err
		}
		return t, nil
	}
}

func Encrypt(key, data []byte) ([]byte, error) {
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

func Decrypt(key, data []byte) (plaintext []byte, err error) {
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
		return nil, errors.New("error data size")
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return
}

func EncodeBase64WebString(in string) string {
	in = strings.ReplaceAll(in, "/", "_")
	return in
}

func DecodeBase64WebString(in string) string {
	in = strings.ReplaceAll(in, "_", "/")
	return in
}
