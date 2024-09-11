package migrate

import (
	"dxkite.cn/meow-web/pkg/config"
	provider "dxkite.cn/nebula/pkg/config"
	"dxkite.cn/nebula/pkg/config/env"
	"dxkite.cn/nebula/pkg/database/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dst []interface{}

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{}
		if err := provider.Bind(env.NewProvider(), cfg); err != nil {
			panic(err)
		}

		ds, err := sqlite.NewSource(cfg.DataPath)
		if err != nil {
			panic(err)
		}

		db := ds.Engine().(*gorm.DB)
		db = db.Debug()
		db.AutoMigrate(dst...)
	},
}
