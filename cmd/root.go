package cmd

import (
	"context"

	app "dxkite.cn/meownest/src"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "meownest",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app.ExecuteContext(cmd.Context())
	},
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
