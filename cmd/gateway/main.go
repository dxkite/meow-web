package main

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/gateway/server"
	"dxkite.cn/log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":2333")
	if err != nil {
		log.Error(err)
	}
	log.Println("server start at", l.Addr())
	cfg := config.NewConfig()
	p := "./conf/config.yml"
	if err := cfg.LoadFromFile("./conf/config.yml"); err != nil {
		log.Error(err)
	}
	if cfg.HotLoad > 0 {
		cfg.SetLoadTime(cfg.HotLoad)
		cfg.OnChange(func(c interface{}) {
			cfg.SetLoadTime(c.(*config.Config).HotLoad)
		})
		cfg.HotLoadIfModify(p)
	}
	log.Println("load config success")
	r := route.NewRoute()
	r.Load(cfg.Routes)
	cfg.OnChange(func(cfg interface{}) {
		r.Load(cfg.(*config.Config).Routes)
	})
	s := server.NewServer(cfg, r)
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}
