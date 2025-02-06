package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"nrmgo/internal/style"
	"nrmgo/internal/table"
)

var (
	// 命令行参数
	verboseOutput bool // 显示详细信息（包括 Home 和 Description）
	allOutput     bool // 与 verbose 相同，显示详细信息（包括 Home 和 Description）
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available registries",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置并创建管理器
		_, manager, err := loadConfigAndCreateManager()
		if err != nil {
			return err
		}

		// 获取所有 registry
		registries := manager.List()

		// 按名称排序
		sort.Slice(registries, func(i, j int) bool {
			return registries[i].Name < registries[j].Name
		})

		// 获取当前使用的 registry
		current, _ := manager.Current()
		currentName := ""
		if current != nil {
			currentName = current.Name
		}

		// 创建表格渲染器
		var headers []string
		if verboseOutput || allOutput {
			headers = []string{"Name", "URL", "Home", "Description"}
		} else {
			headers = []string{"Name", "URL"}
		}
		renderer := table.NewTableRenderer(headers)

		// 添加数据行
		for _, reg := range registries {
			var row []string
			if verboseOutput || allOutput {
				row = []string{
					reg.Name,
					reg.URL,
					reg.Home,
					reg.Description,
				}
			} else {
				row = []string{
					reg.Name,
					reg.URL,
				}
			}

			// 如果是当前使用的 registry，高亮整行
			if reg.Name == currentName {
				for i := range row {
					row[i] = style.Success.Sprint(row[i])
				}
			}

			renderer.MustAddRow(row)
		}

		// 渲染表格
		fmt.Println()
		if err := renderer.Render(); err != nil {
			return fmt.Errorf("\n❌  Failed to render table: %v", err)
		}
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// 添加命令行参数
	lsCmd.Flags().BoolVarP(&verboseOutput, "verbose", "v", false, "Show detailed information including registry's home page and description")
	lsCmd.Flags().BoolVarP(&allOutput, "all", "a", false, "Same as --verbose, show detailed information")
}
