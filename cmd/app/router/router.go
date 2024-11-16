package router

import (
	"net/http"
	"time"

	"dxkite.cn/meow-web/cmd/app/depends"
	"dxkite.cn/meow-web/pkg/middleware"
	"dxkite.cn/nebula/pkg/httpx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Engine = gin.Default()

func init() {
	Engine.ContextWithFallback = true
	Engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	Engine.Use(middleware.DataSource(depends.Scope))
	Engine.Use(middleware.Auth(depends.Scope))
}

func Wrap(handle http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		vars := map[string]string{}
		for _, v := range c.Params {
			vars[v.Key] = v.Value
		}
		httpx.SetPathVars(c.Request, vars)
		handle(c.Writer, c.Request)
	}
}
