package httpserver

import (
	"context"

	"dxkite.cn/meownest/pkg/container"
	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/pkg/httputil/router"
	"dxkite.cn/meownest/src/config"
	"dxkite.cn/meownest/src/monitor"
)

func init() {
	container.Register(func(cfg *config.Config) *monitor.MonitorConfig {
		return &monitor.MonitorConfig{
			Interval:     cfg.MonitorInterval,
			RollInterval: cfg.MonitorRollInterval,
			MaxInterval:  cfg.MonitorRealtimeInterval,
		}
	})
	container.Register(monitor.NewDynamicStatRepository)
	container.Register(monitor.NewMonitorService)
	container.Register(monitor.NewMonitorServer)

	routeCollection.Add(func(ctx context.Context) (router.Collection, error) {
		ds, _ := container.Get[database.DataSource](ctx)
		service, _ := container.Get[monitor.MonitorService](ctx)
		go service.Collection(database.With(ctx, ds))
		return container.Get[*monitor.MonitorServer](ctx)
	})

}
