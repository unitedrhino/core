# 场景联动 + 站内信通知端到端联调
param(
    [string]$BaseUrl = "https://new.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Account = "18059688688",
    [string]$PasswordMd5 = "053a15b5f356b5f6e0e4c9a7b65e1b15",
    [string]$AppCode = "client-app-android",
    [string]$TenantCode = "default",
    [Parameter(Mandatory = $true)][string]$CaptchaCode,
    [string]$CodeId = "",
    [string]$SceneName = "通知",
    [int64]$SceneId = 0,
    [string]$ProductId = "008",
    [string]$DeviceName = "C8586A7C4613",
    [string]$PropertyDataId = "rf_btn1.2",
    [string]$PropertyValue = "3"
)

$ErrorActionPreference = "Stop"
$MsgGroup = [string]::Concat([char]0x573A, [char]0x666F, [char]0x8054, [char]0x52A8, [char]0x901A, [char]0x77E5)

function Invoke-Api {
    param([string]$Uri, [hashtable]$Headers = @{}, [string]$Body = "{}")
    $h = @{ "Content-Type" = "application/json" }
    foreach ($k in $Headers.Keys) { $h[$k] = $Headers[$k] }
    $raw = Invoke-WebRequest -Uri $Uri -Method POST -Headers $h -Body $Body -UseBasicParsing -TimeoutSec 60
    return ($raw.Content | ConvertFrom-Json)
}

Write-Host "========== LOGIN $BaseUrl ==========" -ForegroundColor Cyan
$captcha = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/captcha" -Body '{"type":"image","use":"login"}'
$useCodeId = $captcha.data.codeID
if ($CodeId -ne "") { $useCodeId = $CodeId }
$login = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/login" -Body (@{
    loginType = "pwd"; tenantCode = $TenantCode; account = $Account; password = $PasswordMd5
    pwdType = 2; code = $CaptchaCode; codeID = $useCodeId
} | ConvertTo-Json -Compress)
if ($login.code -ne 200) { throw "login fail code=$($login.code) msg=$($login.msg)" }
$Token = $login.data.token.accessToken
$UserId = $login.data.info.userID
Write-Host "OK userID=$UserId" -ForegroundColor Green

$auth = @{
    "ithings-token"      = $Token
    "ithings-project-id" = $ProjectId
    "ithings-app-code"   = $AppCode
}

Write-Host ""
Write-Host "========== NOTIFY CONFIG ==========" -ForegroundColor Cyan
$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
$c = $cfg.data.list[0]
Write-Host "supportTypes=$($c.supportTypes -join ',') enableTypes=$($c.enableTypes -join ',') isRecord=$($c.isRecord)"

Write-Host ""
Write-Host "========== FIND SCENE ==========" -ForegroundColor Cyan
$sid = $SceneId
if ($sid -eq 0) {
    $scenes = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/index" -Headers $auth -Body '{"page":{"page":1,"size":50}}'
    foreach ($s in $scenes.data.list) {
        if ($s.name -eq $SceneName) { $sid = [int64]$s.id; break }
    }
    if ($sid -eq 0 -and $scenes.data.list.Count -gt 0) {
        foreach ($s in $scenes.data.list) {
            if ($s.type -eq "auto") { $sid = [int64]$s.id; $SceneName = $s.name; break }
        }
    }
}
if ($sid -eq 0) { throw "scene not found" }
Write-Host "scene id=$sid name=$SceneName"

$read = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/read" -Headers $auth -Body ('{"id":' + $sid + '}')
$thenObj = $read.data.then | ConvertFrom-Json
$notify = $thenObj.actions | Where-Object { $_.type -eq "notify" } | Select-Object -First 1
if ($notify) {
    Write-Host "notify type=$($notify.notify.type) code=$($notify.notify.notifyCode) accounts=$($notify.notify.accounts -join ',')"
}

Write-Host ""
Write-Host "========== MESSAGES BEFORE ==========" -ForegroundColor Cyan
$groupJson = '{"page":{"page":1,"size":5},"group":"' + $MsgGroup + '"}'
$allBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5}}'
$beforeTotal = [int64]$allBefore.data.total
Write-Host "all messages total=$beforeTotal"

Write-Host ""
Write-Host "========== MANUAL TRIGGER scene $sid ==========" -ForegroundColor Cyan
$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body ('{"id":' + $sid + '}')
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 5
$allAfterTrigger = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5}}'
$afterTriggerTotal = [int64]$allAfterTrigger.data.total
Write-Host "after manual trigger total=$afterTriggerTotal delta=$($afterTriggerTotal - $beforeTotal)"

Write-Host ""
Write-Host "========== PROPERTY CONTROL simulate trigger ==========" -ForegroundColor Cyan
$dataObj = @{}
$dataObj[$PropertyDataId] = $PropertyValue
$pcBody = @{
    productID     = $ProductId
    deviceName    = $DeviceName
    data          = ($dataObj | ConvertTo-Json -Compress)
    syncTimeout   = 20
    shadowControl = 0
} | ConvertTo-Json -Compress
$pc = Invoke-Api -Uri "$BaseUrl/api/v1/things/device/interact/property-control-send" -Headers $auth -Body $pcBody
Write-Host "property-control code=$($pc.code) msg=$($pc.msg)"
Start-Sleep -Seconds 8
$allAfterProp = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5}}'
$afterPropTotal = [int64]$allAfterProp.data.total
Write-Host "after property total=$afterPropTotal delta=$($afterPropTotal - $beforeTotal)"
if ($allAfterProp.data.list) {
    foreach ($m in $allAfterProp.data.list) {
        Write-Host "latest id=$($m.id) group=$($m.group) notifyCode=$($m.notifyCode) subject=$($m.subject) body=$($m.body)"
    }
}

Write-Host ""
if ($afterPropTotal -gt $beforeTotal -or $afterTriggerTotal -gt $beforeTotal) {
    Write-Host "SUCCESS: new message(s) recorded" -ForegroundColor Green
} else {
    Write-Host "FAIL: no new message (NotifyM fix may not be deployed on thingsEEsvr)" -ForegroundColor Yellow
}
Write-Host "DONE" -ForegroundColor Cyan
