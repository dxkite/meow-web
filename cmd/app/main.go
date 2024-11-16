package app

import (
	"context"
	"time"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/pkg/middleware"
	"dxkite.cn/nebula/pkg/crypto/identity"
	"github.com/gin-contrib/cors"

	"dxkite.cn/meow-web/cmd/app/depends"
	engine "dxkite.cn/meow-web/cmd/app/router"
)

func init() {
	identity.DefaultMask = 1234627081864056831
}

func ExecuteContext(ctx context.Context) {

	cfg, err := depends.Resolve[*config.Config]()
	if err != nil {
		panic(err)
	}

	engine := engine.Engine
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

	engine.Use(middleware.DataSource(depends.Scope))
	engine.Use(middleware.Auth(depends.Scope))

	engine.Run(cfg.Listen)
}
