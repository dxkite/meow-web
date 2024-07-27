package monitor

import (
	"context"
	"net/http"

	"dxkite.cn/meownest/pkg/httputil"
	"dxkite.cn/meownest/pkg/httputil/router"
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
// @Failure      400  {object} httputil.HttpError
// @Failure      500  {object} httputil.HttpError
// @Router       /monitor/dynamic-stat [get]
func (s *MonitorServer) ListDynamicStat(ctx context.Context, req *http.Request, w http.ResponseWriter, vars map[string]string) {
	var param ListDynamicStatRequest

	if err := httputil.ReadQuery(ctx, req, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	if err := httputil.Validate(ctx, &param); err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	rst, err := s.s.ListDynamicStat(ctx, &param)
	if err != nil {
		httputil.Error(ctx, w, err)
		return
	}

	httputil.Result(ctx, w, http.StatusOK, rst)
}

func (s *MonitorServer) Routes() []router.Route {
	return []router.Route{
		router.GET("/monitor/dynamic-stat", s.ListDynamicStat),
	}
}
