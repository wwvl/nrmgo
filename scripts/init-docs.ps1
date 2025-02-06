# 初始化文档目录
$docDirs = @(
    "docs/.vitepress",
    "docs/guide",
    "docs/api",
    "docs/examples",
    "docs/public"
)

foreach ($dir in $docDirs) {
    if (-not (Test-Path $dir)) {
        Write-Host "Creating directory: $dir"
        New-Item -ItemType Directory -Force -Path $dir | Out-Null
    }
}

# 创建VitePress配置文件
$configContent = @"
import { defineConfig } from 'vitepress'

export default defineConfig({
    title: 'NRMGO',
    description: 'NPM Registry Manager GO',
    
    themeConfig: {
        nav: [
            { text: '首页', link: '/' },
            { text: '指南', link: '/guide/' },
            { text: 'API', link: '/api/' },
            { text: '示例', link: '/examples/' }
        ],
        
        sidebar: {
            '/guide/': [
                {
                    text: '入门',
                    items: [
                        { text: '简介', link: '/guide/' },
                        { text: '快速开始', link: '/guide/getting-started' },
                        { text: '安装', link: '/guide/installation' }
                    ]
                },
                {
                    text: '基础',
                    items: [
                        { text: '基本用法', link: '/guide/basic-usage' },
                        { text: '配置', link: '/guide/configuration' },
                        { text: 'Registry管理', link: '/guide/registry-management' }
                    ]
                },
                {
                    text: '进阶',
                    items: [
                        { text: '自定义Registry', link: '/guide/custom-registry' },
                        { text: '备份与恢复', link: '/guide/backup-restore' },
                        { text: '性能优化', link: '/guide/performance' }
                    ]
                }
            ],
            '/api/': [
                {
                    text: 'API参考',
                    items: [
                        { text: '命令行接口', link: '/api/' },
                        { text: '配置文件', link: '/api/configuration' },
                        { text: '环境变量', link: '/api/environment' }
                    ]
                }
            ]
        },
        
        socialLinks: [
            { icon: 'github', link: 'https://github.com/wwvl/nrmgo' }
        ]
    }
})
"@

$configPath = "docs/.vitepress/config.ts"
Write-Host "Creating VitePress config: $configPath"
Set-Content -Path $configPath -Value $configContent

# 创建首页
$indexContent = @"
---
layout: home

hero:
  name: NRMGO
  text: NPM Registry Manager GO
  tagline: 快速、简单、高效的NPM Registry管理工具
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: 在GitHub上查看
      link: https://github.com/wwvl/nrmgo

features:
  - icon: 🚀
    title: 快速切换
    details: 支持快速切换不同的NPM Registry，自动选择最快的镜像源
  - icon: 📦
    title: 多包管理器支持
    details: 支持npm、yarn、pnpm和bun等多种包管理器
  - icon: 🔄
    title: 智能管理
    details: 自动测试Registry延迟，智能选择最快的镜像源
  - icon: 💾
    title: 配置备份
    details: 支持配置文件的备份与恢复，确保数据安全
---
"@

$indexPath = "docs/index.md"
Write-Host "Creating index page: $indexPath"
Set-Content -Path $indexPath -Value $indexContent

# 创建基础文档
$basicDocs = @{
    "docs/guide/index.md" = @"
# 简介

NRMGO是一个用Go语言编写的NPM Registry管理工具，支持快速切换不同的NPM Registry。

## 特性

- 支持多个包管理器（npm、yarn、pnpm、bun）
- 自动测试Registry延迟
- 智能选择最快的Registry
- 缓存优化
- 友好的命令行界面
- 进度显示支持
"@
    
    "docs/guide/getting-started.md" = @"
# 快速开始

## 安装

```bash
go install github.com/wwvl/nrmgo@latest
```

## 基本使用

1. 列出所有可用的Registry：
\`\`\`bash
nrmgo ls
\`\`\`

2. 切换到指定的Registry：
\`\`\`bash
nrmgo use taobao
\`\`\`

3. 自动选择最快的Registry：
\`\`\`bash
nrmgo use
\`\`\`
"@
    
    "docs/api/index.md" = @"
# 命令行接口

NRMGO提供了一系列命令行接口，用于管理NPM Registry。

## 基础命令

- \`nrmgo ls\`: 列出所有可用的Registry
- \`nrmgo use [registry]\`: 切换Registry
- \`nrmgo add <name> <url>\`: 添加自定义Registry
- \`nrmgo del <name>\`: 删除Registry
- \`nrmgo current\`: 显示当前使用的Registry
"@
}

foreach ($doc in $basicDocs.GetEnumerator()) {
    Write-Host "Creating document: $($doc.Key)"
    Set-Content -Path $doc.Key -Value $doc.Value
}

# 创建package.json
$packageJson = @"
{
  "type": "module",
  "scripts": {
    "docs:dev": "vitepress dev docs",
    "docs:build": "vitepress build docs",
    "docs:preview": "vitepress preview docs"
  },
  "devDependencies": {
    "vitepress": "^1.0.0-rc.40"
  }
}
"@

$packageJsonPath = "package.json"
Write-Host "Creating package.json: $packageJsonPath"
Set-Content -Path $packageJsonPath -Value $packageJson

Write-Host "`n✨ Documentation initialization completed!"
Write-Host "Run the following commands to start the development server:"
Write-Host "pnpm install"
Write-Host "pnpm run docs:dev" 