package checker

import (
	"bytes"
	"os/exec"
	"strings"
	"sync"
)

var managers = []string{"npm", "yarn", "pnpm", "bun"}

// execCommand 执行命令并返回输出结果
func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", &CommandError{
			Command: name,
			Args:    strings.Join(args, " "),
			Err:     err,
		}
	}

	return strings.TrimSpace(out.String()), nil
}

// getVersion 获取包管理器版本
func getVersion(name string) (string, error) {
	return execCommand(name, "-v")
}

// detectSingle 检测单个包管理器
func detectSingle(name string) PackageManager {
	pm := PackageManager{Name: name}

	// 获取版本号
	version, err := getVersion(name)
	if err != nil {
		pm.Error = err
		return pm
	}

	pm.Installed = true
	pm.Version = version

	// 获取 registry 和配置文件信息
	if registry, configPath, exists, err := GetRegistry(name); err == nil {
		pm.Registry = registry
		pm.ConfigPath = configPath
		pm.ConfigExists = exists
	}

	return pm
}

// DetectPackageManagers 并行检测所有包管理器
func DetectPackageManagers() []PackageManager {
	var wg sync.WaitGroup
	results := make([]PackageManager, len(managers))

	for i, name := range managers {
		wg.Add(1)
		go func(index int, pmName string) {
			defer wg.Done()
			results[index] = detectSingle(pmName)
		}(i, name)
	}

	wg.Wait()
	return results
}

// GetManager 获取指定包管理器的信息
func GetManager(name string) (PackageManager, bool) {
	results := DetectPackageManagers()
	for _, pm := range results {
		if pm.Name == name {
			return pm, true
		}
	}
	return PackageManager{}, false
}

// GetAvailableManagers 获取所有可用的包管理器
func GetAvailableManagers() []PackageManager {
	results := DetectPackageManagers()
	var available []PackageManager
	for _, pm := range results {
		if pm.Installed {
			available = append(available, pm)
		}
	}
	return available
}

// IsAvailable 检查指定包管理器是否可用
func IsAvailable(name string) bool {
	if pm, found := GetManager(name); found {
		return pm.Installed
	}
	return false
}
