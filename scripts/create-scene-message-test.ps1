param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android",
    [string]$Account = "18059688688"
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

$thenObj = @{
    actions = @(
        @{
            id     = 1
            order  = 1
            type   = "notify"
            status = 1
            notify = @{
                type       = "message"
                notifyCode = "ruleScene"
                userType   = "account"
                accounts   = @($Account)
                params     = @{ body = "scene message debug test"; title = "scene notify" }
            }
        }
    )
}
$thenStr = ($thenObj | ConvertTo-Json -Depth 10 -Compress)

$createBodyObj = @{
    name       = "msg-notify-debug"
    areaID     = "2"
    type       = "manual"
    tag        = "normal"
    deviceMode = "multi"
    if         = '{"triggers":[]}'
    when       = '{"validRanges":[],"invalidRanges":[],"conditions":{"type":"and","terms":[]}}'
    then       = $thenStr
}
$createBody = $createBodyObj | ConvertTo-Json -Compress

Write-Host "=== CREATE SCENE ===" -ForegroundColor Cyan
$created = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/create" -Headers $auth -Body $createBody
Write-Host "code=$($created.code) msg=$($created.msg)"
if ($created.code -ne 200) { exit 1 }
$sceneId = [int64]$created.data.id
Write-Host "sceneId=$sceneId" -ForegroundColor Green

Write-Host ""
Write-Host "=== TRIGGER ===" -ForegroundColor Cyan
$groupJson = '{"page":{"page":1,"size":20},"group":"' + $MsgGroup + '"}'
$mBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$totalBefore = [int]$mBefore.data.total
$triggerBody = '{"id":' + $sceneId + '}'
$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body $triggerBody
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 3
$mAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$totalAfter = [int]$mAfter.data.total
Write-Host "messages before=$totalBefore after=$totalAfter delta=$($totalAfter - $totalBefore)"
if ($totalAfter -gt $totalBefore) {
    Write-Host "SUCCESS" -ForegroundColor Green
    Write-Host "subject=$($mAfter.data.list[0].subject)"
    Write-Host "body=$($mAfter.data.list[0].body)"
} else {
    Write-Host "FAIL - no message" -ForegroundColor Yellow
}

Write-Host "sceneId=$sceneId"
