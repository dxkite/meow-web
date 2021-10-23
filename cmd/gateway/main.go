package main

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/gateway/server"
	"dxkite.cn/log"
	"flag"
	"net"
	"os"
)

func main() {
	conf := flag.String("conf", "./config.yml", "the config file")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	cfg := config.NewConfig()
	if err := cfg.LoadFromFile(*conf); err != nil {
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
			r.ClearAll()
			r.Load(c.(*config.Config).Routes)
			s.ApplyHeaderFilter(c.(*config.Config).HttpAllowHeader)
			s.ApplyCorsConfig(c.(*config.Config).Cors)
		})
		cfg.HotLoadIfModify(*conf)
	}
	cfg.NotifyModify()

	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error(err)
	}
	log.Println("server start at", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}
