package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cleaner 备份清理器
type Cleaner struct {
	ExecPath string // 程序所在路径
}

// NewCleaner 创建清理器
func NewCleaner(execPath string) *Cleaner {
	return &Cleaner{
		ExecPath: execPath,
	}
}

// Clean 清理指定天数前的备份
func (c *Cleaner) Clean(days int) (int, error) {
	// 获取当前时间
	now := time.Now()

	// 获取备份根目录
	backupsRoot := filepath.Join(c.ExecPath, "backups")

	// 如果备份目录不存在，直接返回
	if _, err := os.Stat(backupsRoot); os.IsNotExist(err) {
		return 0, nil
	}

	// 遍历备份目录
	entries, err := os.ReadDir(backupsRoot)
	if err != nil {
		return 0, fmt.Errorf("failed to read backups directory: %w", err)
	}

	// 记录删除的文件数
	removed := 0

	// 处理每个备份目录
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 解析时间戳
		backupTime, err := time.ParseInLocation("20060102_150405", entry.Name(), time.Local)
		if err != nil {
			continue // 跳过无法解析时间戳的目录
		}

		// 检查是否超过指定天数
		if now.Sub(backupTime).Hours() > float64(days*24) {
			// 删除目录
			path := filepath.Join(backupsRoot, entry.Name())
			if err := os.RemoveAll(path); err != nil {
				return removed, fmt.Errorf("failed to remove directory %s: %w", path, err)
			}
			removed++
		}
	}

	return removed, nil
}
