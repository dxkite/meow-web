package cmd

import (
	"context"

	"dxkite.cn/meownest/cmd/httpserver"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "meownest",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		httpserver.ExecuteContext(cmd.Context())
	},
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
