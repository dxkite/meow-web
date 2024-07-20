package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

type Monitor struct {
	s service.Monitor
}

func NewMonitor(s service.Monitor) *Monitor {
	return &Monitor{s: s}
}

// List Dynamic Stat
//
// @Summary      List Dynamic Stat
// @Description  List Dynamic Stat
// @Tags         Monitor
// @Accept       json
// @Produce      json
// @Param        start_time query string false "开始时间。默认-1h"
// @Param		 end_time query string false "结束时间，默认当前时间"
// @Success      200  {object} service.DynamicStatResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /monitor/dynamic-stat [get]
func (s *Monitor) ListDynamicStat(c *gin.Context) {
	var param service.ListDynamicStatParam

	if err := c.ShouldBindQuery(&param); err != nil {
		httputil.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.ListDynamicStat(c, &param)
	if err != nil {
		httputil.ResultError(c, err)
		return
	}

	httputil.Result(c, http.StatusOK, rst)
}

func (s *Monitor) API() httputil.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.GET("/monitor/dynamic-stat", s.ListDynamicStat)
	}
}
