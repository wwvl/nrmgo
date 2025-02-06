package registry

import (
	"context"
	"fmt"
	"sort"

	"nrmgo/internal/checker"
	"nrmgo/internal/config"
	"nrmgo/internal/latency"
)

// manager 实现了 Manager 接口
type manager struct {
	cfg *config.Config
}

// NewManager 创建一个新的 registry 管理器
func NewManager(cfg *config.Config) Manager {
	return &manager{cfg: cfg}
}

// List 列出所有可用的 registry
func (m *manager) List() []*Info {
	// 获取所有内置 registry
	builtinRegs := ListBuiltinRegistries()

	// 合并自定义 registry
	result := make([]*Info, 0, len(builtinRegs)+len(m.cfg.CustomRegistries))
	result = append(result, builtinRegs...)

	// 添加自定义 registry
	for name, reg := range m.cfg.CustomRegistries {
		result = append(result, FromConfig(name, reg))
	}

	// 按名称排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// Get 获取指定名称的 registry
func (m *manager) Get(name string) (*Info, bool) {
	// 先查找内置 registry
	if reg, ok := GetBuiltinRegistry(name); ok {
		return reg, true
	}

	// 再查找自定义 registry
	if reg, ok := m.cfg.CustomRegistries[name]; ok {
		return FromConfig(name, reg), true
	}

	return nil, false
}

// Add 添加自定义 registry
func (m *manager) Add(name string, reg *Info) error {
	// 检查是否是内置 registry
	if _, ok := GetBuiltinRegistry(name); ok {
		return fmt.Errorf("cannot add built-in registry: %s", name)
	}

	// 添加到自定义 registry
	m.cfg.CustomRegistries[name] = reg.ToConfig()

	// 保存配置
	return config.SaveConfig(m.cfg)
}

// Remove 移除自定义 registry
func (m *manager) Remove(name string) error {
	// 检查是否是内置 registry
	if _, ok := GetBuiltinRegistry(name); ok {
		return fmt.Errorf("cannot remove built-in registry: %s", name)
	}

	// 检查是否存在
	if _, ok := m.cfg.CustomRegistries[name]; !ok {
		return fmt.Errorf("registry not found: %s", name)
	}

	// 删除 registry
	delete(m.cfg.CustomRegistries, name)

	// 保存配置
	return config.SaveConfig(m.cfg)
}

// Use 切换当前使用的 registry
func (m *manager) Use(name string) error {
	// 检查 registry 是否存在
	reg, ok := m.Get(name)
	if !ok {
		return fmt.Errorf("registry not found: %s", name)
	}

	// 获取已安装的包管理器
	installedPMs := checker.DetectPackageManagers()

	// 为每个已安装的包管理器设置 registry
	for _, pm := range installedPMs {
		if pm.Installed {
			if err := checker.SetRegistry(pm.Name, reg.URL); err != nil {
				return fmt.Errorf("failed to set %s registry: %v", pm.Name, err)
			}
		}
	}

	return nil
}

// Current 获取当前使用的 registry
func (m *manager) Current() (*Info, error) {
	// 获取当前 npm registry
	registry, _, _, err := checker.GetRegistry("npm")
	if err != nil {
		return nil, fmt.Errorf("failed to get npm registry: %v", err)
	}

	// 查找匹配的 registry
	for _, reg := range m.List() {
		if reg.URL == registry {
			return reg, nil
		}
	}

	// 如果没有找到匹配的，返回一个临时的 Info
	return &Info{
		Name:        "unknown",
		URL:         registry,
		Description: "Current npm registry",
	}, nil
}

// Test 测试指定 registry 的延迟
func (m *manager) Test(names ...string) []*TestResult {
	var registries []*Info

	// 如果没有指定名称，测试所有 registry
	if len(names) == 0 {
		registries = m.List()
	} else {
		// 否则只测试指定的 registry
		for _, name := range names {
			if reg, ok := m.Get(name); ok {
				registries = append(registries, reg)
			}
		}
	}

	// 创建测试目标
	targets := make([]latency.Target, len(registries))
	for i, reg := range registries {
		targets[i] = latency.NewTarget(reg.Name, reg.URL).
			WithTestPath("package.json").
			WithHeaders(map[string][]string{
				"Accept": {"application/json"},
			})
	}

	// 创建测试器并执行测试
	tester := latency.NewTesterFromConfig(m.cfg)
	results := tester.Test(context.Background(), targets)

	// 转换结果
	testResults := make([]*TestResult, len(results))
	for i, result := range results {
		testResults[i] = &TestResult{
			Name:     result.Name,
			URL:      result.URL,
			IsOnline: result.IsOnline,
			Latency:  result.Latency,
			Error:    result.Error,
		}
	}

	return testResults
}

// Rename 重命名 registry
func (m *manager) Rename(oldName, newName string) error {
	// 检查 old registry 是否存在
	oldReg, exists := m.Get(oldName)
	if !exists {
		return fmt.Errorf("registry not found: %s", oldName)
	}

	// 检查是否为内置 registry
	if _, isBuiltin := GetBuiltinRegistry(oldName); isBuiltin {
		return fmt.Errorf("cannot rename built-in registry: %s", oldName)
	}

	// 检查新名称是否已存在
	if _, exists := m.Get(newName); exists {
		return fmt.Errorf("registry already exists: %s", newName)
	}

	// 检查新名称是否为内置 registry 名称
	if _, isBuiltin := GetBuiltinRegistry(newName); isBuiltin {
		return fmt.Errorf("cannot use built-in registry name: %s", newName)
	}

	// 创建新的 registry 信息
	newReg := &Info{
		Name:        newName,
		URL:         oldReg.URL,
		Home:        oldReg.Home,
		Description: oldReg.Description,
	}

	// 获取当前使用的 registry
	current, _ := m.Current()

	// 删除旧的 registry
	delete(m.cfg.CustomRegistries, oldName)

	// 添加新的 registry
	m.cfg.CustomRegistries[newName] = newReg.ToConfig()

	// 保存配置
	if err := config.SaveConfig(m.cfg); err != nil {
		// 如果保存失败，恢复原状态
		delete(m.cfg.CustomRegistries, newName)
		m.cfg.CustomRegistries[oldName] = oldReg.ToConfig()
		return fmt.Errorf("failed to save config: %v", err)
	}

	// 如果重命名的是当前使用的 registry，更新当前使用的 registry
	if current != nil && current.Name == oldName {
		if err := m.Use(newName); err != nil {
			return fmt.Errorf("failed to update current registry: %v", err)
		}
	}

	return nil
}
