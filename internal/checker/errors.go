package checker

import "fmt"

// DetectionError 表示包管理器检测过程中的错误
type DetectionError struct {
	PackageManager string
	Operation      string
	Err            error
}

func (e *DetectionError) Error() string {
	return fmt.Sprintf("%s: %s failed: %v", e.PackageManager, e.Operation, e.Err)
}

// ConfigError 表示配置文件操作错误
type ConfigError struct {
	PackageManager string
	Operation      string
	Path           string
	Err            error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("%s: %s config at %s failed: %v", e.PackageManager, e.Operation, e.Path, e.Err)
}

// NewDetectionError 创建检测错误
func NewDetectionError(pm, op string, err error) error {
	return &DetectionError{
		PackageManager: pm,
		Operation:      op,
		Err:            err,
	}
}

// NewConfigError 创建配置错误
func NewConfigError(pm, op, path string, err error) error {
	return &ConfigError{
		PackageManager: pm,
		Operation:      op,
		Path:           path,
		Err:            err,
	}
}
