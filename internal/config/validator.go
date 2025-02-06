package config

import (
	"fmt"
	"net/url"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateRegistry 验证 registry 配置
func ValidateRegistry(name string, reg *Registry) error {
	if reg == nil {
		return &ValidationError{
			Field:   "registry",
			Message: fmt.Sprintf("registry %s is nil", name),
		}
	}

	// 验证 URL
	if reg.URL == "" {
		return &ValidationError{
			Field:   "registry.url",
			Message: fmt.Sprintf("registry %s: url is required", name),
		}
	}

	// 验证 URL 格式
	if _, err := url.Parse(reg.URL); err != nil {
		return &ValidationError{
			Field:   "registry.url",
			Message: fmt.Sprintf("registry %s: invalid url format: %v", name, err),
		}
	}

	return nil
}

// ValidateConfig 验证整个配置
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return &ValidationError{
			Field:   "config",
			Message: "config is nil",
		}
	}

	// 验证并发请求数
	if cfg.MaxConcurrentRequests < 1 {
		cfg.MaxConcurrentRequests = 5 // 使用默认值
	} else if cfg.MaxConcurrentRequests > 10 {
		return &ValidationError{
			Field:   "max_concurrent_requests",
			Message: "value must be between 1 and 10",
		}
	}

	// 验证所有自定义 registry
	for name, reg := range cfg.CustomRegistries {
		if err := ValidateRegistry(name, reg); err != nil {
			return err
		}
	}

	return nil
}
