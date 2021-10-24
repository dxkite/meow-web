package util

import (
	"context"
	"dxkite.cn/log"
	"os"
)

func SetLogConfig(ctx context.Context, level log.LogLevel, filename string) {
	log.SetLevel(level)
	if len(filename) == 0 {
		return
	}
	if f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err != nil {
		log.Warn("log file open error", filename)
		return
	} else {
		log.SetOutput(log.MultiWriter(log.NewTextWriter(f), log.Writer()))
		go func() {
			<-ctx.Done()
			_ = f.Close()
		}()
	}
}
