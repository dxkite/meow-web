package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"dxkite.cn/log"
	"dxkite.cn/suda"
)

func init() {
	log.SetOutput(log.NewColorWriter(os.Stdout))
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func applyLogConfig(ctx context.Context, cfg *suda.Config) {
	if cfg.LogLevel != 0 {
		log.SetLevel(log.LogLevel(cfg.LogLevel))
	}

	if cfg.LogFile == "" {
		return
	}

	log.Println("log output file", cfg.LogFile)
	filename := cfg.LogFile
	var w io.Writer
	if f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err != nil {
		log.Warn("log file open error", filename)
		return
	} else {
		w = f
		if filepath.Ext(filename) == ".json" {
			w = log.NewJsonWriter(w)
		} else {
			w = log.NewTextWriter(w)
		}
		go func() {
			<-ctx.Done()
			_ = f.Close()
		}()
	}
	log.SetOutput(log.MultiWriter(w, log.Writer()))
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

	app := new(suda.App)
	if err := app.Config("./config.yaml"); err != nil {
		log.Error(err)
		return
	}

	applyLogConfig(ctx, app.Cfg)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
