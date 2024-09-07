package app

import (
	"context"

	"dxkite.cn/meownest/pkg/depends"
	"dxkite.cn/meownest/pkg/httputil/router"
	"dxkite.cn/meownest/src/user"
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
