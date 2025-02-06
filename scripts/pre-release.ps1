# 发布前检查
param (
    [switch]$fix = $false
)

$success = $true

function Write-Step {
    param (
        [string]$message
    )
    Write-Host "`n>> $message" -ForegroundColor Cyan
}

function Write-Result {
    param (
        [string]$message,
        [bool]$passed
    )
    if ($passed) {
        Write-Host "✓ $message" -ForegroundColor Green
    } else {
        Write-Host "✗ $message" -ForegroundColor Red
        $script:success = $false
    }
}

# 检查工作目录
Write-Step "检查工作目录"
$status = git status --porcelain
if ($status) {
    Write-Result "工作目录不干净" $false
    Write-Host $status -ForegroundColor Yellow
    if ($fix) {
        $confirm = Read-Host "是否提交所有更改？(y/N)"
        if ($confirm -eq "y") {
            git add .
            git commit -m "chore: pre-release cleanup"
        }
    }
} else {
    Write-Result "工作目录干净" $true
}

# 检查分支
Write-Step "检查分支"
$branch = git rev-parse --abbrev-ref HEAD
Write-Result "当前分支: $branch" ($branch -eq "main")

# 检查远程同步状态
Write-Step "检查远程同步状态"
git fetch origin
$ahead = git rev-list HEAD..origin/main --count
$behind = git rev-list origin/main..HEAD --count
Write-Result "与远程分支同步" ($ahead -eq 0 -and $behind -eq 0)
if ($ahead -gt 0) {
    Write-Host "本地领先 $ahead 个提交" -ForegroundColor Yellow
}
if ($behind -gt 0) {
    Write-Host "本地落后 $behind 个提交" -ForegroundColor Yellow
    if ($fix) {
        $confirm = Read-Host "是否拉取远程更新？(y/N)"
        if ($confirm -eq "y") {
            git pull origin main
        }
    }
}

# 检查构建
Write-Step "检查构建"
try {
    $buildOutput = make build 2>&1
    Write-Result "构建成功" $true
} catch {
    Write-Result "构建失败" $false
    Write-Host $buildOutput -ForegroundColor Red
}

# 检查代码质量
Write-Step "检查代码质量"
try {
    $lintOutput = make lint 2>&1
    Write-Result "代码检查通过" $true
} catch {
    Write-Result "代码检查失败" $false
    Write-Host $lintOutput -ForegroundColor Red
    if ($fix) {
        $confirm = Read-Host "是否运行代码格式化？(y/N)"
        if ($confirm -eq "y") {
            make fmt
        }
    }
}

# 检查文档
Write-Step "检查文档"
$docFiles = @("README.md", "docs/developer-guide.md")
foreach ($file in $docFiles) {
    if (Test-Path $file) {
        Write-Result "找到 $file" $true
    } else {
        Write-Result "缺少 $file" $false
    }
}

# 总结
if ($success) {
    Write-Host "`n✨ 发布前检查通过！可以继续发布流程。" -ForegroundColor Green
} else {
    Write-Host "`n❌ 发布前检查失败！请修复上述问题后重试。" -ForegroundColor Red
    if (-not $fix) {
        Write-Host "提示：使用 -fix 参数运行脚本可以自动修复部分问题。" -ForegroundColor Yellow
    }
    exit 1
} 