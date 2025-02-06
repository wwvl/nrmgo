package registry

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"nrmgo/internal/config"
)

// Info 表示一个 npm registry 的完整信息
type Info struct {
	Name        string `toml:"name"`        // registry 名称
	URL         string `toml:"url"`         // registry URL
	Home        string `toml:"home"`        // 主页地址
	Description string `toml:"description"` // 描述信息
}

// FromConfig 从配置创建 Info
func FromConfig(name string, cfg *config.Registry) *Info {
	return &Info{
		Name:        name,
		URL:         cfg.URL,
		Home:        cfg.Home,
		Description: cfg.Description,
	}
}

// ToConfig 转换为配置
func (r *Info) ToConfig() *config.Registry {
	return &config.Registry{
		URL:         r.URL,
		Home:        r.Home,
		Description: r.Description,
	}
}

// TestResult 表示 registry 的测试结果
type TestResult struct {
	Name     string        // registry 名称
	URL      string        // registry URL
	IsOnline bool          // 是否在线
	Latency  time.Duration // 延迟时间
	Error    string        // 错误信息
}

// Manager 定义 registry 管理器接口
type Manager interface {
	// List 列出所有可用的 registry
	List() []*Info

	// Get 获取指定名称的 registry
	Get(name string) (*Info, bool)

	// Add 添加自定义 registry
	Add(name string, reg *Info) error

	// Remove 移除自定义 registry
	Remove(name string) error

	// Use 切换当前使用的 registry
	Use(name string) error

	// Current 获取当前使用的 registry
	Current() (*Info, error)

	// Test 测试指定 registry 的延迟
	Test(names ...string) []*TestResult

	// Rename 重命名 registry
	// 规则：
	// 1. oldName 必须存在且不能是内置 registry
	// 2. newName 不能是内置 registry 名称且不能已存在
	// 3. 如果重命名的是当前使用的 registry，会自动更新
	Rename(oldName, newName string) error
}

// NewRegistry 创建一个新的 Registry
func NewRegistry(name, url, home, description string) *Info {
	return &Info{
		Name:        name,
		URL:         url,
		Home:        home,
		Description: description,
	}
}

// IsValidName 验证 registry name 是否合法
// 规则：只允许字母、数字和下划线
func IsValidName(name string) error {
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		return fmt.Errorf("❌ Invalid registry name: %s (only letters, numbers and underscores are allowed)", name)
	}
	return nil
}

// IsValidURL 验证 registry URL 是否合法
// 规则：
// 1. 必须是合法的 URL 格式
// 2. 必须以 http:// 或 https:// 开头
func IsValidURL(registryURL string) error {
	parsedURL, err := url.Parse(registryURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("❌ Invalid registry URL: %s", registryURL)
	}
	if !regexp.MustCompile(`^https?://`).MatchString(registryURL) {
		return fmt.Errorf("❌ Invalid registry URL: %s (must start with http:// or https://)", registryURL)
	}
	return nil
}
