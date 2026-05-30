param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android",
    [string]$UserId = "1802965102482800640"
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

Write-Host "=== UPDATE ruleScene supportTypes add message ===" -ForegroundColor Cyan
$updBody = '{"code":"ruleScene","supportTypes":["message","sms","email","dingTalk","phoneCall"],"enableTypes":["message","sms"]}'
$upd = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/update" -Headers $auth -Body $updBody
Write-Host "update code=$($upd.code) msg=$($upd.msg)"

Write-Host ""
Write-Host "=== DIRECT SEND ruleScene message ===" -ForegroundColor Cyan
$groupJson = '{"page":{"page":1,"size":5},"group":"' + $MsgGroup + '"}'
$mBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$before = [int]$mBefore.data.total
$sendBody = '{"isGlobal":2,"notifyCode":"ruleScene","subject":"direct test","body":"direct notify send test","userIDs":["' + $UserId + '"]}'
$send = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/message/info/send" -Headers $auth -Body $sendBody
Write-Host "send code=$($send.code) msg=$($send.msg)"
Start-Sleep -Seconds 2
$mAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$after = [int]$mAfter.data.total
Write-Host "group messages before=$before after=$after delta=$($after-$before)"

Write-Host ""
Write-Host "=== TRIGGER SCENE 5551 AGAIN ===" -ForegroundColor Cyan
$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body '{"id":5551}'
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 3
$mFinal = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
Write-Host "group total=$($mFinal.data.total)"
if ($mFinal.data.list -and $mFinal.data.list.Count -gt 0) {
    Write-Host "latest subject=$($mFinal.data.list[0].subject) body=$($mFinal.data.list[0].body)"
}
