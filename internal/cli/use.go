package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"nrmgo/internal/checker"
	"nrmgo/internal/registry"
	"nrmgo/internal/style"
	"nrmgo/internal/table"
)

// useCmd 切换 registry
var useCmd = &cobra.Command{
	Use:   "use [registry]",
	Short: "Switch registry for package managers",
	Long: `Switch registry for package managers. If no registry is specified, 
it will automatically test and select the fastest registry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置并创建管理器
		_, manager, err := loadConfigAndCreateManager()
		if err != nil {
			return err
		}

		// 检测已安装的包管理器
		installedPMs := checker.DetectPackageManagers()
		if len(installedPMs) == 0 {
			return fmt.Errorf("\n❌  No package manager installed")
		}

		// 获取所有 registry
		registries := manager.List()
		if len(registries) == 0 {
			return fmt.Errorf("\n❌  No registry found")
		}

		// 如果指定了 registry 名称
		if len(args) > 0 {
			registryName := args[0]
			reg, ok := manager.Get(registryName)
			if !ok {
				return fmt.Errorf("\n❌  Registry '%s' not found", registryName)
			}

			// 设置为当前使用的 registry
			if err := manager.Use(reg.Name); err != nil {
				return fmt.Errorf("\n❌  Failed to set registry: %v", err)
			}

			// 输出成功信息
			installedNames := []string{}
			for _, pm := range installedPMs {
				if pm.Installed {
					installedNames = append(installedNames, pm.Name)
				}
			}
			fmt.Printf("\n✨ Successfully Changed Package Manager(%s) to: %s\n",
				strings.Join(installedNames, ", "),
				style.Success.Sprint(reg.Name))
			return nil
		}

		// 自动测试并选择最快的 registry
		fmt.Println()
		progressbar, err := pterm.DefaultProgressbar.
			WithTotal(len(registries)).
			WithTitle("Testing Registry Latency").
			WithShowCount(true).
			WithElapsedTimeRoundingFactor(time.Millisecond * 10).
			WithMaxWidth(66).
			Start()
		if err != nil {
			return fmt.Errorf("\n❌  Failed to create progress bar: %v", err)
		}

		// 测试每个 registry 的延迟
		testResults := manager.Test()

		// 按延迟排序
		sort.Slice(testResults, func(i, j int) bool {
			// 如果有错误，排在后面
			if testResults[i].Error != "" && testResults[j].Error == "" {
				return false
			}
			if testResults[i].Error == "" && testResults[j].Error != "" {
				return true
			}
			return testResults[i].Latency < testResults[j].Latency
		})

		// 创建表格渲染器
		renderer := table.NewTableRenderer([]string{
			"Name",
			"Registry URL",
			"Latency",
		})

		// 添加数据行
		var fastestReg *registry.Info
		for _, result := range testResults {
			latency := "-"
			if result.Error == "" {
				latency = fmt.Sprintf("%dms", result.Latency.Milliseconds())
				// 记录第一个成功的 registry 为最快的
				if fastestReg == nil {
					reg, _ := manager.Get(result.Name)
					fastestReg = reg
				}
			}

			renderer.MustAddRow([]string{
				result.Name,
				result.URL,
				latency,
			})
			progressbar.Increment()
		}

		// 渲染表格
		fmt.Println()
		if err := renderer.Render(); err != nil {
			return fmt.Errorf("\n❌  Failed to render table: %v", err)
		}
		fmt.Println()

		// 如果没有可用的 registry
		if fastestReg == nil {
			return fmt.Errorf("\n❌  No available registry found")
		}

		// 设置最快的 registry 为当前使用的 registry
		if err := manager.Use(fastestReg.Name); err != nil {
			return fmt.Errorf("\n❌  Failed to set registry: %v", err)
		}

		// 输出成功信息
		installedNames := []string{}
		for _, pm := range installedPMs {
			if pm.Installed {
				installedNames = append(installedNames, pm.Name)
			}
		}
		fmt.Printf("✨ Successfully Changed Package Manager(%s) to: %s\n",
			strings.Join(installedNames, ", "),
			style.Success.Sprint(fastestReg.Name))

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(useCmd)
}
