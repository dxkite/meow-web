package middleware

import (
	"context"
	"net/http"
	"strings"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/src/user"
	"dxkite.cn/nebula/pkg/depends"
	"dxkite.cn/nebula/pkg/errorx"
	"dxkite.cn/nebula/pkg/httpx"
	"github.com/gin-gonic/gin"
)

func getAccessToken(req *http.Request, sessionName string) string {
	auth := req.Header.Get("Authorization")
	if auth != "" {
		tks := strings.SplitN(auth, " ", 2)
		if tks[0] == "Bearer" {
			return tks[1]
		}
	}

	cookie, _ := req.Cookie(sessionName)
	if cookie != nil {
		return cookie.Value
	}

	return ""
}

func Auth(scopeCtx context.Context) gin.HandlerFunc {
	userService, _ := depends.Resolve[user.UserService](scopeCtx)
	cfg, _ := depends.Resolve[*config.Config](scopeCtx)

	return func(ctx *gin.Context) {
		token := getAccessToken(ctx.Request, cfg.SessionName)
		if token == "" {
			ctx.Next()
			return
		}

		scope, err := userService.GetSession(ctx, token)
		if err != nil {
			httpx.Error(ctx.Writer, errorx.System(err))
			ctx.Abort()
			return
		}
		ctx.Request = ctx.Request.WithContext(httpx.WithScope(ctx.Request.Context(), scope))
	}
}
