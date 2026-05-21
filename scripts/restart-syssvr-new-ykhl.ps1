# 重启 new.ykhl.vip 环境（47.94.112.109）上的 syssvr
# 说明：若 apisvr 使用 SysRpc.Mode=direct + RunProxy，syssvr 内嵌在 coresvr 进程内，需同时重启 coresvr 才能加载新配置。
# 用法:
#   cd e:\web-core\core
#   .\scripts\restart-syssvr-new-ykhl.ps1
# 可选环境变量:
#   $env:DEPLOY_HOST = "root@47.94.112.109"
#   $env:DEPLOY_PATH = "/root/run/core"

$ErrorActionPreference = "Stop"
$DeployHost = if ($env:DEPLOY_HOST) { $env:DEPLOY_HOST } else { "root@47.94.112.109" }
$DeployPath = if ($env:DEPLOY_PATH) { $env:DEPLOY_PATH } else { "/root/run/core" }

Write-Host ">> 连接 $DeployHost 重启 syssvr（必要时重启 coresvr）..."
& ssh $DeployHost @"
set -e
cd $DeployPath
echo '--- 重启前进程 ---'
pgrep -af syssvr || true
pgrep -af coresvr || true

# 独立 syssvr 进程（若存在）
if [ -f ./syssvr ]; then
  pkill -f './syssvr' 2>/dev/null || true
  sleep 1
  chmod +x ./syssvr
  setsid ./syssvr > /tmp/syssvr.log 2>&1 &
  sleep 2
  echo '--- 已启动独立 syssvr ---'
  pgrep -af syssvr || echo 'syssvr 未运行'
fi

# direct 模式下 syssvr 内嵌在 coresvr，重启 coresvr 以刷新 OAuth 客户端缓存
if [ -f ./coresvr ]; then
  pkill -f './coresvr' 2>/dev/null || pkill -f 'apisvr' 2>/dev/null || true
  sleep 1
  chmod +x ./coresvr
  setsid ./coresvr > /tmp/coresvr.log 2>&1 &
  sleep 2
  echo '--- 已重启 coresvr（含内嵌 syssvr）---'
  pgrep -af coresvr || echo 'coresvr 未运行'
fi

echo '--- 重启后进程 ---'
pgrep -af 'syssvr|coresvr' || true
"@

Write-Host ">> 完成"
