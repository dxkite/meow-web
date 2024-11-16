package router

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

	Engine.POST("/api/v1/users/session", Wrap(httpx.HandleRet(service.CreateSession)))
	Engine.DELETE("/api/v1/users/session", Wrap(httpx.WrapMiddleware(httpx.Handle(service.DeleteSession), httpx.ScopeRequired())))
}
