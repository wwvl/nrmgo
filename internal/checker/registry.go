package checker

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getConfigPath 获取配置文件的完整路径
func getConfigPath(filename string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, filename), nil
}

// readConfigFile 读取配置文件内容
func readConfigFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return data, err
}

// writeConfigFile 写入配置文件
func writeConfigFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// parseNPMStyleConfig 解析 npm 风格的配置文件
func parseNPMStyleConfig(data []byte) (string, error) {
	config := DefaultNPMConfig() // 使用默认配置

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析配置项
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "registry" {
			config.Registry = value
		}
	}

	return config.Registry, nil
}

// parseYarnConfig 解析 yarn 配置文件
func parseYarnConfig(data []byte) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "registry") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return strings.Trim(parts[1], "\"'"), nil
			}
		}
	}
	return "", nil
}

// parseBunConfig 解析 bun 配置文件
func parseBunConfig(data []byte) (string, error) {
	var inInstallSection bool
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "[install]" {
			inInstallSection = true
			continue
		} else if strings.HasPrefix(line, "[") {
			inInstallSection = false
			continue
		}
		if inInstallSection && strings.HasPrefix(line, "registry") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.Trim(strings.TrimSpace(parts[1]), "\"'"), nil
			}
		}
	}
	return "", nil
}

// writeNPMStyleConfig 写入 npm 风格的配置
func writeNPMStyleConfig(data []byte, registry string) []byte {
	var newLines []string
	var found struct {
		alwaysAuth bool
		strictSSL  bool
	}

	// 获取默认配置
	defaultConfig := DefaultNPMConfig()

	// 如果有现有内容，保留其他配置
	if len(data) > 0 {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			trimmedLine := strings.TrimSpace(line)

			// 跳过空行和注释
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
				newLines = append(newLines, line)
				continue
			}

			switch {
			case strings.HasPrefix(trimmedLine, "registry="):
				// 跳过已存在的 registry 配置
				continue
			case strings.HasPrefix(trimmedLine, "always-auth="):
				found.alwaysAuth = true
			case strings.HasPrefix(trimmedLine, "strict-ssl="):
				found.strictSSL = true
			default:
				newLines = append(newLines, line)
			}
		}
	}

	// 按固定顺序写入配置项
	if len(newLines) > 0 {
		newLines = append(newLines, "")
	}

	// 1. registry - 必须在最前面
	newLines = append(newLines, fmt.Sprintf("registry=%s", registry))

	// 2. always-auth
	if !found.alwaysAuth {
		newLines = append(newLines, fmt.Sprintf("always-auth=%v", defaultConfig.AlwaysAuth))
	}

	// 3. strict-ssl
	if !found.strictSSL {
		newLines = append(newLines, fmt.Sprintf("strict-ssl=%v", defaultConfig.StrictSSL))
	}

	// 确保文件末尾有一个空行
	newLines = append(newLines, "")

	return []byte(strings.Join(newLines, "\n"))
}

// writeYarnConfig 写入 yarn 配置
func writeYarnConfig(data []byte, registry string) []byte {
	var newLines []string
	registryFound := false

	// 如果有现有内容，保留其他配置
	if len(data) > 0 {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.TrimSpace(line), "registry ") {
				newLines = append(newLines, fmt.Sprintf("registry \"%s\"", registry))
				registryFound = true
			} else {
				newLines = append(newLines, line)
			}
		}
	}

	// 如果没有找到 registry 配置，添加新的
	if !registryFound {
		if len(newLines) > 0 {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, fmt.Sprintf("registry \"%s\"", registry))
	}

	return []byte(strings.Join(newLines, "\n") + "\n")
}

