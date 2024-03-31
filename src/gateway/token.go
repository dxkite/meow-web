package gateway

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
)

type Token struct {
	ExpireAt        uint64
	RefreshExpireAt uint64
	Id              uint64
}

func (t Token) Marshal() []byte {
	buf := make([]byte, 24)
	binary.BigEndian.PutUint64(buf, t.Id)
	binary.BigEndian.PutUint64(buf[8:16], t.ExpireAt)
	binary.BigEndian.PutUint64(buf[16:], t.RefreshExpireAt)
	return buf
}

func (t *Token) Unmarshal(buf []byte) error {
	if len(buf) != 24 {
		return errors.New("error size")
	}
	t.Id = binary.BigEndian.Uint64(buf[:8])
	t.ExpireAt = binary.BigEndian.Uint64(buf[8:16])
	t.RefreshExpireAt = binary.BigEndian.Uint64(buf[16:])
	return nil
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
