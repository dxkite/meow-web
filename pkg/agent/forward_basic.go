package agent

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type BasicForwardHandler struct {
	network         string
	address         string
	timeout         time.Duration
	RewriteRequest  func(req *http.Request) error
	RewriteResponse func(resp *http.Response) error
}

func NewForward(network, address string, timeout time.Duration) *BasicForwardHandler {
	return &BasicForwardHandler{network: network, address: address, timeout: timeout}
}

func (h *BasicForwardHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	if err := h.rewriteRequest(req); err != nil {
		http.Error(w, "write request error: "+err.Error(), http.StatusBadGateway)
		return
	}

	rmt, err := net.DialTimeout(h.network, h.address, h.timeout)
	if err != nil {
		http.Error(w, "dial remote error: "+err.Error(), http.StatusBadGateway)
		return
	}

	defer rmt.Close()

	if err := req.WriteProxy(rmt); err != nil {
		http.Error(w, "write proxy error: "+err.Error(), http.StatusBadGateway)
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(rmt), req)
	if err != nil {
		http.Error(w, "read response error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.rewriteResponse(resp); err != nil {
		http.Error(w, "write response error: "+err.Error(), http.StatusBadGateway)
		return
	}

	// 是否升级到websocket
	isWebsocket := h.isUpgradeToWebsocket(req) && resp.StatusCode == http.StatusSwitchingProtocols

	// 普通http请求
	if !isWebsocket {
		h.copyHeader(w, resp.Header)
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			http.Error(w, "write response Error", http.StatusBadGateway)
			return
		}
		return
	}

	// 长链接 websocket 处理
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "attach hijack connection error", http.StatusBadGateway)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "hijack response Error", http.StatusBadGateway)
		return
	}

	defer conn.Close()

	if err := resp.Write(conn); err != nil {
		http.Error(w, "write response error", http.StatusBadGateway)
		return
	}

	if _, _, err = h.transport(conn, rmt); err != nil {
		http.Error(w, "transport error", http.StatusBadGateway)
		return
	}
}

func (h BasicForwardHandler) rewriteRequest(req *http.Request) error {
	if h.RewriteRequest != nil {
		return h.RewriteRequest(req)
	}
	return nil
}

func (h BasicForwardHandler) rewriteResponse(resp *http.Response) error {
	if h.RewriteResponse != nil {
		return h.RewriteResponse(resp)
	}
	return nil
}

func (h BasicForwardHandler) isUpgradeToWebsocket(req *http.Request) bool {
	connection := req.Header.Get("Connection")
	upgrade := req.Header.Get("Upgrade")
	if strings.ToLower(connection) == "upgrade" &&
		strings.ToLower(upgrade) == "websocket" {
		return true
	}
	return false
}

func (h BasicForwardHandler) copyHeader(w http.ResponseWriter, header http.Header) {
	for k, v := range header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
}

func (h BasicForwardHandler) transport(src, dst io.ReadWriter) (up, down int64, err error) {
	var closeCh = make(chan struct{})
	var errCh = make(chan error)

	go func() {
		// remote -> local
		var _err error
		if down, _err = io.Copy(src, dst); _err != nil && _err != io.EOF {
			errCh <- _err
			return
		}
		closeCh <- struct{}{}
	}()

	go func() {
		// local -> remote
		var _err error
		if up, _err = io.Copy(dst, src); _err != nil && _err != io.EOF {
			errCh <- _err
			return
		}
		closeCh <- struct{}{}
	}()

	select {
	case err = <-errCh:
		return
	case <-closeCh:
		return
	}
}
