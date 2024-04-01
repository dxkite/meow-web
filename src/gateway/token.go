package gateway

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

type Token struct {
	Id              uint64
	ExpireAt        uint64
	RefreshExpireAt uint64
}

func (t Token) Marshal() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, t)
	return buf.Bytes()
}

func (t *Token) Unmarshal(buf []byte) error {
	return binary.Read(bytes.NewBuffer(buf), binary.BigEndian, t)
}

func (t *Token) EncodeToString() string {
	return base64.RawURLEncoding.EncodeToString(t.Marshal())
}

func (t *Token) DecodeString(val string) error {
	buf, err := base64.RawURLEncoding.DecodeString(val)
	if err != nil {
		return err
	}
	return t.Unmarshal(buf)
}
