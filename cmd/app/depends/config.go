package depends

import (
	"dxkite.cn/meow-web/pkg/config"
	provider "dxkite.cn/nebula/pkg/config"
	"dxkite.cn/nebula/pkg/config/env"
	"dxkite.cn/nebula/pkg/depends"
)

func init() {
	depends.Register(func() (*config.Config, error) {
		cfg := &config.Config{}
		if err := provider.Bind(env.NewProvider(), cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	})
}
