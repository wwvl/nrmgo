package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nrmgo",
	Short: "NPM Registry Manager GO",
	Long: `NPM Registry Manager GO (nrmgo) is a command-line tool for quickly switching npm registries.
Supports npm, yarn, pnpm and bun package managers.`,
	// 禁用自动生成的使用说明
	DisableAutoGenTag: true,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 禁用自动生成的 completion 子命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 禁用自动排序命令
	cobra.EnableCommandSorting = false

	// 禁用命令排序
	rootCmd.Flags().SortFlags = false
}
