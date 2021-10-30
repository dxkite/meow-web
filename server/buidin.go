package server

import (
	"dxkite.cn/gateway/proto"
	"dxkite.cn/gateway/route"
	"net/http"
	"strconv"
)

type handler struct {
	cfg *route.RouteConfig
	hd  http.Handler
}

type Builder interface {
	Build() proto.Processor
}

func NewHandler(cfg *route.RouteConfig, hd http.Handler) route.RouteItem {
	return &handler{
		cfg: cfg,
		hd:  hd,
	}
}

func (h *handler) Config() *route.RouteConfig {
	return h.cfg
}

func (h *handler) Backend() route.RouteBackend {
	return h
}

func (h *handler) Get() route.Backend {
	return h
}

func (h *handler) BackendType() string {
	return "build-in"
}

func (h *handler) Build() proto.Processor {
	return h
}

// 链接后端
func (h *handler) Do(ctx *proto.BackendContext, w http.ResponseWriter, r *http.Request) (err error) {
	r.Header.Set(ctx.Cfg.UinHeaderName, strconv.Itoa(int(ctx.Uin)))
	r.Header.Set("Authorization", ctx.Ticket)
	h.hd.ServeHTTP(w, r)
	return
}
