# 要删除的文件列表
$filesToDelete = @(
    "internal/checker/detector.go",
    "internal/checker/npm.go",
    "internal/checker/yarn.go",
    "internal/checker/bun.go",
    "internal/checker/config.go"
)

# 遍历并删除文件
foreach ($file in $filesToDelete) {
    if (Test-Path $file) {
        Write-Host "Deleting $file..."
        Remove-Item $file -Force
    }
}

Write-Host "Cleanup completed." 