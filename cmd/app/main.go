package app

import (
	"context"
	"fmt"
	"time"

	provider "dxkite.cn/nebula/pkg/config"
	"dxkite.cn/nebula/pkg/config/env"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/pkg/middleware"
	"dxkite.cn/nebula/pkg/crypto/identity"
	"dxkite.cn/nebula/pkg/database/sqlite"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/httputil/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	identity.DefaultMask = 1234627081864056831
}

var routeCollection = router.NewCollectionBag()

func ExecuteContext(ctx context.Context) {
	cfg := &config.Config{}
	if err := provider.Bind(env.NewProvider(), cfg); err != nil {
		panic(err)
	}

	scopeCtx := depends.NewScopedContext(ctx)

	ds, err := sqlite.NewSource(cfg.DataPath)
	if err != nil {
		panic(err)
	}

	depends.Register(ds)
	depends.Register(cfg)

	engine := gin.Default()
	engine.ContextWithFallback = true
	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	engine.Use(middleware.DataSource(ds))
	engine.Use(middleware.Auth(scopeCtx, cfg))

	const APIBase = "/api/v1"

	routes, err := routeCollection.Build(scopeCtx)
	if err != nil {
		fmt.Println(err)
		return
	}

	applyRoute(engine.Group(APIBase), routes)
	engine.Run(cfg.Listen)
}

func applyRoute(engine *gin.RouterGroup, routes []router.Route) {
	for _, r := range routes {
		handler := r.Handler()
		engine.Handle(r.Method(), r.Path(), func(ctx *gin.Context) {
			vars := map[string]string{}
			for _, v := range ctx.Params {
				vars[v.Key] = v.Value
			}
			handler(ctx, ctx.Request, ctx.Writer, vars)
		})
	}
}
