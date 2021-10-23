package main

import (
	"dxkite.cn/gateway/config"
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
	cfg := &config.Config{}
	if err := cfg.LoadFromFile("./conf/config.yml"); err != nil {
		log.Error(err)
	}
	log.Println("load config success")
	s := server.NewServer(cfg)
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}
