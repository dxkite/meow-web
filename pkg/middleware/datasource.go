package middleware

import (
	"dxkite.cn/nebula/pkg/database"
	"dxkite.cn/nebula/pkg/depends"
	"github.com/gin-gonic/gin"
)

func DataSource() gin.HandlerFunc {
	ds, _ := depends.Resolve[database.DataSource]()
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(database.With(ctx.Request.Context(), ds))
	}
}
