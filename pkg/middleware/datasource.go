package middleware

import (
	"dxkite.cn/nebula/pkg/database"
	"github.com/gin-gonic/gin"
)

func DataSource(ds database.DataSource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(database.With(ctx.Request.Context(), ds))
	}
}
