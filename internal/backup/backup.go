package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Backup 执行备份
func (bm *BackupManager) Backup(managers []string) ([]BackupResult, error) {
	// 创建 backups 根目录
	backupsRoot := filepath.Join(bm.ExecPath, "backups")
	if err := os.MkdirAll(backupsRoot, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backups directory: %w", err)
	}

	// 生成备份目录
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(backupsRoot, timestamp)

	// 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	results := make([]BackupResult, 0)

	// 如果没有指定包管理器，则备份所有
	if len(managers) == 0 {
		managers = bm.GetAllManagers()
	}

	// 执行备份
	for _, name := range managers {
		manager := bm.GetManager(name)
		if manager == nil {
			continue
		}

		// 查找配置文件
		var sourcePath string
		for _, path := range manager.Paths {
			if _, err := os.Stat(path); err == nil {
				sourcePath = path
				break
			}
		}

		result := BackupResult{
			Manager:    manager.Name,
			SourcePath: sourcePath,
		}

		// 如果找不到配置文件
		if sourcePath == "" {
			result.Success = false
			result.Error = fmt.Errorf("config file not found")
			results = append(results, result)
			continue
		}

		// 设置备份路径
		backupPath := filepath.Join(backupDir, filepath.Base(sourcePath))
		result.BackupPath = backupPath

		// 复制文件
		if err := copyFile(sourcePath, backupPath); err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to copy file: %w", err)
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results, nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 获取源文件权限
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 创建目标文件
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
