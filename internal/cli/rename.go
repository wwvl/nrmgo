package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"nrmgo/internal/registry"
	"nrmgo/internal/style"
)

// 定义全局变量
var (
	forceRename bool
)

// renameCmd 重命名 registry 命令
var renameCmd = &cobra.Command{
	Use:   "rename <old-name> <new-name>",
	Short: "Rename a registry in the configuration",
	Long: `Rename a registry in the configuration.
Note: Built-in registries cannot be renamed.`,
	Example: `  # Rename a custom registry
  nrmgo rename old-registry new-registry

  # Rename a registry without confirmation
  nrmgo rename old-registry new-registry --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 验证参数
		if len(args) != 2 {
			return fmt.Errorf("\n❌  Both old and new registry names are required")
		}

		oldName := args[0]
		newName := args[1]

		// 加载配置并创建管理器
		_, manager, err := loadConfigAndCreateManager()
		if err != nil {
			return err
		}

		// 获取当前使用的 registry
		current, _ := manager.Current()

		// 检查 old registry 是否存在
		_, exists := manager.Get(oldName)
		if !exists {
			return fmt.Errorf("\n❌  Registry '%s' not found", oldName)
		}

		// 检查是否为内置 registry
		if _, isBuiltin := registry.GetBuiltinRegistry(oldName); isBuiltin {
			return fmt.Errorf("\n❌  Cannot rename built-in registry: %s", oldName)
		}

		// 检查新名称是否已存在
		if _, exists := manager.Get(newName); exists {
			return fmt.Errorf("\n❌  Registry '%s' already exists", newName)
		}

		// 检查新名称是否为内置 registry 名称
		if _, isBuiltin := registry.GetBuiltinRegistry(newName); isBuiltin {
			return fmt.Errorf("\n❌  Cannot use built-in registry name: %s", newName)
		}

		// 验证新名称格式
		if err := registry.IsValidName(newName); err != nil {
			return fmt.Errorf("\n%v", err)
		}

		// 如果正在使用中且未使用 force 参数，需要确认
		if !forceRename && current != nil && current.Name == oldName {
			fmt.Printf("\nRegistry '%s' is currently in use. Are you sure to rename it? [y/N]: ", oldName)
			var answer string
			if _, err := fmt.Scanln(&answer); err != nil {
				answer = "n"
			}
			if !strings.EqualFold(answer, "y") {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		// 执行重命名
		if err := manager.Rename(oldName, newName); err != nil {
			return fmt.Errorf("\n❌  %v", err)
		}

		fmt.Printf("\n✨ Successfully renamed registry from '%s' to '%s'\n",
			style.Success.Sprint(oldName),
			style.Success.Sprint(newName))
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// 添加命令行参数
	flags := renameCmd.Flags()
	flags.BoolVarP(&forceRename, "force", "f", false, "Force rename without confirmation")
}
