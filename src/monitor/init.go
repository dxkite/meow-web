package monitor

import (
	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/nebula/pkg/depends"
)

func init() {
	depends.Register(func(cfg *config.Config) *MonitorConfig {
		return &MonitorConfig{
			Interval:     cfg.MonitorInterval,
			RollInterval: cfg.MonitorRollInterval,
			MaxInterval:  cfg.MonitorRealtimeInterval,
		}
	})
	depends.Register(NewDynamicStatRepository)
	depends.Register(NewMonitorService)
}
