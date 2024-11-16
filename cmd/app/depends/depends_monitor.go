package depends

import (
	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/src/monitor"
	"dxkite.cn/nebula/pkg/depends"
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
}
