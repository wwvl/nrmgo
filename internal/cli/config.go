package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"nrmgo/internal/config"
	"nrmgo/internal/style"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage nrmgo configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 如果没有子命令，显示帮助
		return cmd.Help()
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取默认配置模板
		template, err := config.GetDefaultTemplate()
		if err != nil {
			return fmt.Errorf("❌  Failed to get default template: %v", err)
		}

		// 获取程序所在目录
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("❌  Failed to get executable path: %v", err)
		}
		configPath := config.GetConfigPath(execPath)

		// 写入配置文件
		if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
			return fmt.Errorf("❌  Failed to write config file: %v", err)
		}

		fmt.Print(style.Success.Sprintf("\n✨  Configuration file initialized successfully at %s\n", configPath))
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置
		cfg, err := config.LoadConfig()
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("⚠️  Config file not found, please use 'nrmgo config init' to create it ~")
			}
			return fmt.Errorf("❌  Failed to load config: %v", err)
		}

		// 使用树形渲染器显示配置
		if err := config.RenderConfig(cfg); err != nil {
			return fmt.Errorf("❌  Failed to render config: %v", err)
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
}
