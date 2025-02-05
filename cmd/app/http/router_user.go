package http

import (
	"dxkite.cn/meow-web/src/user"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/httpx"
)

func init() {

	service, err := depends.Resolve[user.UserService]()
	if err != nil {
		panic(err)
	}

	engine.POST("/api/v1/users/session", Handle(httpx.HandleRet(service.CreateSession)))
	engine.DELETE("/api/v1/users/session", Handle(httpx.Handle(service.DeleteSession), httpx.ScopeRequired()))
}
