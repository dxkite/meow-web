package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"dxkite.cn/log"
	"dxkite.cn/meownest/src/bootstrap"
)

func init() {
	log.SetOutput(log.NewColorWriter(os.Stdout))
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func main() {
	ctx, exit := context.WithCancel(context.Background())
	defer exit()
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			log.Error("[panic error]", r)
			log.Error(string(buf[:n]))
			name := fmt.Sprintf("crash-%s.log", time.Now().Format("20060102150405"))
			panicErr := string(buf[:n])
			_ = os.WriteFile(name, []byte(panicErr), os.ModePerm)
		}
	}()

	cfg := "./config.yaml"
	if len(os.Args) >= 2 {
		cfg = os.Args[1]
	}

	if err := bootstrap.ServerGateway(ctx, cfg); err != nil {
		log.Error(err)
	}
}
