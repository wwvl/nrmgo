package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"nrmgo/internal/registry"
	"nrmgo/internal/style"
	"nrmgo/internal/table"
)

var addCmd = &cobra.Command{
	Use:   "add <registry-name> <registry-url> [registry-home] [registry-description]",
	Short: "Add a custom registry",
	Long: `Add a custom registry with name and URL.
Registry name can only contain letters, numbers and underscores.
Registry URL must be a valid HTTP/HTTPS URL.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取参数
		name := args[0]
		url := args[1]
		home := ""
		description := ""
		if len(args) > 2 {
			home = args[2]
		}
		if len(args) > 3 {
			description = args[3]
		}

		// 加载配置并创建管理器
		_, manager, err := loadConfigAndCreateManager()
		if err != nil {
			return err
		}

		// 检查名称是否已存在
		if _, exists := manager.Get(name); exists {
			return fmt.Errorf("\n❌  Registry '%s' already exists", name)
		}

		// 检查是否为内置 registry 名称
		if _, isBuiltin := registry.GetBuiltinRegistry(name); isBuiltin {
			return fmt.Errorf("\n❌  Cannot use built-in registry name: %s", name)
		}

		// 验证名称格式
		if err := registry.IsValidName(name); err != nil {
			return fmt.Errorf("\n%v", err)
		}

		// 验证 URL 格式
		if err := registry.IsValidURL(url); err != nil {
			return fmt.Errorf("\n%v", err)
		}

		// 创建新的 registry
		reg := registry.NewRegistry(name, url, home, description)

		// 添加 registry
		if err := manager.Add(name, reg); err != nil {
			return fmt.Errorf("\n❌  Failed to add registry: %v", err)
		}

		// 创建表格渲染器
		renderer := table.NewTableRenderer([]string{
			"Name",
			"URL",
			"Home",
			"Description",
		})

		// 添加数据行
		renderer.MustAddRow([]string{
			reg.Name,
			reg.URL,
			reg.Home,
			reg.Description,
		})

		// 渲染表格
		fmt.Println()
		if err := renderer.Render(); err != nil {
			return fmt.Errorf("\n❌  Failed to render table: %v", err)
		}

		// 输出成功信息
		fmt.Printf("\n✨ Successfully Added Registry: %s => %s\n",
			style.Success.Sprint(name),
			style.Success.Sprint(url))

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
