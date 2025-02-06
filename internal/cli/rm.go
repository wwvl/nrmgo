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
	forceRemove bool
	removeAll   bool
)

// rmCmd 删除 registry 命令
var rmCmd = &cobra.Command{
	Use:   "rm <registry-name>",
	Short: "Remove a registry from the configuration",
	Long: `Remove a registry from the configuration.
Note: Built-in registries cannot be removed.`,
	Example: `  # Remove a custom registry
  nrmgo rm my-registry

  # Remove a registry without confirmation
  nrmgo rm my-registry --force

  # Remove all custom registries
  nrmgo rm --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置并创建管理器
		_, manager, err := loadConfigAndCreateManager()
		if err != nil {
			return err
		}

		// 获取当前使用的 registry
		current, _ := manager.Current()

		// 执行删除操作
		if removeAll {
			return removeAllCustomRegistries(manager, forceRemove)
		}

		if len(args) == 0 {
			return fmt.Errorf("\n❌  Registry name is required")
		}

		return removeRegistry(manager, args[0], forceRemove, current)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// removeRegistry 删除单个 registry
func removeRegistry(manager registry.Manager, name string, force bool, current *registry.Info) error {
	// 检查 registry 是否存在
	if _, exists := manager.Get(name); !exists {
		return fmt.Errorf("\n❌  Registry '%s' not found", name)
	}

	// 检查是否为内置 registry
	if _, isBuiltin := registry.GetBuiltinRegistry(name); isBuiltin {
		return fmt.Errorf("\n❌  Cannot remove built-in registry: %s", name)
	}

	// 如果正在使用中且未使用 force 参数，需要确认
	if !force && current != nil && current.Name == name {
		fmt.Printf("\nRegistry '%s' is currently in use. Are you sure to remove it? [y/N]: ", name)
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil {
			answer = "n"
		}
		if !strings.EqualFold(answer, "y") {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// 执行删除
	if err := manager.Remove(name); err != nil {
		return fmt.Errorf("\n❌  Failed to remove registry: %v", err)
	}

	fmt.Printf("\n✨ Successfully removed registry: %s\n", style.Success.Sprint(name))
	return nil
}

// removeAllCustomRegistries 删除所有自定义 registry
func removeAllCustomRegistries(manager registry.Manager, force bool) error {
	// 获取所有 registry
	registries := manager.List()

	// 过滤出自定义 registry
	var customRegs []*registry.Info
	for _, reg := range registries {
		if _, isBuiltin := registry.GetBuiltinRegistry(reg.Name); !isBuiltin {
			customRegs = append(customRegs, reg)
		}
	}

	// 如果没有自定义 registry
	if len(customRegs) == 0 {
		return fmt.Errorf("\n⚠️  No custom registries found")
	}

	// 如果未使用 force 参数，需要确认
	if !force {
		fmt.Printf("\nFound %d custom registries. Are you sure to remove them all? [y/N]: ", len(customRegs))
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil {
			answer = "n"
		}
		if !strings.EqualFold(answer, "y") {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// 执行删除
	var removed []string
	for _, reg := range customRegs {
		if err := manager.Remove(reg.Name); err != nil {
			fmt.Printf("\n❌  Failed to remove registry '%s': %v\n", reg.Name, err)
			continue
		}
		removed = append(removed, reg.Name)
	}

	// 显示结果
	if len(removed) > 0 {
		fmt.Printf("\n✨ Successfully removed %d registries: %s\n",
			len(removed),
			style.Success.Sprint(strings.Join(removed, ", ")))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// 添加命令行参数
	flags := rmCmd.Flags()
	flags.BoolVarP(&forceRemove, "force", "f", false, "Force remove without confirmation")
	flags.BoolVarP(&removeAll, "all", "a", false, "Remove all custom registries")
}
