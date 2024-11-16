package middleware

import (
	"context"

	"dxkite.cn/nebula/pkg/database"
	"dxkite.cn/nebula/pkg/depends"
	"github.com/gin-gonic/gin"
)

func DataSource(scopeCtx context.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ds, _ := depends.Resolve[database.DataSource](scopeCtx)
		ctx.Request = ctx.Request.WithContext(database.With(ctx.Request.Context(), ds))
	}
}
