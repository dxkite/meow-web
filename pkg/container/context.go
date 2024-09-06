package container

import "context"

type Scope interface {
	Get(id InstanceId) any
	Set(id InstanceId, val any)
}

type scope struct {
	scope map[InstanceId]any
}

func (s *scope) Get(id InstanceId) any {
	return s.scope[id]
}

func (s *scope) Set(id InstanceId, val any) {
	s.scope[id] = val
}

func NewScope() Scope {
	return &scope{scope: make(map[InstanceId]any)}
}

type scopeKey struct{}

func NewScopedContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, scopeKey{}, NewScope())
}
