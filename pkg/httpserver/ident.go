package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IdentityConfig struct {
	Ident func(ctx *gin.Context) (id uint64, scopes []string, err error)
}

func Identity(cfg IdentityConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, scopes, err := cfg.Ident(ctx)
		if err != nil {
			Error(ctx, http.StatusBadRequest, "invalid_token", err.Error())
			ctx.Abort()
			return
		}
		ctx.Set("ident", id)
		ctx.Set("scopes", scopes)
	}
}

func inStringSlice(a string, arr []string) bool {
	for _, v := range arr {
		if v == a {
			return true
		}
	}
	return false
}

func ScopeRequired(scopes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		curScopes := ctx.GetStringSlice("scopes")
		// 任意权限
		if inStringSlice("*", curScopes) {
			return
		}
		// 检查权限列表
		for _, scope := range scopes {
			if !inStringSlice(scope, curScopes) {
				Error(ctx, http.StatusUnauthorized, "invalid_scope", fmt.Sprintf("scope %s required", scope))
				ctx.Abort()
				return
			}
		}
	}
}

func IdentityRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ident := ctx.GetUint64("ident")
		if ident == 0 {
			Error(ctx, http.StatusUnauthorized, "invalid_ident", "identity required")
			ctx.Abort()
			return
		}
	}
}
