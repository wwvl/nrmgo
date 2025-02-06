package config

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// RenderConfig 将配置渲染为树形结构
func RenderConfig(cfg *Config) error {
	fmt.Println("\n🌈  NRMGO Configuration")

	// 创建树形结构的数据
	leveledList := pterm.LeveledList{
		{Level: 0, Text: "📦 Configuration"},
	}

	// 添加 custom_registries
	if len(cfg.CustomRegistries) > 0 {
		leveledList = append(leveledList, pterm.LeveledListItem{
			Level: 1,
			Text:  "📂 custom_registries",
		})

		// 添加每个 registry
		for name, reg := range cfg.CustomRegistries {
			leveledList = append(leveledList, pterm.LeveledListItem{
				Level: 2,
				Text:  fmt.Sprintf("📦 %s", name),
			})

			// 添加 registry 的字段
			leveledList = append(leveledList,
				pterm.LeveledListItem{Level: 3, Text: fmt.Sprintf("🔗 url: %q", reg.URL)},
			)

			// 可选字段
			if reg.Home != "" {
				leveledList = append(leveledList,
					pterm.LeveledListItem{Level: 3, Text: fmt.Sprintf("🏠 home: %q", reg.Home)},
				)
			}
			if reg.Description != "" {
				leveledList = append(leveledList,
					pterm.LeveledListItem{Level: 3, Text: fmt.Sprintf("📝 description: %q", reg.Description)},
				)
			}
		}
	}

	// 添加其他配置项
	leveledList = append(leveledList,
		pterm.LeveledListItem{Level: 1, Text: fmt.Sprintf("🔢 max_concurrent_requests: %d", cfg.MaxConcurrentRequests)},
	)

	// 添加空行
	fmt.Println()

	// 创建并渲染树形结构
	return pterm.DefaultTree.WithRoot(putils.TreeFromLeveledList(leveledList)).Render()
}
