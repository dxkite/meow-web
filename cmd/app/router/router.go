package router

import (
	"net/http"

	"dxkite.cn/nebula/pkg/httpx"
	"github.com/gin-gonic/gin"
)

var Engine = gin.Default()

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
