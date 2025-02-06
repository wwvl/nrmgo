package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version 是应用程序的语义化版本号
	Version = "dev"

	// Commit 是构建时的 Git 提交哈希
	Commit = "none"

	// BuildTime 是构建时的时间戳
	BuildTime = "unknown"

	// GoVersion 是构建时使用的 Go 版本
	GoVersion = runtime.Version()
)

// VersionInfo 返回格式化的版本信息
func VersionInfo() string {
	return fmt.Sprintf(`NRMGO Version Information:
  Version:    %s
  Commit:     %s
  Built:      %s
  Go version: %s
  OS/Arch:    %s/%s`,
		Version,
		Commit,
		BuildTime,
		GoVersion,
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// ShortVersion 返回简短的版本信息
func ShortVersion() string {
	return fmt.Sprintf("nrmgo version %s", Version)
}

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nrmgo version %s\n", Version)
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}
