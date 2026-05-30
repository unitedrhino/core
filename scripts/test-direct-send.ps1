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

$allBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5}}'
$tb = [int64]$allBefore.data.total
Write-Host "total before=$tb"

$sendBody = '{"isGlobal":2,"notifyCode":"ruleScene","subject":"direct test 2","body":"body test 2","userIDs":["' + $UserId + '"]}'
$send = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/message/info/send" -Headers $auth -Body $sendBody
Write-Host "send code=$($send.code) msg=$($send.msg)"

Start-Sleep -Seconds 2
$allAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5}}'
$ta = [int64]$allAfter.data.total
Write-Host "total after=$ta delta=$($ta-$tb)"
if ($allAfter.data.list) {
    foreach ($item in $allAfter.data.list) {
        Write-Host "group=$($item.group) notifyCode=$($item.notifyCode) subject=$($item.subject) body=$($item.body)"
    }
}

$st = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
Write-Host "stats:" ($st.data.list | ConvertTo-Json -Compress)
