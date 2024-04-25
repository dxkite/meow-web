package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"dxkite.cn/log"
	"dxkite.cn/meownest/src/application/router"
)

func init() {
	log.SetOutput(log.NewColorWriter(os.Stdout))
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func main() {
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

	r := router.New()
	r.Run(":2333")

	// app := application.New(application.WithRouter(router.New()))

	// if err := app.Serve(); err != nil {
	// 	log.Error(err)
	// }
}
