package httpserver

import (
	"context"
	"strings"
	"time"

	"dxkite.cn/meownest/pkg/config/env"
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

func ExecuteContext(ctx context.Context) {
	configProvider, err := env.NewDotEnvConfig()
	if err != nil {
		panic(err)
	}

	cfg := config.Config{}
	configProvider.Bind(&cfg)

	ds, err := sqlite.Open(cfg.DataPath)
	if err != nil {
		panic(err)
	}

	db := ds.Engine().(*gorm.DB)
	db.AutoMigrate(user.User{}, monitor.DynamicStat{})

	SessionIdName := cfg.SessionName

	userRepository := user.NewUserRepository()
	sessionRepository := user.NewSessionRepository()
	userService := user.NewUserService(userRepository, sessionRepository, []byte(cfg.SessionCryptoKey))
	userServer := user.NewUserHttpServer(userService, SessionIdName)

	monitorRepository := monitor.NewDynamicStatRepository()
	// 5秒 统计一次，记录最新1小时数据，5分钟聚合一次
	monitorService := monitor.NewMonitorService(&monitor.MonitorConfig{
		Interval:     cfg.MonitorInterval,
		RollInterval: cfg.MonitorRollInterval,
		MaxInterval:  cfg.MonitorRealtimeInterval,
	}, monitorRepository)
	monitorServer := monitor.NewMonitorServer(monitorService)

	go monitorService.Collection(database.With(context.Background(), ds))

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
		cookie, _ := ctx.Cookie(SessionIdName)
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
		scope, err := userService.GetSession(ctx, tks[1])
		if err != nil {
			httputil.Error(ctx, ctx.Writer, errors.System(err))
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(httputil.WithScope(ctx.Request.Context(), scope))
	})

	const APIBase = "/api/v1"
	applyRoute(engine.Group(APIBase), userServer.Routes())
	applyRoute(engine.Group(APIBase), monitorServer.Routes())
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
