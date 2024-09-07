package middleware

import (
	"context"
	"strings"

	"dxkite.cn/meownest/pkg/config"
	"dxkite.cn/meownest/src/user"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/errors"
	"dxkite.cn/nebula/pkg/httputil"
	"github.com/gin-gonic/gin"
)

func Auth(scopeCtx context.Context, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userService, _ := depends.Resolve[user.UserService](scopeCtx)

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

		scope, err := userService.GetSession(ctx, tks[1])
		if err != nil {
			httputil.Error(ctx, ctx.Writer, errors.System(err))
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(httputil.WithScope(ctx.Request.Context(), scope))
	}
}
