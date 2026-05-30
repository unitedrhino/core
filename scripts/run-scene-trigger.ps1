param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android",
    [int64]$SceneId = 5550
)

$ErrorActionPreference = "Stop"
$MsgGroup = [string]::Concat([char]0x573A, [char]0x666F, [char]0x8054, [char]0x52A8, [char]0x901A, [char]0x77E5)

function Invoke-Api {
    param([string]$Uri, [hashtable]$Headers = @{}, [string]$Body = "{}")
    $h = @{ "Content-Type" = "application/json" }
    foreach ($k in $Headers.Keys) { $h[$k] = $Headers[$k] }
    $raw = Invoke-WebRequest -Uri $Uri -Method POST -Headers $h -Body $Body -UseBasicParsing -TimeoutSec 30
    return ($raw.Content | ConvertFrom-Json)
}

$auth = @{
    "ithings-token"      = $Token
    "ithings-project-id" = $ProjectId
    "ithings-app-code"   = $AppCode
}

Write-Host "=== SCENE READ id=$SceneId ===" -ForegroundColor Cyan
$readBody = '{"id":' + $SceneId + '}'
$scene = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/read" -Headers $auth -Body $readBody
if ($scene.code -ne 200) {
    Write-Host "read fail $($scene.msg)" -ForegroundColor Red
    exit 1
}
Write-Host "name=$($scene.data.name) status=$($scene.data.status)"
$thenJson = $scene.data.then | ConvertTo-Json -Depth 20
Write-Host "then=$thenJson"

Write-Host ""
Write-Host "=== NOTIFY CONFIG ===" -ForegroundColor Cyan
$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
if ($cfg.data.list) {
    $c = $cfg.data.list[0]
    Write-Host "supportTypes=$($c.supportTypes -join ',') enableTypes=$($c.enableTypes -join ',')"
}

Write-Host ""
Write-Host "=== TEMPLATE ===" -ForegroundColor Cyan
$tpl = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/index" -Headers $auth -Body '{"notifyCode":"ruleScene","type":"message","page":{"page":1,"size":10}}'
if ($tpl.data.list) {
    Write-Host "id=$($tpl.data.list[0].id) subject=$($tpl.data.list[0].subject)"
}

Write-Host ""
Write-Host "=== TRIGGER ===" -ForegroundColor Cyan
$groupJson = '{"page":{"page":1,"size":20},"group":"' + $MsgGroup + '"}'
$mBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$totalBefore = [int]$mBefore.data.total
Write-Host "messages before=$totalBefore"

$triggerBody = '{"id":' + $SceneId + '}'
$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body $triggerBody
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 3
$mAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$totalAfter = [int]$mAfter.data.total
Write-Host "messages after=$totalAfter delta=$($totalAfter - $totalBefore)"
if ($mAfter.data.list -and $mAfter.data.list.Count -gt 0) {
    Write-Host "latest subject=$($mAfter.data.list[0].subject)"
    Write-Host "latest body=$($mAfter.data.list[0].body)"
}
