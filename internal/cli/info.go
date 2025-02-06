package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"nrmgo/internal/checker"
	"nrmgo/internal/style"
	"nrmgo/internal/table"
)

// infoCmd 显示包管理器信息
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show package manager information",
	Long:  "Show package manager information, including installation status, version, registry and config file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 检测包管理器
		managers := checker.DetectPackageManagers()

		// 创建表格渲染器
		renderer := table.NewTableRenderer([]string{
			"Package Manager",
			"Status",
			"Version",
			"Registry",
			"Config",
		})

		// 添加数据行
		for _, pm := range managers {
			// 状态显示
			status := style.Error.Sprint("❌")
			if pm.Installed {
				status = style.Success.Sprint("✅")
			}

			// 版本显示
			version := "-"
			if pm.Installed {
				version = pm.Version
			}

			// Registry 显示
			registry := "-"
			if pm.Installed {
				if reg, _, _, err := checker.GetRegistry(pm.Name); err == nil {
					registry = reg
				}
			}

			// Config 显示
			config := "-"
			if pm.Installed {
				home, _ := os.UserHomeDir()
				if pm.ConfigPath != "" {
					config = strings.Replace(pm.ConfigPath, home, "$HOME", 1)
					if !pm.ConfigExists {
						config = fmt.Sprintf("%s %s", config, style.Error.Sprint("(missing)"))
					}
				}
			}

			renderer.MustAddRow([]string{
				pm.Name,
				status,
				version,
				registry,
				config,
			})
		}

		// 输出表格
		fmt.Println()
		if err := renderer.Render(); err != nil {
			return fmt.Errorf("❌  Failed to render table: %v", err)
		}

		// 显示配置文件提示
		var configMap = make(map[string][]string)
		for _, pm := range managers {
			if pm.Installed && !pm.ConfigExists {
				configMap[pm.ConfigPath] = append(configMap[pm.ConfigPath], pm.Name)
			}
		}

		if len(configMap) > 0 {
			fmt.Printf("\n%s Missing package manager configuration files. You can set them up using '%s':\n",
				"💡",
				style.Success.Sprint("nrmgo use <registry>"))
			home, _ := os.UserHomeDir()
			for path, pms := range configMap {
				displayPath := strings.Replace(path, home, "$HOME", 1)
				fmt.Printf("  📄 %s: %s\n",
					strings.Join(pms, "/"),
					displayPath)
			}
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
