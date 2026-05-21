# 重启 new.ykhl.vip 上的 coresvr（含内嵌 syssvr）
# 用法: cd e:\web-core\core ; .\scripts\restart-coresvr-new-ykhl.ps1

$ErrorActionPreference = "Stop"
$DeployHost = if ($env:DEPLOY_HOST) { $env:DEPLOY_HOST } else { "root@47.94.112.109" }
$DeployPath = if ($env:DEPLOY_PATH) { $env:DEPLOY_PATH } else { "/root/run/core" }

Write-Host ">> 重启 ${DeployHost}:${DeployPath} 上的 coresvr ..."
& ssh $DeployHost @"
cd $DeployPath
echo '--- before ---'
pgrep -af coresvr || true
pkill -f './coresvr' 2>/dev/null || pkill -f 'apisvr' 2>/dev/null || true
sleep 1
chmod +x coresvr
setsid ./coresvr > /tmp/coresvr.log 2>&1 &
sleep 2
echo '--- after ---'
pgrep -af coresvr || echo 'coresvr not running'
tail -3 /tmp/coresvr.log 2>/dev/null || true
"@

Write-Host ">> 完成"
