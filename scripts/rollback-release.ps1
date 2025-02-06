# 回滚发布
param (
    [string]$version = "",
    [switch]$force = $false
)

function Get-LatestTag {
    $tags = git tag -l "v*" | Sort-Object -Property {[version]($_ -replace '^v','')} -Descending
    return $tags[0]
}

# 如果没有指定版本，使用最新的标签
if (-not $version) {
    $version = Get-LatestTag
    if (-not $version) {
        Write-Host "错误：找不到任何版本标签。" -ForegroundColor Red
        exit 1
    }
}

Write-Host "准备回滚版本 $version ..." -ForegroundColor Yellow

# 确认回滚
if (-not $force) {
    Write-Host "`n警告：此操作将删除版本标签并可能影响其他开发者。" -ForegroundColor Red
    $confirm = Read-Host "是否确定要回滚？(y/N)"
    if ($confirm -ne "y") {
        Write-Host "取消回滚。" -ForegroundColor Yellow
        exit 0
    }
}

# 检查标签是否存在
$tag = git tag -l $version
if (-not $tag) {
    Write-Host "错误：找不到版本标签 $version" -ForegroundColor Red
    exit 1
}

try {
    # 删除本地标签
    Write-Host "删除本地标签..." -NoNewline
    git tag -d $version
    Write-Host " 完成" -ForegroundColor Green

    # 删除远程标签
    Write-Host "删除远程标签..." -NoNewline
    git push origin :refs/tags/$version
    Write-Host " 完成" -ForegroundColor Green

    # 清理构建产物
    if (Test-Path "bin") {
        Write-Host "清理构建产物..." -NoNewline
        Remove-Item -Recurse -Force bin
        Write-Host " 完成" -ForegroundColor Green
    }

    Write-Host "`n✨ 版本 $version 已成功回滚！" -ForegroundColor Green
    
    # 显示当前最新版本
    $currentVersion = Get-LatestTag
    if ($currentVersion) {
        Write-Host "当前最新版本：$currentVersion" -ForegroundColor Cyan
    }
} catch {
    Write-Host "`n❌ 回滚过程中出错：" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
} 