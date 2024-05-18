package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type Monitor struct {
	s service.Monitor
}

func NewMonitor(s service.Monitor) *Monitor {
	return &Monitor{s: s}
}

func (s *Monitor) GetLoadStat(c *gin.Context) {
	stat, err := s.s.GetRealtimeStat(c)

	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.Result(c, http.StatusOK, stat)
}

func (s *Monitor) RegisterToHttp(r gin.IRouter) {
	r.GET("/monitor/realtime-stat", s.GetLoadStat)
}
