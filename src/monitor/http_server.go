package monitor

import (
	"net/http"

	"dxkite.cn/nebula/pkg/httpx"
	"dxkite.cn/nebula/pkg/httpx/router"
)

type MonitorServer struct {
	s MonitorService
}

func NewMonitorServer(s MonitorService) *MonitorServer {
	return &MonitorServer{s: s}
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
// @Success      200  {object} DynamicStatResult
// @Failure      400  {object} httpx.HttpError
// @Failure      500  {object} httpx.HttpError
// @Router       /monitor/dynamic-stat [get]
func (s *MonitorServer) ListDynamicStat(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var param ListDynamicStatRequest

	if err := httpx.ReadQuery(req, &param); err != nil {
		httpx.Error(w, err)
		return
	}

	if err := httpx.Validate(&param); err != nil {
		httpx.Error(w, err)
		return
	}

	rst, err := s.s.ListDynamicStat(ctx, &param)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.Result(w, http.StatusOK, rst)
}

func (s *MonitorServer) Routes() []router.Route {
	return []router.Route{
		router.GET("/api/v1/monitor/dynamic-stat", s.ListDynamicStat),
	}
}
