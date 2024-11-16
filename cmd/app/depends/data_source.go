package depends

import (
	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/nebula/pkg/database"
	"dxkite.cn/nebula/pkg/database/sqlite"
	"dxkite.cn/nebula/pkg/depends"
)

func init() {
	depends.Register(func(cfg *config.Config) (database.DataSource, error) {
		return sqlite.NewSource(cfg.DataPath, sqlite.WithDebug(cfg.Env == config.EnvDevelopment))
	})
}
