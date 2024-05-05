package token

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

type Crypto interface {
	Encrypt(data []byte) (encryptData []byte, err error)
	Decrypt(encryptData []byte) (data []byte, err error)
}

type BinaryToken struct {
	Id       uint64
	ExpireAt uint64
}

func (t BinaryToken) Marshal() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, t)
	return buf.Bytes()
}

func (t *BinaryToken) Unmarshal(buf []byte) error {
	return binary.Read(bytes.NewBuffer(buf), binary.BigEndian, t)
}

func (t *BinaryToken) EncodeToString() string {
	return base64.RawURLEncoding.EncodeToString(t.Marshal())
}

func (t *BinaryToken) DecodeString(val string) error {
	buf, err := base64.RawURLEncoding.DecodeString(val)
	if err != nil {
		return err
	}
	return t.Unmarshal(buf)
}

func (t *BinaryToken) Encrypt(c Crypto) (string, error) {
	tok := t.Marshal()
	tokEnc, err := c.Encrypt(tok)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tokEnc), nil
}

func (t *BinaryToken) Decrypt(token string, c Crypto) error {
	buf, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return err
	}
	tokDec, err := c.Decrypt(buf)
	if err != nil {
		return err
	}
	return t.Unmarshal(tokDec)
}
