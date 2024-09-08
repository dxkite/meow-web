package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"dxkite.cn/meow-web/cmd"
)

//go:generate swag init -o ./docs -g main.go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		cmd.ExecuteContext(ctx)
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		fmt.Println("receive stop signal")
	case <-done:
	}

	cancel()
}
