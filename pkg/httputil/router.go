package httputil

import (
	"context"
	"net/http"
)

type RouteHandler func(ctx context.Context, req *http.Request, resp http.ResponseWriter, vars map[string]string)

type Route interface {
	Path() string
	Method() string
	Handler() RouteHandler
}

type RouteWrapper func(r Route) Route

type route struct {
	path    string
	method  string
	handler RouteHandler
}

func (r *route) Path() string {
	return r.path
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Handler() RouteHandler {
	return r.handler
}

func NewRoute(method, path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	var r Route = &route{method: method, path: path, handler: handler}
	for _, w := range wrappers {
		r = w(r)
	}
	return r
}

func NewGetRoute(path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	return NewRoute(http.MethodGet, path, handler, wrappers...)
}

func NewPostRoute(path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	return NewRoute(http.MethodPost, path, handler, wrappers...)
}

func NewDeleteRoute(path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	return NewRoute(http.MethodDelete, path, handler, wrappers...)
}

func NewPutRoute(path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	return NewRoute(http.MethodPut, path, handler, wrappers...)
}

func NewOptionRoute(path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	return NewRoute(http.MethodOptions, path, handler, wrappers...)
}

func NewHeadRoute(path string, handler RouteHandler, wrappers ...RouteWrapper) Route {
	return NewRoute(http.MethodHead, path, handler, wrappers...)
}
