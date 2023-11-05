package suda

import (
	"strings"
)

type RouteTarget interface {
	RouteName() string
}

type Router struct {
	routes map[string][]RouteTarget
}

func NewRouter() *Router {
	r := &Router{}
	r.routes = map[string][]RouteTarget{}
	return r
}

func (r *Router) Add(uri string, target RouteTarget) *Router {
	if r.routes[uri] == nil {
		r.routes[uri] = []RouteTarget{}
	}
	r.routes[uri] = append(r.routes[uri], target)
	return r
}

func (r Router) MatchAll(uri string) (string, []RouteTarget) {
	// 完整路由
	if routes, ok := r.routes[uri]; ok && routes != nil {
		return uri, routes
	}

	// 前缀路由
	for k, v := range r.routes {
		if strings.HasPrefix(uri, k) && v != nil {
			return k, v
		}
	}
	return "", nil
}

func (r Router) Match(uri string) (string, RouteTarget) {
	p, rr := r.MatchAll(uri)
	if rr == nil {
		return p, nil
	}
	i := intn(len(rr))
	return p, rr[i]
}
