# 检查必需的工具和环境
$requiredTools = @{
    "go" = "go version"
    "git" = "git --version"
    "make" = "make --version"
    "golangci-lint" = "golangci-lint --version"
}

$success = $true

Write-Host "检查开发环境..." -ForegroundColor Cyan

foreach ($tool in $requiredTools.Keys) {
    Write-Host "检查 $tool..." -NoNewline
    try {
        $version = Invoke-Expression $requiredTools[$tool] 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host " OK" -ForegroundColor Green
            Write-Host "  $version" -ForegroundColor Gray
        } else {
            Write-Host " 未安装" -ForegroundColor Red
            $success = $false
        }
    } catch {
        Write-Host " 未安装" -ForegroundColor Red
        $success = $false
    }
}

# 检查 Go 环境
Write-Host "`n检查 Go 环境..." -ForegroundColor Cyan
$goVersion = go version 2>&1
if ($LASTEXITCODE -eq 0) {
    if ($goVersion -match "go1\.(1[5-9]|[2-9][0-9])") {
        Write-Host "Go 版本符合要求" -ForegroundColor Green
    } else {
        Write-Host "Go 版本过低，需要 1.15 或更高版本" -ForegroundColor Red
        $success = $false
    }
} else {
    Write-Host "无法检查 Go 版本" -ForegroundColor Red
    $success = $false
}

# 检查 GOPATH 和 GOROOT
Write-Host "`n检查 Go 路径..." -ForegroundColor Cyan
$gopath = go env GOPATH
$goroot = go env GOROOT
Write-Host "GOPATH: $gopath"
Write-Host "GOROOT: $goroot"

# 检查项目依赖
Write-Host "`n检查项目依赖..." -ForegroundColor Cyan
$modResult = go mod verify 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "依赖检查通过" -ForegroundColor Green
} else {
    Write-Host "依赖检查失败:" -ForegroundColor Red
    Write-Host $modResult -ForegroundColor Red
    $success = $false
}

if (-not $success) {
    Write-Host "`n环境检查失败！请安装缺失的工具或修复问题。" -ForegroundColor Red
    exit 1
} else {
    Write-Host "`n✨ 环境检查通过！" -ForegroundColor Green
} 