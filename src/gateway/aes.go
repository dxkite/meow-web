package gateway

import (
	"encoding/base64"
	"net/http"

	"dxkite.cn/meownest/src/crypto"
)

type HttpAesHandler struct {
	Key    string
	Query  []string
	Header []string
	Cookie []string
}

func (s *HttpAesHandler) GetTokenFromRequest(req *http.Request) (token string, err error) {
	encryptedToken, err := s.getToken(req)
	if err != nil {
		return "", err
	}

	encryptedData, err := base64.RawURLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	if tok, err := crypto.AesDecrypt([]byte(s.Key), encryptedData); err != nil {
		return "", err
	} else {
		return string(tok), nil
	}
}

func (s *HttpAesHandler) getToken(req *http.Request) (token string, err error) {

	for _, name := range s.Cookie {
		if c, err := req.Cookie(name); err != http.ErrNoCookie {
			return "", err
		} else if c != nil {
			return c.Value, nil
		}
	}

	for _, name := range s.Query {
		if v := req.URL.Query().Get(name); v != "" {
			return v, nil
		}
	}

	for _, name := range s.Header {
		if v := req.Header.Get(name); v != "" {
			return v, nil
		}
	}

	return "", nil
}
