package cmd

import (
	"context"

	"dxkite.cn/meownest/cmd/main"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nest-cli",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		main.ExecuteContext(cmd.Context())
	},
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
