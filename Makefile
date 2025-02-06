.PHONY: all build lint clean

# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go 参数
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO_FLAGS := -v
GO_LDFLAGS := -X nrmgo/internal/cli.Version=$(VERSION) \
              -X nrmgo/internal/cli.Commit=$(COMMIT) \
              -X nrmgo/internal/cli.BuildTime=$(BUILD_TIME)

# 输出目录
OUT_DIR := bin
BINARY := nrmgo
CONFIG_FILE := internal/config/config.toml

# 平台列表
PLATFORMS := linux-amd64 linux-arm64 linux-386 \
            windows-amd64 windows-arm64 windows-386 \
            darwin-amd64 darwin-arm64

# 所有目标
all: clean build lint

# 构建单个平台
build-platform:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_FLAGS) -ldflags "$(GO_LDFLAGS)" \
		-o $(OUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)/$(BINARY)$(if $(filter windows,$(GOOS)),.exe,)
	@cp $(CONFIG_FILE) $(OUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)/

# 构建当前平台
build:
	@$(MAKE) build-platform

# 构建所有平台
build-all:
	@echo "Building for all platforms..."
	@$(foreach platform,$(PLATFORMS),\
		$(eval os = $(word 1,$(subst -, ,$(platform)))) \
		$(eval arch = $(word 2,$(subst -, ,$(platform)))) \
		GOOS=$(os) GOARCH=$(arch) $(MAKE) build-platform; \
	)

# 代码检查
lint:
	@echo "Running linters..."
	golangci-lint run

# 清理构建产物
clean:
	@echo "Cleaning..."
	rm -rf $(OUT_DIR)
	go clean

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# 检查依赖更新
check-updates:
	@echo "Checking for dependency updates..."
	go list -u -m all

# 更新依赖
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# 创建新版本
release:
	@echo "Creating new release..."
	@read -p "Enter version tag (e.g., v1.0.0): " tag; \
	git tag -a $$tag -m "Release $$tag"; \
	git push origin $$tag

# 显示版本信息
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Built:   $(BUILD_TIME)"

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all          - Clean, build, and lint"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  lint         - Run linters"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  check-updates - Check for dependency updates"
	@echo "  update-deps  - Update dependencies"
	@echo "  release      - Create new release"
	@echo "  version      - Show version information"
	@echo "  help         - Show this help" 