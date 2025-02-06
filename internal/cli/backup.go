package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nrmgo/internal/backup"
	"nrmgo/internal/checker"
	"nrmgo/internal/style"

	"github.com/spf13/cobra"
)

// backupCmd 备份命令
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup package manager configurations",
	Long:  `Backup package manager configurations to backups/{timestamp} directory`,
	Run:   runBackup,
}

func runBackup(cmd *cobra.Command, args []string) {
	// 获取程序所在目录
	execPath, err := os.Executable()
	if err != nil {
		style.Error.Printf("❌ Failed to get executable path: %v\n", err)
		return
	}
	execPath = filepath.Dir(execPath)

	// 检查是否需要清理
	if days, _ := cmd.Flags().GetInt("clean"); days > 0 {
		cleaner := backup.NewCleaner(execPath)
		removed, err := cleaner.Clean(days)
		if err != nil {
			style.Error.Printf("❌ Failed to clean up: %v\n", err)
			return
		}
		style.Success.Printf("🧹 Successfully cleaned up %d directories older than %d days\n", removed, days)
		return
	}

	// 获取需要备份的包管理器
	var managers []string
	if all, _ := cmd.Flags().GetBool("all"); all {
		managers = append(managers, []string{"npm", "yarn", "pnpm", "bun"}...)
	} else {
		if npm, _ := cmd.Flags().GetBool("npm"); npm {
			managers = append(managers, "npm")
		}
		if yarn, _ := cmd.Flags().GetBool("yarn"); yarn {
			managers = append(managers, "yarn")
		}
		if pnpm, _ := cmd.Flags().GetBool("pnpm"); pnpm {
			managers = append(managers, "pnpm")
		}
		if bun, _ := cmd.Flags().GetBool("bun"); bun {
			managers = append(managers, "bun")
		}
	}

	// 如果没有指定任何包管理器，默认备份所有
	if len(managers) == 0 {
		managers = append(managers, []string{"npm", "yarn", "pnpm", "bun"}...)
	}

	// 获取已安装的包管理器
	installedManagers := make(map[string]bool)
	for _, pm := range checker.GetAvailableManagers() {
		installedManagers[pm.Name] = true
	}

	// 创建备份管理器
	bm := backup.NewBackupManager(execPath)

	// 执行备份
	results, err := bm.Backup(managers)
	if err != nil {
		style.Error.Printf("❌ Failed to backup: %v\n", err)
		return
	}

	// 显示备份目录
	if len(results) > 0 {
		backupPath := filepath.Dir(results[0].BackupPath)
		relPath := strings.TrimPrefix(backupPath, execPath)
		relPath = strings.TrimPrefix(relPath, string(os.PathSeparator))
		fmt.Println()
		style.Info.Printf("📂  Backup directory: %s\n", relPath)
	}

	// 创建结果映射
	resultMap := make(map[string]backup.BackupResult)
	for _, result := range results {
		resultMap[result.Manager] = result
	}

	// 收集不同状态的包管理器
	var (
		successList      []string
		notInstalledList []string
		failedList       []string
		failedDetails    = make(map[string]string)
	)

	for _, name := range managers {
		if !installedManagers[name] {
			notInstalledList = append(notInstalledList, name)
			continue
		}

		result, exists := resultMap[name]
		if !exists || !result.Success {
			failedList = append(failedList, name)
			if exists {
				failedDetails[name] = result.Error.Error()
			}
		} else {
			successList = append(successList, name)
		}
	}

	// 显示成功信息
	if len(successList) > 0 {
		style.Success.Printf("\n🎉  Successfully backed up: %s\n", strings.Join(successList, ", "))
	}

	// 显示未安装信息
	if len(notInstalledList) > 0 {
		style.Warning.Printf("\n⚠️   Not installed: %s\n", strings.Join(notInstalledList, ", "))
	}

	// 显示失败信息
	for _, name := range failedList {
		if errMsg, ok := failedDetails[name]; ok {
			style.Error.Printf("\n❌  Backup failed: %s (error: %s)\n", name, errMsg)
		} else {
			style.Error.Printf("\n❌  Backup failed: %s\n", name)
		}
	}
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// 添加命令行参数
	backupCmd.Flags().BoolP("all", "a", false, "Backup all package manager configurations")
	backupCmd.Flags().Bool("npm", false, "Backup npm configuration")
	backupCmd.Flags().Bool("yarn", false, "Backup yarn configuration")
	backupCmd.Flags().Bool("pnpm", false, "Backup pnpm configuration")
	backupCmd.Flags().Bool("bun", false, "Backup bun configuration")
	backupCmd.Flags().Int("clean", 0, "Clean up backups older than specified days")
}
