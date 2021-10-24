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
			cc := c.(*config.Config)
			cfg.SetLoadTime(cc.HotLoad)
			r.ClearAll()
			r.Load(cc.Routes)
			s.ApplyHeaderFilter(cc.HttpAllowHeader)
			s.ApplyCorsConfig(cc.Cors)
		})
		cfg.HotLoadIfModify(*conf)
	}

	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error(err)
	}
	log.Println("server start at", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}
