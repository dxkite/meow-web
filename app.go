package suda

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"dxkite.cn/log"
)

type App struct {
	Cfg    *Config
	mods   []*ModuleConfig
	router *Router
}

func (app *App) Config(name string) error {
	app.Cfg = &Config{}
	if err := loadYaml(name, app.Cfg); err != nil {
		return err
	}

	// 加载模块
	if err := app.loadModules(app.Cfg.ModuleConfig); err != nil {
		return err
	}

	app.registerModules()

	return nil
}

func (app *App) Run() error {
	go app.execModules()
	return app.web()
}

func (app *App) internal() error {
	return nil
}

func (app *App) web() error {
	l, err := net.Listen("tcp", app.Cfg.Addr)
	if err != nil {
		log.Debug("Listen", err)
		return err
	}

	return http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.forward(w, r)
	}))
}

func (app *App) forward(w http.ResponseWriter, req *http.Request) {
	uri := req.URL.Path
	log.Debug("forward", uri)
	_, route := app.router.Match(uri)

	if route == nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	info, ok := route.(*RouteInfo)
	if !ok {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if info.Auth {
		v := app.getAuthToken(req)
		if v == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !matchScope(req.URL.Path, v.Scope) {
			http.Error(w, "unauthorized scope", http.StatusUnauthorized)
			return
		}

		req.Header.Set(app.Cfg.Auth.Header, v.Value)
		log.Debug("auth header", app.Cfg.Auth.Header, v.Value)
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

	if err := app.forwardEndpoint(w, req, endpoint, uri); err != nil {
		return
	}

	return
}

func (app *App) getAuthToken(req *http.Request) *Token {
	return app.getAuthTokenAes(req)
}

func (app *App) getAuthTokenAes(req *http.Request) *Token {
	if app.Cfg.Auth.Type != "aes" {
		return nil
	}

	b := readAuthData(req, app.Cfg.Auth.Source)
	enc, err := base64.RawURLEncoding.DecodeString(b)
	if err != nil {
		return nil
	}

	data, err := AesDecrypt([]byte(app.Cfg.Auth.Aes.Key), enc)
	if err != nil {
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

func (_ *App) forwardEndpoint(w http.ResponseWriter, req *http.Request, endpoint, uri string) error {
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

func (app *App) loadModules(p string) error {
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
		if err := app.loadModuleConfig(p, name); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) loadModuleConfig(p, name string) error {
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

	app.mods = append(app.mods, cfg)
	return nil
}

func (app *App) registerModules() {
	router := NewRouter()
	for _, mod := range app.mods {
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
	app.router = router
}

func (app *App) execModules() {
	for _, mod := range app.mods {
		if len(mod.Exec) > 0 {
			go func(mod *ModuleConfig) {
				err := app.execModule(mod)
				if err != nil {
					log.Error("execModule", err)
				}
			}(mod)
		}
	}
}

func (app *App) execModule(cfg *ModuleConfig) error {
	ap, err := filepath.Abs(cfg.Exec[0])
	if err != nil {
		log.Error("exec", cfg.Exec, err)
		return err
	}

	bp := filepath.Dir(ap)
	cfg.Exec[0] = ap

	w := MakeNameLoggerWriter(cfg.Name)
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
