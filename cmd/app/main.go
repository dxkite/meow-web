package app

import (
	"context"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/nebula/pkg/crypto/identity"

	"dxkite.cn/meow-web/cmd/app/depends"
	engine "dxkite.cn/meow-web/cmd/app/router"
)

func init() {
	identity.DefaultMask = 1234627081864056831
}

func ExecuteContext(ctx context.Context) {

	cfg, err := depends.Resolve[*config.Config]()
	if err != nil {
		panic(err)
	}

	engine.Engine.Run(cfg.Listen)
}
