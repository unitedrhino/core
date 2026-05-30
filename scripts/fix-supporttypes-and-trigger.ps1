param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android",
    [int64]$SceneId = 5551
)

$ErrorActionPreference = "Stop"
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

$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
$c = $cfg.data.list[0]
Write-Host "config id=$($c.id) supportTypes=$($c.supportTypes -join ',')"

$updBody = @{
    id           = $c.id
    code         = $c.code
    name         = $c.name
    group        = $c.group
    isRecord     = $c.isRecord
    supportTypes = @("message","sms","email","dingTalk","phoneCall")
    enableTypes  = @("message","sms")
} | ConvertTo-Json -Compress
$upd = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/update" -Headers $auth -Body $updBody
Write-Host "update code=$($upd.code) msg=$($upd.msg)"

$cfg2 = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
Write-Host "after supportTypes=$($cfg2.data.list[0].supportTypes -join ',')"

$st1 = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
$count1 = 0
foreach ($item in $st1.data.list) {
    if ($item.group -notlike "*feedBack*" -and $item.count -lt 200) { $count1 = [int]$item.count }
}
$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body ('{"id":' + $SceneId + '}')
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 3
$st2 = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
$count2 = 0
foreach ($item in $st2.data.list) {
    if ($item.group -notlike "*feedBack*" -and $item.count -lt 200) { $count2 = [int]$item.count }
}
Write-Host "scene-notify count before=$count1 after=$count2 delta=$($count2-$count1)"
