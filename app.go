package suda

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"os/exec"
	"path"
	"regexp"
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
	log.Debug("run at", app.Cfg.Addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Debug("Accept", err)
			continue
		}
		go func(c net.Conn) {
			if err := app.handleConn(c); err != nil {
				log.Debug("handleConn", err)
			}
		}(conn)
	}
}

func (app *App) handleConn(conn net.Conn) error {
	defer conn.Close()

	r := bufio.NewReader(conn)
	req, err := http.ReadRequest(r)
	if err != nil {
		return err
	}

	if err := app.forward(req, conn); err != nil {
		return err
	}

	return nil
}

func (app *App) forward(req *http.Request, conn net.Conn) error {
	uri := req.URL.Path
	log.Debug("forward", uri)
	_, route := app.router.Match(uri)

	if route == nil {
		return writeBody(conn, http.StatusNotFound, "404 not found")
	}

	info, ok := route.(*RouteInfo)
	if !ok {
		return writeBody(conn, http.StatusNotFound, "404 not found")
	}

	if info.Auth {
		v := app.getAuthToken(req)
		if v == nil {
			return writeBody(conn, http.StatusUnauthorized, "unauthorized")
		}

		if !matchScope(req.URL.Path, v.Scope) {
			return writeBody(conn, http.StatusUnauthorized, "unauthorized scope")
		}

		req.Header.Set(app.Cfg.Auth.Header, v.Value)
		log.Debug("auth header", app.Cfg.Auth.Header, v.Value)
	}

	endpoint := info.EndPoints[intn(len(info.EndPoints))]

	if len(info.Rewrite.Regex) >= 2 {
		reg, _ := regexp.Compile(info.Rewrite.Regex)
		if reg != nil {
			uri = reg.ReplaceAllString(uri, info.Rewrite.Replace)
		}
	}

	log.Debug("uri", req.URL.Path, uri)

	if err := app.forwardEndpoint(req, conn, endpoint, uri); err != nil {
		return err
	}

	return nil
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

func (_ *App) forwardEndpoint(req *http.Request, conn net.Conn, endpoint, uri string) error {
	rmt, err := dial(endpoint)
	if err != nil {
		log.Debug("Dial", err)
		return err
	}

	reqId := genRequestId()

	req.URL.Path = uri
	req.RequestURI = req.URL.String()
	req.Header.Set("X-Forward-Endpoint", endpoint)
	req.Header.Set("Request-Id", reqId)

	// write to remote
	if err := req.WriteProxy(rmt); err != nil {
		log.Debug("WriteProxy", err)
		return err
	}

	resp, err := http.ReadResponse(bufio.NewReader(rmt), req)
	if err != nil {
		log.Debug("http.ReadResponse", err)
		return err
	}

	resp.Header.Set("Request-Id", reqId)
	resp.Header.Set("X-Powered-By", "suda/1.0")

	if err := resp.Write(conn); err != nil {
		log.Debug("resp.Write", err)
		return err
	}

	r, w, err := transport(conn, rmt)
	if err != nil {
		log.Debug("transport", err)
		return err
	}

	log.Debug("transport", r, w)
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
		cfg := &ModuleConfig{}
		if err := loadYaml(path.Join(p, name), cfg); err != nil {
			return err
		}

		app.mods = append(app.mods, cfg)
	}

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
					log.Debug(err)
				}
			}(mod)
		}
	}
}

func (app *App) execModule(cfg *ModuleConfig) error {
	cmd := exec.Command(cfg.Exec[0], cfg.Exec[1:]...)
	w := makeLoggerWriter(cfg.Name)
	cmd.Stderr = w
	cmd.Stdout = w
	return cmd.Run()
}
