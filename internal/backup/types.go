package backup

// Manager 包管理器配置
type Manager struct {
	Name       string   // 包管理器名称
	ConfigFile string   // 配置文件名
	Paths      []string // 可能的配置文件路径
}

// BackupManager 备份管理器
type BackupManager struct {
	ExecPath string              // 程序所在路径
	Managers map[string]*Manager // 包管理器配置
}

// BackupResult 备份结果
type BackupResult struct {
	Manager    string // 包管理器名称
	SourcePath string // 源文件路径
	BackupPath string // 备份路径
	Success    bool   // 是否成功
	Error      error  // 错误信息
}
