package gateway

import (
	"net/http"
	"sort"
	"strings"

	"dxkite.cn/meownest/src/utils"
)

type Router struct {
	routes map[string][]http.Handler
	uris   []string
}

type RouteMatcher interface {
	MatchRequest(req *http.Request) bool
}

func NewRouter() *Router {
	r := &Router{}
	r.routes = map[string][]http.Handler{}
	r.uris = []string{}
	return r
}

func (r *Router) Add(uri string, target http.Handler) *Router {
	if r.routes[uri] == nil {
		r.routes[uri] = []http.Handler{}
		r.uris = append(r.uris, uri)
	}
	r.routes[uri] = append(r.routes[uri], target)
	return r
}

func (r *Router) sort() {
	sort.Slice(r.uris, func(a, b int) bool {
		if len(r.uris[a]) > len(r.uris[b]) {
			return true
		}
		return r.uris[a] > r.uris[b]
	})
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.sort()
	if _, m := r.match(req); m != nil {
		m.ServeHTTP(w, req)
		return
	}
	http.NotFound(w, req)
}

func (r Router) matchAll(uri string) (string, []http.Handler) {
	// 完整路由
	if route, ok := r.routes[uri]; ok && route != nil {
		return uri, route
	}
	// 前缀路由
	for _, v := range r.uris {
		if strings.HasPrefix(uri, v) {
			return v, r.routes[v]
		}
	}
	return "", nil
}

func (r Router) match(req *http.Request) (string, http.Handler) {
	p, rr := r.matchAll(req.URL.Path)
	if rr == nil {
		return p, nil
	}

	for _, v := range rr {
		if m, ok := v.(RouteMatcher); ok && m.MatchRequest(req) {
			return p, v
		}
	}
	i := utils.Intn(len(rr))
	return p, rr[i]
}
