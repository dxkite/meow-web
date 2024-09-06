package httpserver

import (
	"context"
	"strings"
	"time"

	"dxkite.cn/meownest/pkg/config/env"
	"dxkite.cn/meownest/pkg/container"
	"dxkite.cn/meownest/pkg/crypto/identity"
	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/pkg/database/sqlite"
	"dxkite.cn/meownest/pkg/errors"
	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/pkg/httputil/router"
	"dxkite.cn/meownest/src/config"
	"dxkite.cn/meownest/src/monitor"
	"dxkite.cn/meownest/src/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	identity.DefaultMask = 1234627081864056831
}

var routeCollection = router.NewCollectionBag()

func ExecuteContext(ctx context.Context) {
	configProvider, err := env.NewDotEnvConfig()
	if err != nil {
		panic(err)
	}

	instanceCtx := container.NewScopedContext(ctx)

	cfg := config.Config{}
	configProvider.Bind(&cfg)

	ds, err := sqlite.Open(cfg.DataPath)
	if err != nil {
		panic(err)
	}

	db := ds.Engine().(*gorm.DB)
	db.AutoMigrate(user.User{}, monitor.DynamicStat{})

	container.Register(func() database.DataSource { return ds })
	container.Register(func() *config.Config { return &cfg })

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

	engine.Use(func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(database.With(ctx.Request.Context(), ds))
	})

	engine.Use(func(ctx *gin.Context) {
		cookie, _ := ctx.Cookie(cfg.SessionName)
		if cookie == "" {
			ctx.Next()
			return
		}

		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.Next()
			return
		}

		tks := strings.SplitN(auth, " ", 2)
		if tks[0] != "Bearer" {
			httputil.Error(ctx, ctx.Writer, errors.Unauthorized(errors.Errorf("invalid token type %s", tks[0])))
			ctx.Abort()
			return
		}
		userService, _ := container.Get[user.UserService](instanceCtx)
		scope, err := userService.GetSession(ctx, tks[1])
		if err != nil {
			httputil.Error(ctx, ctx.Writer, errors.System(err))
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(httputil.WithScope(ctx.Request.Context(), scope))
	})

	const APIBase = "/api/v1"

	routes, err := routeCollection.Build(instanceCtx)
	if err != nil {
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
