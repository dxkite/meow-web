package cmd

import (
	"os/user"

	"dxkite.cn/meow-web/pkg/config"
	"dxkite.cn/meow-web/src/monitor"
	provider "dxkite.cn/nebula/pkg/config"
	"dxkite.cn/nebula/pkg/config/env"
	"dxkite.cn/nebula/pkg/database/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
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
		db.AutoMigrate(user.User{}, monitor.DynamicStat{})
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
