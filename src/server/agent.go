package server

import (
	"dxkite.cn/meownest/pkg/agent"
	"github.com/gin-gonic/gin"
)

type Agent struct {
	svr *agent.Server
}

func NewAgent(svr *agent.Server) *Agent {
	return &Agent{svr: svr}
}

func (s *Agent) Run(addr string) {
	s.svr.Run(addr)
}

func (s *Agent) RegisterToHttp(group gin.IRouter) {
}
