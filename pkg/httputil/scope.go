package httputil

import (
	"context"
	"net/http"

	"dxkite.cn/meownest/pkg/errors"
	"dxkite.cn/meownest/pkg/httputil/router"
)

type ScopeName string

type ScopeContext interface {
	Scopes() []ScopeName
	Identity() uint64
}

type scopeContextKeyType string

var scopeContextKey scopeContextKeyType = "scope-context"

func WithScope(ctx context.Context, scope ScopeContext) context.Context {
	return context.WithValue(ctx, scopeContextKey, scope)
}

func Scope(ctx context.Context) ScopeContext {
	if v, ok := ctx.Value(scopeContextKey).(ScopeContext); ok {
		return v
	}
	return &scopeContext{}
}

func NewScope(ident uint64, scopes ...ScopeName) ScopeContext {
	return &scopeContext{id: ident, scopes: scopes}
}

type scopeContext struct {
	id     uint64
	scopes []ScopeName
}

func (s *scopeContext) Scopes() []ScopeName {
	return s.scopes
}

func (s *scopeContext) Identity() uint64 {
	return s.id
}

func ScopeRequired(scopes ...ScopeName) router.Wrapper {
	return func(r router.Route) router.Route {
		handler := r.Handler()
		return router.New(r.Method(), r.Path(), router.Handler(func(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
			scopeCtx := Scope(ctx)
			// 必须登录状态
			if scopeCtx.Identity() == 0 {
				Error(ctx, w, errors.Unauthorized(errors.Errorf("identity required")))
				return
			}

			// 是否需要进行权限判断
			if len(scopes) == 0 {
				// 执行 next
				handler(ctx, req, w, vars)
				return
			}

			// 任意权限
			if inScopeSlice("*", scopeCtx.Scopes()) {
				// 执行 next
				handler(ctx, req, w, vars)
				return
			}

			// 检查权限列表
			for _, scope := range scopes {
				if !inScopeSlice(scope, scopeCtx.Scopes()) {
					Error(ctx, w, errors.Unauthorized(errors.Errorf("scope %s required", scope)))
					return
				}
			}
			// 执行 next
			handler(ctx, req, w, vars)
		}))
	}
}

func inScopeSlice(a ScopeName, arr []ScopeName) bool {
	for _, v := range arr {
		if v == a {
			return true
		}
	}
	return false
}
