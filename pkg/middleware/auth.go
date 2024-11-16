package middleware

import (
	"context"
	"strings"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/src/user"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/errorx"
	"dxkite.cn/nebula/pkg/httpx"
	"github.com/gin-gonic/gin"
)

func Auth(scopeCtx context.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userService, _ := depends.Resolve[user.UserService](scopeCtx)
		cfg, _ := depends.Resolve[*config.Config](scopeCtx)

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
			httpx.Error(ctx.Writer, errorx.Unauthorized(errorx.Errorf("invalid token type %s", tks[0])))
			ctx.Abort()
			return
		}

		scope, err := userService.GetSession(ctx, tks[1])
		if err != nil {
			httpx.Error(ctx.Writer, errorx.System(err))
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(httpx.WithScope(ctx.Request.Context(), scope))
	}
}
