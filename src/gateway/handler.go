package gateway

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dxkite.cn/log"
	"dxkite.cn/meownest/src/utils"
)

type AuthorizationHandler interface {
	GetTokenFromRequest(req *http.Request) (token string, err error)
}

type RewriteConfig struct {
	Regex   string
	Replace string
}

type HttpForwardHandler struct {
	Name           string
	Rewrite        *RewriteConfig
	AuthCheck      bool
	AuthHandler    AuthorizationHandler
	IdAssignHeader string
	Endpoints      []string
}

func (h HttpForwardHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.IdAssignHeader != "" {
		req.Header.Del(h.IdAssignHeader)
	}

	var token *Token
	if h.AuthCheck {
		if t, err := h.checkAuth(req); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			token = t
		}
	}

	endpoint := h.Endpoints[utils.Intn(len(h.Endpoints))]
	h.forwardEndpoint(w, req, endpoint, token)
}

func (h *HttpForwardHandler) forwardEndpoint(w http.ResponseWriter, req *http.Request, endpoint string, token *Token) error {
	requestId := utils.GenerateRequestId()
	h.rewriteRequest(req, requestId, endpoint, token)

	log.Debug("dial", h.Name, endpoint, requestId, req.URL)

	uri, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	timeout := 500 * time.Millisecond
	if t := uri.Query().Get("timeout"); t != "" {
		tm, _ := strconv.ParseUint(t, 10, 64)
		timeout = time.Duration(tm) * time.Millisecond
	}

	rmt, err := utils.DialTimeout(uri, timeout)
	if err != nil {
		log.Error("Dial", err)
		http.Error(w, "Unavailable Service", http.StatusInternalServerError)
		return err
	}

	defer rmt.Close()

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

	// 重写响应
	h.rewriteResponse(resp, requestId)

	// 是否升级到websocket
	isWebsocket := h.isUpgradeToWebsocket(req) && resp.StatusCode == http.StatusSwitchingProtocols

	// 普通http请求
	if !isWebsocket {
		h.copyHeader(w, resp.Header)
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

	rb, wb, err := utils.Transport(conn, rmt)
	if err != nil {
		log.Error("transport", err)
		http.Error(w, "Transport Error", http.StatusInternalServerError)
		return err
	}

	log.Debug("transport", rb, wb)
	return nil
}

var ErrUnauthorized = errors.New("Unauthorized")

func (h HttpForwardHandler) checkAuth(req *http.Request) (*Token, error) {
	tok, err := h.AuthHandler.GetTokenFromRequest(req)
	if err != nil {
		return nil, ErrUnauthorized
	}
	token := &Token{}
	if err := token.Unmarshal([]byte(tok)); err != nil {
		return nil, ErrUnauthorized
	}

	if uint64(time.Now().Unix()) > token.ExpireAt {
		return nil, ErrUnauthorized
	}
	return token, nil
}

func (h HttpForwardHandler) rewriteRequest(req *http.Request, requestId, endpoint string, token *Token) error {
	uri := req.URL.Path
	if h.Rewrite != nil {
		if h.Rewrite.Regex != "" {
			if v, err := utils.StringReplaceAll(h.Rewrite.Regex, uri, h.Rewrite.Replace); err != nil {
				log.Error("StringReplaceAll", err)
			} else {
				log.Debug("rewrite url", uri, "->", v, h.Rewrite.Regex)
				uri = v
			}
		}
	}
	req.URL.Path = uri
	req.RequestURI = req.URL.String()
	req.Header.Set("Request-Id", requestId)
	req.Header.Set("X-Forward-Endpoint", endpoint)
	if token != nil {
		req.Header.Set(h.IdAssignHeader, strconv.FormatUint(token.Id, 10))
	}
	return nil
}

func (h HttpForwardHandler) rewriteResponse(resp *http.Response, requestId string) {
	resp.Header.Set("Request-Id", requestId)
	resp.Header.Set("X-Powered-By", "meow")
	resp.Header.Del(h.IdAssignHeader)
}

func (h HttpForwardHandler) isUpgradeToWebsocket(req *http.Request) bool {
	connection := req.Header.Get("Connection")
	upgrade := req.Header.Get("Upgrade")
	if strings.ToLower(connection) == "upgrade" &&
		strings.ToLower(upgrade) == "websocket" {
		return true
	}
	return false
}

func (h HttpForwardHandler) copyHeader(w http.ResponseWriter, header http.Header) {
	for k, v := range header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
}
