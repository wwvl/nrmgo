package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"nrmgo/internal/checker"
	"nrmgo/internal/style"
)

// unuseCmd 恢复默认 registry 命令
var unuseCmd = &cobra.Command{
	Use:   "unuse",
	Short: "Restore package managers to their default registries",
	Long: `Restore package managers to their default registries.

Default registries:
  npm/pnpm/bun:  https://registry.npmjs.org/
  yarn: https://registry.yarnpkg.com/

Examples:
  # Restore all package managers to their default registries
  nrmgo unuse

  # Restore specific package managers
  nrmgo unuse --npm --pnpm

  # Same as 'nrmgo unuse'
  nrmgo unuse --all`,
	RunE: runUnuse,
}

func runUnuse(cmd *cobra.Command, args []string) error {
	// 获取需要恢复的包管理器列表
	var managers []string
	if all, _ := cmd.Flags().GetBool("all"); all {
		managers = append(managers, []string{"npm", "yarn", "pnpm", "bun"}...)
	} else {
		hasFlag := false
		if npm, _ := cmd.Flags().GetBool("npm"); npm {
			managers = append(managers, "npm")
			hasFlag = true
		}
		if yarn, _ := cmd.Flags().GetBool("yarn"); yarn {
			managers = append(managers, "yarn")
			hasFlag = true
		}
		if pnpm, _ := cmd.Flags().GetBool("pnpm"); pnpm {
			managers = append(managers, "pnpm")
			hasFlag = true
		}
		if bun, _ := cmd.Flags().GetBool("bun"); bun {
			managers = append(managers, "bun")
			hasFlag = true
		}

		// 如果没有指定任何参数，默认恢复所有
		if !hasFlag {
			managers = append(managers, []string{"npm", "yarn", "pnpm", "bun"}...)
		}
	}

	// 检测已安装的包管理器
	installedPMs := checker.DetectPackageManagers()
	installedMap := make(map[string]bool)
	for _, pm := range installedPMs {
		installedMap[pm.Name] = pm.Installed
	}

	// 收集不同状态的包管理器
	var (
		successList      []string
		notInstalledList []string
		failedList       []string
		failedDetails    = make(map[string]string)
	)

	// 执行恢复操作
	for _, name := range managers {
		if !installedMap[name] {
			notInstalledList = append(notInstalledList, name)
			continue
		}

		// 获取默认 registry
		defaultRegistry, _, _, err := checker.GetDefaultRegistry(name)
		if err != nil {
			failedList = append(failedList, name)
			failedDetails[name] = err.Error()
			continue
		}

		if err := checker.SetRegistry(name, defaultRegistry); err != nil {
			failedList = append(failedList, name)
			failedDetails[name] = err.Error()
		} else {
			successList = append(successList, name)
		}
	}

	// 显示成功信息
	if len(successList) > 0 {
		fmt.Println()
		style.Success.Printf("✅  Successfully restored: %s\n", strings.Join(successList, ", "))
	}

	// 显示未安装信息
	if len(notInstalledList) > 0 {
		fmt.Println()
		style.Warning.Printf("⚠️   Not installed: %s\n", strings.Join(notInstalledList, ", "))
	}

	// 显示失败信息
	for _, name := range failedList {
		fmt.Println()
		if errMsg, ok := failedDetails[name]; ok {
			style.Error.Printf("❌  Restore failed: %s (error: %s)\n", name, errMsg)
		} else {
			style.Error.Printf("❌  Restore failed: %s\n", name)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(unuseCmd)

	// 添加命令行参数
	unuseCmd.Flags().BoolP("all", "a", false, "Restore all package managers to their default registries")
	unuseCmd.Flags().Bool("npm", false, "Restore npm to its default registry (https://registry.npmjs.org/)")
	unuseCmd.Flags().Bool("yarn", false, "Restore yarn to its default registry (https://registry.yarnpkg.com/)")
	unuseCmd.Flags().Bool("pnpm", false, "Restore pnpm to its default registry (https://registry.npmjs.org/)")
	unuseCmd.Flags().Bool("bun", false, "Restore bun to its default registry (https://registry.npmjs.org/)")
}
