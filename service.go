package suda

import (
	"net/http"
	"os/exec"
	"path/filepath"

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
	router, err := srv.router.Build(&srv.Cfg.Auth)
	if err != nil {
		return err
	}

	listen := func(port Port) func() error {
		return func() error {
			l, err := Listen(port)
			if err != nil {
				return err
			}
			log.Info("listen", port.String())
			if err := http.Serve(l, router); err != nil {
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

func (srv *Service) registerRouters() {
	router := NewRouter()
	for _, route := range srv.Cfg.Routes {
		for _, uri := range route.Paths {
			log.Debug("register", srv.Cfg.Ports, uri)
			router.Add(uri, ForwardTarget{
				Name:      srv.Cfg.Name + ":" + route.Name,
				Auth:      route.Auth,
				Match:     route.Match,
				Rewrite:   route.Rewrite,
				Endpoints: route.EndPoints,
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
