package cli

import (
	"fmt"
	"os"

	"nrmgo/internal/config"
	"nrmgo/internal/registry"
)

// loadConfigAndCreateManager 加载配置并创建 registry 管理器
// 返回配置对象、registry 管理器和可能的错误
func loadConfigAndCreateManager() (*config.Config, registry.Manager, error) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("\n⚠️  Config file not found, please use 'nrmgo config init' to create it")
		}
		return nil, nil, fmt.Errorf("\n❌  Failed to load config: %v", err)
	}

	// 创建 registry 管理器
	manager := registry.NewManager(cfg)

	return cfg, manager, nil
}
