package suda

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"dxkite.cn/log"
)

type Service struct {
	Cfg    *ServiceConfig
	router *Router
}

func (srv *Service) Config(cfg *ServiceConfig) error {
	srv.Cfg = cfg
	srv.registerRouters()
	return nil
}

func (srv *Service) Run() error {
	return srv.serve()
}

func (srv *Service) serve() error {
	listen := func(port Port) func() error {
		return func() error {
			l, err := Listen(port)
			if err != nil {
				return err
			}
			log.Info("listen", port.String())
			if err := http.Serve(l, http.HandlerFunc(srv.forward)); err != nil {
				return err
			}
			return nil
		}
	}

	execChain := ExecChain{}

	for _, port := range srv.Cfg.Ports {
		execChain = append(execChain, listen(port))
	}

	return execChain.Run()
}

func (srv *Service) forward(w http.ResponseWriter, req *http.Request) {
	uri := req.URL.Path
	log.Debug("forward", uri)
	_, routes := srv.router.MatchAll(uri)

	if len(routes) == 0 {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	var info *RouteInfo

	// 匹配路由
	if v, err := matchRouteTarget(req, routes); err != nil {
		log.Error("matchRouteTarget", err)
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else {
		info = v
	}

	// 清除请求头
	req.Header.Del(srv.Cfg.Auth.Header)

	if info.Auth {
		v := srv.getAuthToken(req)
		if v == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		req.Header.Set(srv.Cfg.Auth.Header, v.Value)
		log.Debug("auth header", srv.Cfg.Auth.Header, v.Value)
	}

	if len(info.Rewrite.Regex) >= 2 {
		if v, err := regexReplaceAll(info.Rewrite.Regex, uri, info.Rewrite.Replace); err != nil {
			log.Error("regexReplaceAll", err)
		} else {
			uri = v
		}
	}

	log.Debug("uri", strconv.Quote(req.URL.Path), strconv.Quote(uri))

	endpoint := matchEndpoint(req, info.EndPoints)

	if err := srv.forwardEndpoint(w, req, endpoint, uri); err != nil {
		return
	}

	return
}

func (srv *Service) getAuthToken(req *http.Request) *Token {
	return srv.getAuthTokenAes(req)
}

func (srv *Service) getAuthTokenAes(req *http.Request) *Token {
	if srv.Cfg.Auth.Type != "aes" {
		return nil
	}

	b := readAuthData(req, srv.Cfg.Auth.Source)
	log.Debug("read auth data", b)

	if b == "" {
		return nil
	}

	enc, err := base64.RawURLEncoding.DecodeString(b)
	if err != nil {
		log.Error("decode token error", err)
		return nil
	}

	data, err := AesDecrypt([]byte(srv.Cfg.Auth.Aes.Key), enc)
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

func (_ *Service) forwardEndpoint(w http.ResponseWriter, req *http.Request, endpoint *Endpoint, uri string) error {
	log.Debug("dial", endpoint, uri)
	rmt, err := Dial(endpoint.Port)
	if err != nil {
		log.Error("Dial", err)
		http.Error(w, "Unavailable Server", http.StatusInternalServerError)
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

func (srv *Service) registerRouters() {
	router := NewRouter()
	for _, route := range srv.Cfg.Routes {
		for _, uri := range route.Paths {
			log.Debug("register", srv.Cfg.Ports, uri)
			router.Add(uri, &RouteInfo{
				Name:        srv.Cfg.Name + ":" + route.Name,
				RouteConfig: &route,
			})
		}
	}
	srv.router = router
}

func execInstance(ins *InstanceConfig) error {
	ap, err := filepath.Abs(ins.Exec[0])
	if err != nil {
		log.Error("exec", ins.Exec, err)
		return err
	}

	bp := filepath.Dir(ap)
	ins.Exec[0] = ap

	w := MakeNameLoggerWriter(ins.Name)
	cmd := &exec.Cmd{
		Path:   ap,
		Dir:    bp,
		Args:   ins.Exec,
		Stderr: w,
		Stdout: w,
	}

	if err := cmd.Start(); err != nil {
		log.Error("exec", ins.Exec, err)
		return err
	}

	log.Info("exec", ins.Exec, "pid", cmd.Process.Pid)

	defer func() {
		cmd.Process.Kill()
	}()

	return cmd.Wait()
}
