package main

import (
	"fmt"

	"dxkite.cn/suda"
)

func main() {
	app := new(suda.App)
	if err := app.Config("./config.yaml"); err != nil {
		fmt.Println(err)
		return
	}
	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