// writeBunConfig 写入 bun 配置
func writeBunConfig(data []byte, registry string) []byte {
	var newLines []string
	var inInstallSection bool
	registryFound := false
	installSectionFound := false

	// 如果有现有内容，保留其他配置
	if len(data) > 0 {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if line == "[install]" {
				inInstallSection = true
				installSectionFound = true
			} else if strings.HasPrefix(line, "[") {
				inInstallSection = false
			}

			if inInstallSection && strings.HasPrefix(strings.TrimSpace(line), "registry") {
				newLines = append(newLines, fmt.Sprintf("registry = \"%s\"", registry))
				registryFound = true
			} else {
				newLines = append(newLines, line)
			}
		}
	}

	// 如果没有找到 registry 配置
	if !registryFound {
		if !installSectionFound {
			if len(newLines) > 0 {
				newLines = append(newLines, "")
			}
			newLines = append(newLines, "[install]")
		}
		newLines = append(newLines, fmt.Sprintf("registry = \"%s\"", registry))
	}

	return []byte(strings.Join(newLines, "\n") + "\n")
}

// registryConfigs 定义支持的包管理器配置
var registryConfigs = map[string]RegistryConfig{
	"npm": {
		Name:         "npm",
		ConfigFile:   ".npmrc",
		DefaultValue: "https://registry.npmjs.org/",
		Parser:       parseNPMStyleConfig,
		Writer:       writeNPMStyleConfig,
	},
	"yarn": {
		Name:         "yarn",
		ConfigFile:   ".yarnrc",
		DefaultValue: "https://registry.yarnpkg.com/",
		Parser:       parseYarnConfig,
		Writer:       writeYarnConfig,
	},
	"bun": {
		Name:         "bun",
		ConfigFile:   ".bunfig.toml",
		DefaultValue: "https://registry.npmjs.org/",
		Parser:       parseBunConfig,
		Writer:       writeBunConfig,
	},
}

// GetRegistry 获取包管理器的 registry 配置
func GetRegistry(name string) (registry string, configPath string, exists bool, err error) {
	// 如果是 pnpm，直接使用 npm 的配置
	if name == "pnpm" {
		return GetRegistry("npm")
	}

	config, ok := registryConfigs[name]
	if !ok {
		return "", "", false, fmt.Errorf("unsupported package manager: %s", name)
	}

	// 获取配置文件路径
	configPath, err = getConfigPath(config.ConfigFile)
	if err != nil {
		return config.DefaultValue, configPath, false, nil
	}

	// 读取配置文件
	data, err := readConfigFile(configPath)
	if err != nil {
		return config.DefaultValue, configPath, false, nil
	}

	// 如果文件不存在
	if data == nil {
		return config.DefaultValue, configPath, false, nil
	}

	// 解析配置文件
	registry, err = config.Parser(data)
	if err != nil || registry == "" {
		return config.DefaultValue, configPath, true, nil
	}

	return registry, configPath, true, nil
}

// SetRegistry 设置包管理器的 registry
func SetRegistry(name, registry string) error {
	if name == "pnpm" {
		name = "npm" // pnpm 使用 npm 的配置
	}

	config, exists := registryConfigs[name]
	if !exists {
		return fmt.Errorf("unsupported package manager: %s", name)
	}

	// 获取配置文件路径
	configPath, err := getConfigPath(config.ConfigFile)
	if err != nil {
		return err
	}

	// 读取现有配置
	data, err := readConfigFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 生成新的配置内容
	newData := config.Writer(data, registry)

	// 写入配置文件
	return writeConfigFile(configPath, newData)
}

// GetDefaultRegistry 获取包管理器的默认 registry 配置
func GetDefaultRegistry(name string) (registry string, configPath string, exists bool, err error) {
	// 如果是 pnpm，直接使用 npm 的配置
	if name == "pnpm" {
		name = "npm"
	}

	config, ok := registryConfigs[name]
	if !ok {
		return "", "", false, fmt.Errorf("unsupported package manager: %s", name)
	}

	// 获取配置文件路径
	configPath, err = getConfigPath(config.ConfigFile)
	if err != nil {
		return config.DefaultValue, configPath, false, nil
	}

	return config.DefaultValue, configPath, true, nil
}
