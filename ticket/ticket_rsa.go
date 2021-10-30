package ticket

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"dxkite.cn/gateway/util"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"os"
	"time"
)

type RsaTicket struct {
	pri *rsa.PrivateKey
	pub *rsa.PublicKey
}

type signPackage struct {
	data []byte
	sign []byte
}

const sizeLen = 4

func (t *signPackage) Marshal() []byte {
	head := make([]byte, sizeLen)
	binary.BigEndian.PutUint32(head, uint32(len(t.data)))
	buf := &bytes.Buffer{}
	buf.Write(head)
	buf.Write(t.data)
	buf.Write(t.sign)
	return buf.Bytes()
}

func (t *signPackage) Unmarshal(buf []byte) error {
	if len(buf) <= sizeLen {
		return errors.New("error size")
	}
	size := binary.BigEndian.Uint32(buf)
	t.data = buf[sizeLen : size+sizeLen]
	t.sign = buf[size+sizeLen:]
	return nil
}

func (r *RsaTicket) Encode(session *SessionData) (string, error) {
	if session.CreateTime == 0 {
		session.CreateTime = uint64(time.Now().Unix())
	}
	data := session.Marshal()
	if sign, err := r.sign(data); err != nil {
		return "", err
	} else {
		p := &signPackage{
			data: data,
			sign: sign,
		}
		return base64.RawStdEncoding.EncodeToString(p.Marshal()), nil
	}
}

func (r *RsaTicket) Decode(ticket string) (*SessionData, error) {
	rawData, err := base64.RawStdEncoding.DecodeString(ticket)
	if err != nil {
		return nil, err
	}
	p := &signPackage{}
	if err := p.Unmarshal(rawData); err != nil {
		return nil, err
	}
	if err := r.check(p.data, p.sign); err != nil {
		return nil, err
	} else {
		t := &SessionData{}
		if err := t.Unmarshal(p.data); err != nil {
			return nil, err
		}
		return t, nil
	}
}

func (r *RsaTicket) sign(data []byte) ([]byte, error) {
	if r.pri == nil {
		return nil, errors.New("empty rsa private key")
	}
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.pri, crypto.SHA256, d)
}

func (r *RsaTicket) check(message []byte, sign []byte) error {
	if r.pub == nil {
		return errors.New("empty rsa public key")
	}
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.pub, crypto.SHA256, d, sign)
}

func NewRsaTicket(pri, cert string) TicketEnDecoder {
	certPEMBlock, _ := os.ReadFile(pri)
	keyPEMBlock, _ := os.ReadFile(cert)
	priKey, _ := util.ParsePrivateKey(certPEMBlock)
	pubKey, _ := util.ParsePublicKeyFromCertificate(keyPEMBlock)
	return &RsaTicket{
		pri: priKey,
		pub: pubKey,
	}
}
