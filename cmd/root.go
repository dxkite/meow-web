package cmd

import (
	"context"

	"dxkite.cn/meow-web/cmd/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "meow-cli",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app.ExecuteContext(cmd.Context())
	},
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
