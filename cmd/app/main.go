package app

import (
	"context"
	"fmt"
	"time"

	provider "dxkite.cn/nebula/pkg/config"
	"dxkite.cn/nebula/pkg/config/env"
	"dxkite.cn/nebula/pkg/httpx"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/pkg/middleware"
	"dxkite.cn/nebula/pkg/crypto/identity"
	"dxkite.cn/nebula/pkg/database/sqlite"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/httpx/router"
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

	ds, err := sqlite.NewSource(cfg.DataPath, sqlite.WithDebug(cfg.Env == config.EnvDevelopment))
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

	routes, err := routeCollection.Build(scopeCtx)
	if err != nil {
		fmt.Println(err)
		return
	}

	applyRoute(engine, routes)
	engine.Run(cfg.Listen)
}

func applyRoute(engine *gin.Engine, routes []router.Route) {
	for _, r := range routes {
		handle := r.Handle()
		engine.Handle(r.Method(), r.Path(), func(c *gin.Context) {
			vars := map[string]string{}
			for _, v := range c.Params {
				vars[v.Key] = v.Value
			}
			httpx.SetVars(c.Request, vars)
			handle(c.Writer, c.Request)
		})
	}
}
