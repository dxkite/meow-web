package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type Agent struct {
	s service.Agent
}

func NewAgent(s service.Agent) *Agent {
	return &Agent{s: s}
}

func (s *Agent) Reload(c *gin.Context) {
	s.s.LoadRoute(c)
	c.Status(http.StatusOK)
}

func (s *Agent) RegisterToHttp(r gin.IRouter) {
	r.POST("/agent/reload", s.Reload)
}
