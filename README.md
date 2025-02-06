# NRMGO - NPM Registry Manager Go

[![CI/CD](https://github.com/wwvl/nrmgo/actions/workflows/ci.yml/badge.svg)](https://github.com/wwvl/nrmgo/actions/workflows/ci.yml)
[![GitHub release](https://img.shields.io/github/release/wwvl/nrmgo.svg)](https://github.com/wwvl/nrmgo/releases)
[![License](https://img.shields.io/github/license/wwvl/nrmgo.svg)](https://github.com/wwvl/nrmgo/blob/main/LICENSE)

NRMGO 是一个用 Go 语言编写的 NPM Registry 管理工具，支持快速切换不同的 NPM Registry。

## ✨ 特性

- 支持多个包管理器（npm、yarn、pnpm、bun）
- 自动测试 Registry 延迟
- 智能选择最快的 Registry
- 友好的命令行界面
- 进度显示支持

## ⚙️ 配置

### 配置文件

配置文件在程序目录下：`./config.toml`

```toml
# HTTP concurrent requests for latency testing
# Default: 5, Recommended range: 1-10
max_concurrent_requests = 5

# User-defined registry list
[custom_registries]

# Custom registry example
# [custom_registries.example] # Registry name
# url = "https://example.com"  # Registry URL
# home = "https://example.com" # Registry homepage (optional)
# description = "example"      # Registry description (optional)
```

## 📦 安装

```bash
go install github.com/wwvl/nrmgo@latest
```

## 🚀 使用方法

### 列出所有可用的 Registry

```bash
$ nrmgo ls # 列出所有可用的 Registry

  NAME        REGISTRY URL                                   HOME PAGE
─────────────────────────────────────────────────────────────────────────────────────────────
  huawei      https://repo.huaweicloud.com/repository/npm/   https://www.huaweicloud.com/special/npm-jingxiang.html
  nju         https://repo.nju.edu.cn/repository/npm/        https://doc.nju.edu.cn/books/35f4a/page/npm
```

### 显示当前使用的 Registry

```bash
$ nrmgo current # 显示当前使用的 Registry

  PACKAGE MANAGER   STATUS   VERSION   REGISTRY   CONFIG FILE
───────────────────────────────────────────────────────────────────────────────────────
  npm               ✓        10.9.2    ustc       C:\Users\Administrator\.npmrc
  yarn              ✗        -         -          -
  pnpm              ✓        10.1.0    ustc       C:\Users\Administrator\.npmrc
  bun               ✓        1.2.0     ustc       C:\Users\Administrator\.bunfig.toml
```

### 自动选择最快的 Registry

```bash
$ nrmgo use  # 自动测试并选择最快的 Registry

Testing Registry latency...

Progress: 8/8

  NAME        REGISTRY URL                                   LATENCY
──────────────────────────────────────────────────────────────────────
  npm         https://registry.npmjs.org/                    1174 ms
  ustc        https://npmreg.proxy.ustclug.org/              1473 ms
  yarn        https://registry.yarnpkg.com/                  1585 ms
  taobao      https://registry.npmmirror.com/                1758 ms
  nju         https://repo.nju.edu.cn/repository/npm/        1818 ms
  npmMirror   https://skimdb.npmjs.com/registry/             1851 ms
  huawei      https://repo.huaweicloud.com/repository/npm/   2234 ms
  tencent     https://mirrors.tencent.com/npm/               Error

Using fastest Registry: npm


✓ Successfully switched 3 package manager(s) to npm Registry
```

### 切换到指定的 Registry

```bash
$ nrmgo use taobao  # 切换到淘宝 Registry

✓ Successfully switched 3 package manager(s) to taobao Registry
```

### 查看 Registry 详细信息

```bash
$ nrmgo info  # 查看所有 Registry 的详细信息

  NAME        HOMEPAGE                                                 REGISTRY URL                                   DESCRIPTION
─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  huawei      https://www.huaweicloud.com/special/npm-jingxiang.html   https://repo.huaweicloud.com/repository/npm/   Huawei Cloud npm registry mirror
  nju         https://doc.nju.edu.cn/books/35f4a/page/npm              https://repo.nju.edu.cn/repository/npm/        Nanjing University npm registry mirror
  npm         https://www.npmjs.org                                    https://registry.npmjs.org/                    Official npm registry
  npmMirror   https://skimdb.npmjs.com/                                https://skimdb.npmjs.com/registry/             npm registry mirror
  taobao      https://npmmirror.com                                    https://registry.npmmirror.com/                Taobao npm registry mirror
  tencent     https://mirrors.tencent.com/npm/                         https://mirrors.tencent.com/npm/               Tencent Cloud npm registry mirror
  ustc        https://mirrors.ustc.edu.cn/help/npm.html                https://npmreg.proxy.ustclug.org/              University of Science and Technology of China npm registry mirror
  yarn        https://yarnpkg.com                                      https://registry.yarnpkg.com/                  Official yarn registry

```

### 查看指定 Registry 的详细信息

```bash
nrmgo info taobao  # 查看淘宝 Registry 的详细信息

  NAME     HOMEPAGE                REGISTRY URL                      DESCRIPTION
─────────────────────────────────────────────────────────────────────────────────────────────────
  taobao   https://npmmirror.com   https://registry.npmmirror.com/   Taobao npm registry mirror
```

### 重置为默认镜像源

```bash
$ nrmgo reset

  PACKAGE MANAGER   STATUS   BACKUP FILE
──────────────────────────────────────────────────────────────────────────────────────
  npm               ✓        C:\Users\Administrator\.npmrc.20250202_002035.bak
  yarn              -        -
  pnpm              ✓        C:\Users\Administrator\.npmrc.20250202_002035.bak
  bun               ✓        C:\Users\Administrator\.bunfig.toml.20250202_002035.bak
```

### 高级功能

```bash
# 添加自定义镜像源
nrmgo add custom https://custom.registry.com/

# 删除镜像源
nrmgo del custom

# 查看镜像源详细信息
nrmgo info taobao

# 性能分析
nrmgo test --profile
```

## 📚 支持的 Registry

| 名称      | 说明              | Registry URL                                 |
| --------- | ----------------- | -------------------------------------------- |
| npm       | 官方 Registry     | https://registry.npmjs.org/                  |
| npmMirror | npm 镜像          | https://skimdb.npmjs.com/registry/           |
| yarn      | 官方 Yarn 镜像    | https://registry.yarnpkg.com/                |
| taobao    | 淘宝 NPM 镜像     | https://registry.npmmirror.com/              |
| huawei    | 华为云 NPM 镜像   | https://repo.huaweicloud.com/repository/npm/ |
| tencent   | 腾讯云 NPM 镜像   | http://mirrors.tencent.com/npm/              |
| ustc      | 中科大 NPM 镜像   | https://npmreg.proxy.ustclug.org/            |
| nju       | 南京大学 NPM 镜像 | https://repo.nju.edu.cn/repository/npm/      |

## 🔧 开发

### 环境要求

- Go 1.21 或更高版本
- Git
- Make
- golangci-lint

详细的开发指南请参考 [开发者指南](docs/developer-guide.md)。

### 构建

```bash
# 安装依赖
make deps

# 本地构建
make build

# 运行代码检查
make lint
```

## 🤝 贡献指南

欢迎贡献代码！请查看 [贡献指南](CONTRIBUTING.md) 了解详情。

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📝 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解详细的更新历史。

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [nrm](https://github.com/Pana/nrm) - 提供了部分镜像源数据
- [Go 社区](https://golang.org/) - 提供了优秀的开发语言和工具
- [npm](https://mirrors.tencent.com/help/npm.html) - 腾讯云 npm 镜像源
- [npmmirror 镜像站](https://npmmirror.com/) - 淘宝 npm 镜像源
- [NPM 镜像 - 华为云](https://www.huaweicloud.com/special/npm-jingxiang.html) - 华为云 npm 镜像源
- [npm | e-Science Document](https://doc.nju.edu.cn/books/e1654/page/npm) - 南京大学 npm 镜像源
- [NPM 反向代理 - USTC Mirror Help](https://mirrors.ustc.edu.cn/help/npm.html) - 中科大 npm 镜像源

## 📞 支持

- 提交 Issue: [GitHub Issues](https://github.com/wwvl/nrmgo/issues)
- 文档：[GoDoc](https://godoc.org/github.com/wwvl/nrmgo)
- 讨论：[GitHub Discussions](https://github.com/wwvl/nrmgo/discussions)
