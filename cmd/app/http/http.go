package http

import (
	"net/http"
	"time"

	"dxkite.cn/meow-web/pkg/middleware"
	"dxkite.cn/nebula/pkg/httpx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "dxkite.cn/meow-web/cmd/app/depends"
)

var engine = gin.Default()

func init() {
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

	engine.Use(middleware.DataSource())
	engine.Use(middleware.Auth())
}

func Handle(handler http.HandlerFunc, middleware ...httpx.Middleware) gin.HandlerFunc {
	handle := httpx.WrapMiddleware(handler, middleware...)
	return func(c *gin.Context) {
		vars := map[string]string{}
		for _, v := range c.Params {
			vars[v.Key] = v.Value
		}
		httpx.SetPathVars(c.Request, vars)
		handle(c.Writer, c.Request)
	}
}

func Run(addr string) error {
	return engine.Run(addr)
}
