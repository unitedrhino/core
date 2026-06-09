# 部署 apisvr 到 new.ykhl.vip 环境（47.94.112.109）
# 用法（在已安装 Go、可 SSH 到服务器的机器上执行）:
#   cd e:\web-core\core
#   .\scripts\deploy-apisvr-new-ykhl.ps1
# 可选环境变量:
#   $env:DEPLOY_HOST = "root@47.94.112.109"
#   $env:DEPLOY_PATH = "/root/run/core"

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
if (-not (Test-Path "$Root\Makefile")) {
    throw "请在 core 仓库根目录下使用本脚本（未找到 Makefile）"
}

$DeployHost = if ($env:DEPLOY_HOST) { $env:DEPLOY_HOST } else { "root@47.94.112.109" }
$DeployPath = if ($env:DEPLOY_PATH) { $env:DEPLOY_PATH } else { "/root/run/core" }

Set-Location $Root
Write-Host ">> 编译 apisvr（含内嵌 syssvr，目标 linux/amd64）..."
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
& go build -tags no_k8s -o ./cmd/coresvr ./service/apisvr
if (-not (Test-Path "./cmd/coresvr") -and -not (Test-Path "./cmd/coresvr.exe")) {
    throw "编译失败，请确认 go 在 PATH 中"
}

# 复制配置（与 Makefile cp.etc 一致）
New-Item -ItemType Directory -Force -Path ./cmd/etc | Out-Null
Copy-Item -Force ./service/apisvr/etc/* ./cmd/etc/

$BinName = if (Test-Path "./cmd/coresvr") { "coresvr" } else { "coresvr.exe" }
Write-Host ">> 上传到 ${DeployHost}:${DeployPath} ..."
& scp ./cmd/coresvr "${DeployHost}:${DeployPath}/coresvr"
& scp -r ./cmd/etc "${DeployHost}:${DeployPath}/"

Write-Host ">> 重启 apisvr ..."
$remoteCmd = "cd $DeployPath && chmod +x coresvr && pkill -f './coresvr' 2>/dev/null || pkill -f 'apisvr' 2>/dev/null || true && sleep 1 && setsid ./coresvr > /tmp/coresvr.log 2>&1 & sleep 2 && pgrep -af coresvr || pgrep -af apisvr"
& ssh $DeployHost $remoteCmd

Write-Host ">> 完成。验证 Apple 登录枚举:"
Write-Host 'Invoke-RestMethod -Uri "https://new.ykhl.vip/api/v1/system/user/self/login" -Method Post -ContentType "application/json" -Body ''{"loginType":"apple","code":"test"}'''
