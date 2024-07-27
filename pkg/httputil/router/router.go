package router

import (
	"context"
	"net/http"
)

type Handler func(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string)

type Route interface {
	Path() string
	Method() string
	Handler() Handler
}

type Wrapper func(r Route) Route

type route struct {
	path    string
	method  string
	handler Handler
}

func (r *route) Path() string {
	return r.path
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Handler() Handler {
	return r.handler
}

func New(method, path string, handler Handler, wrappers ...Wrapper) Route {
	var r Route = &route{method: method, path: path, handler: handler}
	for _, w := range wrappers {
		r = w(r)
	}
	return r
}

func GET(path string, handler Handler, wrappers ...Wrapper) Route {
	return New(http.MethodGet, path, handler, wrappers...)
}

func POST(path string, handler Handler, wrappers ...Wrapper) Route {
	return New(http.MethodPost, path, handler, wrappers...)
}

func DELETE(path string, handler Handler, wrappers ...Wrapper) Route {
	return New(http.MethodDelete, path, handler, wrappers...)
}

func PUT(path string, handler Handler, wrappers ...Wrapper) Route {
	return New(http.MethodPut, path, handler, wrappers...)
}

func OPTIONS(path string, handler Handler, wrappers ...Wrapper) Route {
	return New(http.MethodOptions, path, handler, wrappers...)
}

func HEAD(path string, handler Handler, wrappers ...Wrapper) Route {
	return New(http.MethodHead, path, handler, wrappers...)
}
