# 发布新版本
param (
    [string]$version = "",
    [switch]$dryRun = $false
)

function Test-ValidVersion {
    param ([string]$v)
    return $v -match '^v\d+\.\d+\.\d+$'
}

function Get-LatestTag {
    $tags = git tag -l "v*" | Sort-Object -Property {[version]($_ -replace '^v','')} -Descending
    return $tags[0]
}

function Show-ReleaseInfo {
    param (
        [string]$version,
        [string]$commit,
        [string]$changes
    )
    Write-Host "`n发布信息预览：" -ForegroundColor Cyan
    Write-Host "版本号: $version" -ForegroundColor Yellow
    Write-Host "提交号: $commit" -ForegroundColor Yellow
    Write-Host "`n更改日志：`n$changes" -ForegroundColor Gray
}

# 检查工作目录是否干净
$status = git status --porcelain
if ($status) {
    Write-Host "错误：工作目录不干净，请提交或存储更改。" -ForegroundColor Red
    exit 1
}

# 获取当前分支
$branch = git rev-parse --abbrev-ref HEAD
if ($branch -ne "main") {
    Write-Host "警告：当前不在 main 分支上。" -ForegroundColor Yellow
    $continue = Read-Host "是否继续？(y/N)"
    if ($continue -ne "y") {
        exit 0
    }
}

# 拉取最新代码
git pull origin main

# 获取版本号
if (-not $version) {
    $latestTag = Get-LatestTag
    Write-Host "当前最新版本：$latestTag" -ForegroundColor Cyan
    $version = Read-Host "请输入新版本号 (例如 v1.0.0)"
}

# 验证版本号格式
if (-not (Test-ValidVersion $version)) {
    Write-Host "错误：无效的版本号格式。必须符合 'v1.0.0' 格式。" -ForegroundColor Red
    exit 1
}

# 检查版本号是否已存在
$existingTag = git tag -l $version
if ($existingTag) {
    Write-Host "错误：版本 $version 已存在。" -ForegroundColor Red
    exit 1
}

# 获取上次发布以来的提交日志
$lastTag = Get-LatestTag
$changes = git log --pretty=format:"- %s" "${lastTag}..HEAD"
$commit = git rev-parse --short HEAD

# 显示发布信息
Show-ReleaseInfo $version $commit $changes

# 确认发布
if (-not $dryRun) {
    $confirm = Read-Host "`n是否确认发布？(y/N)"
    if ($confirm -ne "y") {
        Write-Host "取消发布。" -ForegroundColor Yellow
        exit 0
    }

    # 创建版本标签
    git tag -a $version -m "Release $version"
    git push origin $version

    Write-Host "`n✨ 版本 $version 发布成功！" -ForegroundColor Green
    Write-Host "请等待 CI/CD 流程完成..." -ForegroundColor Cyan
} else {
    Write-Host "`n这是一个演习，没有实际创建发布。" -ForegroundColor Yellow
} 