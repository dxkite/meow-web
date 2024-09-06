package httpserver

import (
	"context"

	"dxkite.cn/meownest/pkg/container"
	"dxkite.cn/meownest/pkg/httputil/router"
	"dxkite.cn/meownest/src/user"
)

func init() {

	container.Register(user.NewUserRepository)
	container.Register(user.NewSessionRepository)
	container.Register(user.NewUserService)
	container.Register(user.NewUserHttpServer)

	routeCollection.Add(func(ctx context.Context) (router.Collection, error) {
		return container.Get[*user.UserHttpServer](ctx)
	})
}
