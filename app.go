package suda

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path"
	"regexp"
)

type App struct {
	cfg    *Config
	mods   []*ModuleConfig
	router *Router
}

func (app *App) Config(name string) error {
	app.cfg = &Config{}
	if err := loadYaml(name, app.cfg); err != nil {
		return err
	}

	// 加载模块
	if err := app.loadModules(app.cfg.ModuleConfig); err != nil {
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
	l, err := net.Listen("tcp", app.cfg.Addr)
	if err != nil {
		fmt.Println("Listen", err)
		return err
	}
	fmt.Println("run at", app.cfg.Addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept", err)
			continue
		}
		go func(c net.Conn) {
			if err := app.handleConn(c); err != nil {
				fmt.Println("handleConn", err)
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
	fmt.Println("forward", uri)
	_, route := app.router.Match(uri)

	if route == nil {
		return writeBody(conn, http.StatusNotFound, "404 not found")
	}

	info, ok := route.(*RouteInfo)
	if !ok {
		return writeBody(conn, http.StatusNotFound, "404 not found")
	}

	endpoint := info.EndPoints[intn(len(info.EndPoints))]

	if len(info.Rewrite.Regex) >= 2 {
		reg, _ := regexp.Compile(info.Rewrite.Regex)
		if reg != nil {
			uri = reg.ReplaceAllString(uri, info.Rewrite.Replace)
		}
	}

	fmt.Println("uri", req.URL.Path, uri)

	if err := app.forwardEndpoint(req, conn, endpoint, uri); err != nil {
		return err
	}

	return nil
}

func (_ *App) forwardEndpoint(req *http.Request, conn net.Conn, endpoint, uri string) error {
	rmt, err := dial(endpoint)
	if err != nil {
		fmt.Println("Dial", err)
		return err
	}

	reqId := genRequestId()

	req.URL.Path = uri
	req.RequestURI = req.URL.String()
	req.Header.Set("X-Forward-Endpoint", endpoint)
	req.Header.Set("Request-Id", reqId)

	// write to remote
	if err := req.WriteProxy(rmt); err != nil {
		fmt.Println("WriteProxy", err)
		return err
	}

	resp, err := http.ReadResponse(bufio.NewReader(rmt), req)
	if err != nil {
		fmt.Println("http.ReadResponse", err)
		return err
	}

	resp.Header.Set("Request-Id", reqId)
	resp.Header.Set("X-Powered-By", "suda/1.0")

	if err := resp.Write(conn); err != nil {
		fmt.Println("resp.Write", err)
		return err
	}

	r, w, err := transport(conn, rmt)
	if err != nil {
		fmt.Println("transport", err)
		return err
	}

	fmt.Println("transport", r, w)
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
		fmt.Println("load", p, name)
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
				fmt.Println("register", uri)
				router.Add(uri, &RouteInfo{
					Name:      mod.Name + ":" + route.Name,
					Auth:      route.Auth,
					Rewrite:   route.Rewrite,
					EndPoints: mod.EndPoints,
				})
			}
		}
	}
	fmt.Println("registerModules", router)
	app.router = router
}

func (app *App) execModules() {
	for _, mod := range app.mods {
		if len(mod.Exec) > 0 {
			go func(mod *ModuleConfig) {
				err := app.execModule(mod)
				if err != nil {
					fmt.Println(err)
				}
			}(mod)
		}
	}
}

func (app *App) execModule(cfg *ModuleConfig) error {
	cmd := exec.Command(cfg.Exec[0], cfg.Exec[1:]...)
	return cmd.Run()
}
