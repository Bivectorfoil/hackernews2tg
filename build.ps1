# 设置环境变量以禁用 CGO
$env:CGO_ENABLED = "0"

# 创建输出目录
New-Item -ItemType Directory -Force -Path "dist"

# 设置编译时的优化标志
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$ldflags = "-w -s"  # -w 删除调试信息, -s 删除符号表

Write-Host "Building windows version..."
go build -ldflags $ldflags -o "dist/newsboy.exe"

# 如果需要同时构建 Linux 版本
$env:GOOS = "linux"
Write-Host "Building linux version..."
go build -ldflags $ldflags -o "dist/newsboy"

Write-Host "Build completed!"
Write-Host "Windows version: dist/newsboy.exe"
Write-Host "Linux version: dist/newsboy"
