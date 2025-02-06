package registry

import (
	"fmt"
)

// ErrRegistryNotFound 表示未找到指定的 registry
type ErrRegistryNotFound struct {
	Name string
}

func (e *ErrRegistryNotFound) Error() string {
	return fmt.Sprintf("registry not found: %s", e.Name)
}

// ErrRegistryExists 表示 registry 已存在
type ErrRegistryExists struct {
	Name string
}

func (e *ErrRegistryExists) Error() string {
	return fmt.Sprintf("registry already exists: %s", e.Name)
}

// ErrBuiltinRegistry 表示试图修改内置 registry
type ErrBuiltinRegistry struct {
	Name string
}

func (e *ErrBuiltinRegistry) Error() string {
	return fmt.Sprintf("cannot modify builtin registry: %s", e.Name)
}

// ErrInvalidRegistry 表示 registry 配置无效
type ErrInvalidRegistry struct {
	Name   string
	Reason string
}

func (e *ErrInvalidRegistry) Error() string {
	return fmt.Sprintf("invalid registry %s: %s", e.Name, e.Reason)
}
