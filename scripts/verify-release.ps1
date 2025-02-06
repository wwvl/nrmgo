# 验证发布
param (
    [string]$version = ""
)

function Test-Binary {
    param (
        [string]$path,
        [string]$expectedVersion
    )
    
    if (-not (Test-Path $path)) {
        Write-Host "错误：找不到二进制文件 $path" -ForegroundColor Red
        return $false
    }

    try {
        $output = & $path version 2>&1
        if ($LASTEXITCODE -eq 0 -and $output -match $expectedVersion) {
            Write-Host "✓ 版本验证通过: $path" -ForegroundColor Green
            return $true
        } else {
            Write-Host "✗ 版本验证失败: $path" -ForegroundColor Red
            Write-Host "  预期版本: $expectedVersion" -ForegroundColor Yellow
            Write-Host "  实际输出: $output" -ForegroundColor Yellow
            return $false
        }
    } catch {
        Write-Host "✗ 执行失败: $path" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        return $false
    }
}

# 获取要验证的版本
if (-not $version) {
    $version = git describe --tags --abbrev=0
}

Write-Host "开始验证版本 $version ..." -ForegroundColor Cyan

# 验证 git tag
$tag = git tag -l $version
if (-not $tag) {
    Write-Host "错误：找不到版本标签 $version" -ForegroundColor Red
    exit 1
}

# 检查构建产物
$binaries = @(
    "bin/nrmgo-linux-amd64",
    "bin/nrmgo-windows-amd64.exe",
    "bin/nrmgo-darwin-amd64"
)

$success = $true

foreach ($binary in $binaries) {
    if (-not (Test-Binary $binary $version)) {
        $success = $false
    }
}

# 验证 CI/CD 状态
$repoUrl = git config --get remote.origin.url
if ($repoUrl -match "github.com[:/]([^/]+)/([^/]+)\.git$") {
    $owner = $matches[1]
    $repo = $matches[2]
    
    Write-Host "`n检查 GitHub Actions 运行状态..." -ForegroundColor Cyan
    $workflowUrl = "https://github.com/$owner/$repo/actions/workflows/ci.yml"
    Write-Host "请访问以下链接查看 CI/CD 状态：" -ForegroundColor Yellow
    Write-Host $workflowUrl -ForegroundColor Blue
}

if ($success) {
    Write-Host "`n✨ 版本 $version 验证通过！" -ForegroundColor Green
} else {
    Write-Host "`n❌ 版本 $version 验证失败！" -ForegroundColor Red
    exit 1
} 