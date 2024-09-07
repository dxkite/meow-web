package cmd

import (
	"os/user"

	provider "dxkite.cn/meownest/pkg/config"
	"dxkite.cn/meownest/pkg/config/env"
	"dxkite.cn/meownest/pkg/database/sqlite"
	"dxkite.cn/meownest/src/config"
	"dxkite.cn/meownest/src/monitor"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{}
		if err := provider.Bind(env.Name, cfg); err != nil {
			panic(err)
		}

		ds, err := sqlite.Open(cfg.DataPath)
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
