package agent

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"dxkite.cn/meownest/pkg/crypto/token"
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

func (a *binaryAuth) validateToken(tokStr string) (*token.BinaryToken, error) {
	tok := &token.BinaryToken{}

	if err := tok.Decrypt(tokStr, token.NewAesCrypto([]byte(a.key))); err != nil {
		return nil, errors.New("invalid token")
	}

	if uint64(time.Now().Unix()) > tok.ExpireAt {
		return nil, errors.New("invalid token expire")
	}

	return tok, nil
}
