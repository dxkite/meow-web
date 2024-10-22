package app

import (
	"context"

	"dxkite.cn/meow-web/src/user"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/httpx/router"
)

func init() {
	depends.Register(user.NewUserRepository)
	depends.Register(user.NewSessionRepository)
	depends.Register(user.NewUserService)
	depends.Register(user.NewUserHttpServer)

	routeCollection.Add(func(ctx context.Context) (router.Collection, error) {
		return depends.Resolve[*user.UserHttpServer](ctx)
	})
}
