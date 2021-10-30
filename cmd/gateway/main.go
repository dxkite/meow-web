package main

import (
	"context"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/server"
	"dxkite.cn/log"
	"flag"
	"net"
	"os"
)

func init() {
	log.SetOutput(log.NewColorWriter())
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func main() {
	ctx, _ := context.WithCancel(context.Background())
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

	s := server.NewPortable(ctx, cfg)

	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error(err)
	}

	log.Println("gateway start at", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}
