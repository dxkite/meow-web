package proto

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"net/http"
)

type BackendContext struct {
	Cfg     *config.Config
	Uin     uint64
	Ticket  string
	Route   *route.RouteInfo
	Backend route.Backend
}

type Processor interface {
	// 链接后端
	Do(ctx *BackendContext, w http.ResponseWriter, req *http.Request) (err error)
}
