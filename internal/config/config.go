package config

import (
	_ "embed" // 用于嵌入模板文件
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	configFileName = "config.toml"
)

//go:embed config.toml
var defaultTemplate string

// GetConfigPath 获取配置文件路径
func GetConfigPath(execPath string) string {
	return filepath.Join(filepath.Dir(execPath), configFileName)
}

// GetDefaultTemplate 获取默认的配置文件内容
func GetDefaultTemplate() (string, error) {
	return defaultTemplate, nil
}

// Config 配置结构
type Config struct {
	// CustomRegistries 用户自定义的 registry 列表
	CustomRegistries map[string]*Registry `toml:"custom_registries"`

	// MaxConcurrentRequests HTTP 并发请求数，用于延迟测试
	// 默认值：5，建议范围：1-10
	MaxConcurrentRequests int `toml:"max_concurrent_requests"`
}

// Registry 注册表信息
type Registry struct {
	// URL registry 的地址
	URL string `toml:"url"`

	// Home registry 的主页地址（可选）
	Home string `toml:"home,omitempty"`

	// Description registry 的描述信息（可选）
	Description string `toml:"description,omitempty"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 获取程序所在目录
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	configPath := GetConfigPath(execPath)

	// 如果配置文件不存在，返回错误
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	cfg := &Config{
		CustomRegistries: make(map[string]*Registry),
	}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 验证配置
	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// SaveConfig 保存配置
func SaveConfig(cfg *Config) error {
	// 验证配置
	if err := ValidateConfig(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 获取程序所在目录
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	configPath := GetConfigPath(execPath)

	// 使用 toml.Marshal 序列化配置
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
