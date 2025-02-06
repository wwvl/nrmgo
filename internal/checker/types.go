package checker

// PackageManager 定义包管理器信息
type PackageManager struct {
	Name         string // 包管理器名称
	Installed    bool   // 是否已安装
	Version      string // 版本号
	Registry     string // registry 配置
	Error        error  // 错误信息
	ConfigExists bool   // 配置文件是否存在
	ConfigPath   string // 配置文件路径
}

// RegistryConfig 定义包管理器的 registry 配置
type RegistryConfig struct {
	Name         string                       // 包管理器名称
	ConfigFile   string                       // 配置文件名
	DefaultValue string                       // 默认 registry
	Parser       func([]byte) (string, error) // 配置文件解析函数
	Writer       func([]byte, string) []byte  // 配置文件写入函数
}

// CommandError 定义命令执行错误
type CommandError struct {
	Command string // 执行的命令
	Args    string // 命令参数
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *CommandError) Error() string {
	return e.Err.Error()
}

// NPMConfig 定义 npm 配置结构
type NPMConfig struct {
	Registry   string // registry 地址
	AlwaysAuth bool   // 是否总是进行身份验证
	StrictSSL  bool   // 是否启用严格的 SSL
}

// DefaultNPMConfig 返回默认的 npm 配置
func DefaultNPMConfig() NPMConfig {
	return NPMConfig{
		AlwaysAuth: false,
		StrictSSL:  true,
	}
}
