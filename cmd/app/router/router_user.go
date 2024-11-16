package router

import (
	"dxkite.cn/meow-web/cmd/app/depends"
	"dxkite.cn/meow-web/src/user"
	"dxkite.cn/nebula/pkg/httpx"
)

func init() {

	service, err := depends.Resolve[user.UserService]()
	if err != nil {
		panic(err)
	}

	Engine.POST("/api/v1/user/session", Wrap(httpx.HandleRet(service.CreateSession)))
}
