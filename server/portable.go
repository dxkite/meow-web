package server

import (
	"context"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/gateway/util"
	"dxkite.cn/log"
)

type Portable struct {
	*Server
	Config *config.Config
	Route  *route.Route
}

func NewPortable(ctx context.Context, cfg *config.Config) *Portable {
	pp := &Portable{}

	pp.Config = cfg
	util.SetLogConfig(ctx, cfg.LogConfig.Level, cfg.LogConfig.Path)
	log.Println("load config success")
	r := route.NewRoute()
	pp.Route = r
	r.Load(cfg.Routes)

	s := NewServer(cfg, r)
	pp.Server = s

	cfg.OnChange(func(c interface{}) {
		cc := c.(*config.Config)
		cfg.SetLastLoadTime(cc.HotLoad)
		r.ClearAll()
		r.Load(cc.Routes)
		r.ApplyDynamic()
		s.ApplyHeaderFilter(cc.HttpAllowHeader)
		s.ApplyCorsConfig(cc.Cors)
		err := s.InitTicketMode(cc.Session().Mode)
		if err != nil {
			log.Error("load ticket mode error", err)
		}
	})

	if cfg.HotLoad > 0 {
		cfg.SetLastLoadTime(cfg.HotLoad)
		cfg.HotLoadIfModify()
	}

	// 触发配置刷新
	cfg.NotifyModify()
	return pp
}
