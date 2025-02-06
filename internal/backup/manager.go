package backup

import (
	"path/filepath"

	"nrmgo/internal/checker"
)

// NewBackupManager 创建备份管理器
func NewBackupManager(execPath string) *BackupManager {
	// 初始化管理器
	bm := &BackupManager{
		ExecPath: execPath,
		Managers: make(map[string]*Manager),
	}

	// 获取所有可用的包管理器
	pms := checker.GetAvailableManagers()

	// 根据检查器初始化包管理器配置
	for _, pm := range pms {
		bm.Managers[pm.Name] = &Manager{
			Name:       pm.Name,
			ConfigFile: filepath.Base(pm.ConfigPath),
			Paths:      []string{pm.ConfigPath},
		}
	}

	return bm
}

// GetAllManagers 获取所有包管理器名称
func (bm *BackupManager) GetAllManagers() []string {
	managers := make([]string, 0, len(bm.Managers))
	for name := range bm.Managers {
		managers = append(managers, name)
	}
	return managers
}

// GetManager 获取指定包管理器配置
func (bm *BackupManager) GetManager(name string) *Manager {
	return bm.Managers[name]
}
