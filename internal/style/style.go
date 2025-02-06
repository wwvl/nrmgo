package style

import "github.com/pterm/pterm"

// Success 成功样式（绿色加粗）
var Success = pterm.NewStyle(pterm.FgGreen, pterm.Bold)

// Error 错误样式（红色加粗）
var Error = pterm.NewStyle(pterm.FgRed, pterm.Bold)

// Warning 警告样式（黄色加粗）
var Warning = pterm.NewStyle(pterm.FgYellow, pterm.Bold)

// Info 信息样式（蓝色）
var Info = pterm.NewStyle(pterm.FgBlue)
