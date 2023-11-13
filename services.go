package suda

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"dxkite.cn/log"
)

type Service struct {
	Cfg    *ServiceConfig
	mods   []*ModuleConfig
	router *Router
}

func (srv *Service) Config(cfg *ServiceConfig) error {
	srv.Cfg = cfg

	// 加载模块
	if err := srv.loadModules(srv.Cfg.ModuleConfig); err != nil {
		return err
	}

	srv.registerModules()
	return nil
}

func (srv *Service) Run() error {
	go srv.execModules()
	return srv.web()
}

func (srv *Service) web() error {
	log.Info("listen", srv.Cfg.Addr)
	return ListenAndServe(srv.Cfg.Addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.forward(w, r)
	}))
}

func (srv *Service) forward(w http.ResponseWriter, req *http.Request) {
	uri := req.URL.Path
	log.Debug("forward", uri)
	_, route := srv.router.Match(uri)

	if route == nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	info, ok := route.(*RouteInfo)
	if !ok {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	// 清除请求头
	req.Header.Del(srv.Cfg.Auth.Header)

	if info.Auth {
		v := srv.getAuthToken(req)
		if v == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		req.Header.Set(srv.Cfg.Auth.Header, v.Value)
		log.Debug("auth header", srv.Cfg.Auth.Header, v.Value)
	}

	endpoint := info.EndPoints[intn(len(info.EndPoints))]

	if len(info.Rewrite.Regex) >= 2 {
		if v, err := regexReplaceAll(info.Rewrite.Regex, uri, info.Rewrite.Replace); err != nil {
			log.Error("regexReplaceAll", err)
		} else {
			uri = v
		}
	}

	log.Debug("uri", req.URL.Path, uri)

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

func (_ *Service) forwardEndpoint(w http.ResponseWriter, req *http.Request, endpoint, uri string) error {
	log.Debug("dial", endpoint, uri)
	rmt, err := dial(endpoint)
	if err != nil {
		log.Error("Dial", err)
		return err
	}
	defer rmt.Close()

	reqId := genRequestId()

	req.URL.Path = uri
	req.RequestURI = req.URL.String()
	req.Header.Set("X-Forward-Endpoint", endpoint)
	req.Header.Set("Request-Id", reqId)

	// write to remote
	if err := req.WriteProxy(rmt); err != nil {
		log.Error("WriteProxy", err)
		return err
	}

	resp, err := http.ReadResponse(bufio.NewReader(rmt), req)
	if err != nil {
		log.Error("http.ReadResponse", err)
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
			return err
		}
		return nil
	}

	log.Info("handle websocket")

	// websocket 请求
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return errors.New("error to attach http.Hijacker")
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		log.Error("Hijack error", err)
		return err
	}
	defer conn.Close()

	if err := resp.Write(conn); err != nil {
		log.Error("resp.Write", err)
		return err
	}

	rb, wb, err := transport(conn, rmt)
	if err != nil {
		log.Error("transport", err)
		return err
	}

	log.Debug("transport", rb, wb)
	return nil
}

func (srv *Service) loadModules(p string) error {
	names, err := readDirNames(p)
	if err != nil {
		return err
	}

	for _, name := range names {
		ext := path.Ext(name)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		log.Debug("load", p, name)
		if err := srv.loadModuleConfig(p, name); err != nil {
			return err
		}
	}

	return nil
}

func (srv *Service) loadModuleConfig(p, name string) error {
	cfg := &ModuleConfig{}
	if err := loadYaml(path.Join(p, name), cfg); err != nil {
		return err
	}

	if len(cfg.Exec) > 0 {
		if !filepath.IsAbs(cfg.Exec[0]) {
			cfg.Exec[0] = filepath.Join(p, cfg.Exec[0])
		}
	}

	for i, ep := range cfg.EndPoints {
		if strings.HasPrefix(ep, "unix://") {
			sock := ep[7:]
			if !filepath.IsAbs(sock) {
				sock = filepath.Join(p, sock)
			}
			cfg.EndPoints[i] = "unix://" + sock
		}
	}

	srv.mods = append(srv.mods, cfg)
	return nil
}

func (srv *Service) registerModules() {
	router := NewRouter()
	for _, mod := range srv.mods {
		for _, route := range mod.Routes {
			for _, uri := range route.Paths {
				log.Debug("register", uri)
				router.Add(uri, &RouteInfo{
					Name:      mod.Name + ":" + route.Name,
					Auth:      route.Auth,
					Rewrite:   route.Rewrite,
					EndPoints: mod.EndPoints,
				})
			}
		}
	}
	log.Debug("registerModules", router)
	srv.router = router
}

func (srv *Service) execModules() {
	for _, mod := range srv.mods {
		if len(mod.Exec) > 0 {
			go func(mod *ModuleConfig) {
				err := srv.execModule(mod)
				if err != nil {
					log.Error("execModule", err)
				}
			}(mod)
		}
	}
}

func (srv *Service) execModule(cfg *ModuleConfig) error {
	ap, err := filepath.Abs(cfg.Exec[0])
	if err != nil {
		log.Error("exec", cfg.Exec, err)
		return err
	}

	bp := filepath.Dir(ap)
	cfg.Exec[0] = ap

	w := MakeNameLoggerWriter(srv.Cfg.Name + ":" + cfg.Name)
	cmd := &exec.Cmd{
		Path:   ap,
		Dir:    bp,
		Args:   cfg.Exec,
		Stderr: w,
		Stdout: w,
	}

	if err := cmd.Start(); err != nil {
		log.Error("exec", cfg.Exec, err)
		return err
	}

	log.Info("exec", cfg.Exec, "pid", cmd.Process.Pid)

	defer func() {
		cmd.Process.Kill()
	}()

	return cmd.Wait()
}
