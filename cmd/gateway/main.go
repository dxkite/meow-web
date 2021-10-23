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

	log.Println("load config success")
	r := route.NewRoute()
	r.Load(cfg.Routes)

	s := server.NewServer(cfg, r)
	if cfg.HotLoad > 0 {
		cfg.SetLoadTime(cfg.HotLoad)
		cfg.OnChange(func(c interface{}) {
			cfg.SetLoadTime(c.(*config.Config).HotLoad)
			r.Load(c.(*config.Config).Routes)
			s.ApplyHeaderFilter(c.(*config.Config).HttpAllowHeader)
			s.ApplyCorsConfig(c.(*config.Config).Cors)
		})
		cfg.HotLoadIfModify(p)
	}
	cfg.NotifyModify()

	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}
