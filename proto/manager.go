package proto

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"io"
	"net/http"
)

type BackendContext struct {
	Cfg     *config.Config
	Req     *http.Request
	Writer  http.ResponseWriter
	Route   *route.RouteInfo
	Backend route.Backend
}

type Processor interface {
	// 链接后端
	Do(uin uint64, ticket string) (user uint64, status int, header http.Header, body io.ReadCloser, err error)
}
