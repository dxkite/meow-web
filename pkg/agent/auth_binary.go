package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net/http"
	"strconv"
	"time"

	"dxkite.cn/meownest/pkg/crypto"
)

func NewBinaryAuth(key, header string, source []*AuthorizeSource) AuthorizeHandler {
	return &binaryAuth{key: key, header: header, source: source}
}

type AuthorizeSource struct {
	Source string
	Name   string
}

type binaryAuth struct {
	key    string
	header string
	source []*AuthorizeSource
}

func (a *binaryAuth) HandleAuthorizeCheck(w http.ResponseWriter, req *http.Request) bool {
	req.Header.Del(a.header)

	if a.key == "" {
		return true
	}

	for _, v := range a.source {
		tok := VarFrom(req, v.Source, v.Name)
		if token, err := a.validateToken(tok); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return false
		} else {
			req.Header.Add(a.header, strconv.FormatUint(token.Id, 10))
			return true
		}
	}

	return false
}

func (a *binaryAuth) validateToken(tokStr string) (*BinaryToken, error) {
	encryptedData, err := base64.RawURLEncoding.DecodeString(tokStr)
	if err != nil {
		return nil, err
	}

	tok, err := crypto.AesDecrypt([]byte(a.key), encryptedData)
	if err != nil {
		return nil, err
	}

	token := &BinaryToken{}
	if err := token.Unmarshal([]byte(tok)); err != nil {
		return nil, errors.New("invalid token")
	}

	if uint64(time.Now().Unix()) > token.ExpireAt {
		return nil, errors.New("invalid token expire")
	}

	return token, nil
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
