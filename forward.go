package meownest

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"dxkite.cn/log"
)

type ForwardTarget struct {
	Name       string
	Auth       bool
	AuthConfig *AuthConfig
	Match      []RouteMatch
	Rewrite    RewriteConfig
	Endpoints  []Endpoint
}

func (target ForwardTarget) MatchRequest(req *http.Request) bool {
	return matchRequest(req, target.Match)
}

func (target ForwardTarget) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	uri := req.URL.Path

	// 清除请求头
	req.Header.Del(target.AuthConfig.Header)

	if target.Auth && target.AuthConfig != nil {
		v := target.getAuthToken(req)
		if v == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		req.Header.Set(target.AuthConfig.Header, v.Value)
		log.Debug("auth header", target.AuthConfig.Header, v.Value)
	}

	if target.Rewrite.Regex != "" {
		if v, err := regexReplaceAll(target.Rewrite.Regex, uri, target.Rewrite.Replace); err != nil {
			log.Error("regexReplaceAll", err)
		} else {
			log.Debug("rewrite url", uri, "->", v, target.Rewrite.Regex)
			uri = v
		}
	}

	log.Debug("match", target.Name, "uri", strconv.Quote(req.URL.Path), strconv.Quote(uri))

	endpoint := matchEndpoint(req, target.Endpoints)

	if err := target.forwardEndpoint(w, req, endpoint, uri); err != nil {
		return
	}

	return
}

func (_ *ForwardTarget) forwardEndpoint(w http.ResponseWriter, req *http.Request, endpoint *Endpoint, uri string) error {
	log.Debug("dial", endpoint, uri)

	timeout := 500 * time.Millisecond
	if endpoint.Timeout != 0 {
		timeout = time.Duration(endpoint.Timeout) * time.Millisecond
	}

	rmt, err := DialTimeout(endpoint.Port, timeout)
	if err != nil {
		log.Error("Dial", err)
		http.Error(w, "Unavailable Service", http.StatusInternalServerError)
		return err
	}

	defer rmt.Close()

	reqId := genRequestId()

	req.URL.Path = uri
	req.RequestURI = req.URL.String()
	req.Header.Set("X-Forward-Endpoint", endpoint.String())
	req.Header.Set("Request-Id", reqId)

	// write to remote
	if err := req.WriteProxy(rmt); err != nil {
		log.Error("WriteProxy", err)
		http.Error(w, "Write Proxy Error", http.StatusInternalServerError)
		return err
	}

	resp, err := http.ReadResponse(bufio.NewReader(rmt), req)
	if err != nil {
		log.Error("http.ReadResponse", err)
		http.Error(w, "Read Response Error", http.StatusInternalServerError)
		return err
	}

	resp.Header.Set("Request-Id", reqId)
	resp.Header.Set("X-Powered-By", "suda")

	// 是否升级到websocket
	isWebsocket := isUpgradeToWebsocket(req) && resp.StatusCode == http.StatusSwitchingProtocols

	// 普通http请求
	if !isWebsocket {
		copyHeader(w, resp.Header)
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Error("copy error", err)
			http.Error(w, "Write Response Error", http.StatusInternalServerError)
			return err
		}
		return nil
	}

	log.Info("handle websocket")

	// websocket 请求
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Attach Hijack Connection Error", http.StatusInternalServerError)
		return errors.New("error to attach http.Hijacker")
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		log.Error("Hijack error", err)
		http.Error(w, "Hijack Response Error", http.StatusInternalServerError)
		return err
	}
	defer conn.Close()

	if err := resp.Write(conn); err != nil {
		log.Error("resp.Write", err)
		http.Error(w, "Write Response Error", http.StatusInternalServerError)
		return err
	}

	rb, wb, err := Transport(conn, rmt)
	if err != nil {
		log.Error("transport", err)
		http.Error(w, "Transport Error", http.StatusInternalServerError)
		return err
	}

	log.Debug("transport", rb, wb)
	return nil
}

func (f *ForwardTarget) getAuthToken(req *http.Request) *Token {
	return f.getAuthTokenAes(req)
}

func (f *ForwardTarget) getAuthTokenAes(req *http.Request) *Token {
	if f.AuthConfig.Type != "aes" {
		return nil
	}

	b := readAuthData(req, f.AuthConfig.Source)
	log.Debug("read auth data", b)

	if b == "" {
		return nil
	}

	enc, err := base64.RawURLEncoding.DecodeString(b)
	if err != nil {
		log.Error("decode token error", err)
		return nil
	}

	data, err := AesDecrypt([]byte(f.AuthConfig.Aes.Key), enc)
	if err != nil {
		log.Error("decrypt token error", err)
		return nil
	}

	token := &Token{}
	if err := json.Unmarshal([]byte(data), token); err != nil {
		return nil
	}

	if time.Now().Unix() > token.ExpireAt {
		return nil
	}

	return token
}
