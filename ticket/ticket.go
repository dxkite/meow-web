package ticket

import (
	"encoding/binary"
	"errors"
	"strings"
)

type TicketType int

const (
	TicketRandomAes = iota
	TicketAes
	TicketRsaSign
)

type TicketEnDecoder interface {
	Encode(session *SessionData) (string, error)
	Decode(ticket string) (*SessionData, error)
}

type SessionData struct {
	Uin        uint64
	CreateTime uint64
}

func (t *SessionData) Marshal() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf, t.Uin)
	binary.BigEndian.PutUint64(buf[8:], t.CreateTime)
	return buf
}

func (t *SessionData) Unmarshal(buf []byte) error {
	if len(buf) != 16 {
		return errors.New("error size")
	}
	t.Uin = binary.BigEndian.Uint64(buf[:8])
	t.CreateTime = binary.BigEndian.Uint64(buf[8:])
	return nil
}

func EncodeBase64WebString(in string) string {
	in = strings.ReplaceAll(in, "/", "_")
	return in
}

func DecodeBase64WebString(in string) string {
	in = strings.ReplaceAll(in, "_", "/")
	return in
}
