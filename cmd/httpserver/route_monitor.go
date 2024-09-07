package httpserver

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/pkg/depends"
	"dxkite.cn/meownest/pkg/httputil/router"
	"dxkite.cn/meownest/src/config"
	"dxkite.cn/meownest/src/monitor"
)

func init() {
	depends.Register(func(cfg *config.Config) *monitor.MonitorConfig {
		return &monitor.MonitorConfig{
			Interval:     cfg.MonitorInterval,
			RollInterval: cfg.MonitorRollInterval,
			MaxInterval:  cfg.MonitorRealtimeInterval,
		}
	})
	depends.Register(monitor.NewDynamicStatRepository)
	depends.Register(monitor.NewMonitorService)
	depends.Register(monitor.NewMonitorServer)

	routeCollection.Add(func(ctx context.Context) (router.Collection, error) {
		ds, err := depends.Resolve[database.DataSource](ctx)
		if err != nil {
			return nil, err
		}

		service, err := depends.Resolve[monitor.MonitorService](ctx)
		if err != nil {
			return nil, err
		}

		go service.Collection(database.With(ctx, ds))

		return depends.Resolve[*monitor.MonitorServer](ctx)
	})

}
