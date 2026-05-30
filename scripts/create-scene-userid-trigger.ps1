param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android",
    [string]$UserId = "1802965102482800640"
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
                userType   = "userID"
                userIDs    = @($UserId)
                params     = @{ body = "userID notify test"; title = "scene userID" }
            }
        }
    )
}
$createBodyObj = @{
    name       = "msg-notify-debug-userid"
    areaID     = "2"
    type       = "manual"
    tag        = "normal"
    deviceMode = "multi"
    if         = '{"triggers":[]}'
    when       = '{"validRanges":[],"invalidRanges":[],"conditions":{"type":"and","terms":[]}}'
    then       = ($thenObj | ConvertTo-Json -Depth 10 -Compress)
}
$created = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/create" -Headers $auth -Body ($createBodyObj | ConvertTo-Json -Compress)
Write-Host "create code=$($created.code) id=$($created.data.id)"
$sid = [int64]$created.data.id

$st1 = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
$c1 = 0
foreach ($item in $st1.data.list) { if ($item.count -lt 200 -and $item.group -notlike "*feedBack*") { $c1 = [int]$item.count } }

$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body ('{"id":' + $sid + '}')
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 3

$read = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/read" -Headers $auth -Body ('{"id":' + $sid + '}')
Write-Host "action status=$($read.data.then.actions[0].status) reason=$($read.data.then.actions[0].reason) isAbnormal=$($read.data.then.actions[0].isAbnormal)"

$st2 = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
$c2 = 0
foreach ($item in $st2.data.list) { if ($item.count -lt 200 -and $item.group -notlike "*feedBack*") { $c2 = [int]$item.count } }
Write-Host "count before=$c1 after=$c2 delta=$($c2-$c1)"
