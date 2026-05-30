param(
    [Parameter(Mandatory=$true)][string]$CaptchaCode,
    [string]$BaseUrl = "https://new.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Account = "18059688688",
    [string]$PasswordMd5 = "053a15b5f356b5f6e0e4c9a7b65e1b15",
    [string]$AppCode = "client-app-android"
)

$ErrorActionPreference = "Stop"
$MsgGroup = [string]::Concat([char]0x573A, [char]0x666F, [char]0x8054, [char]0x52A8, [char]0x901A, [char]0x77E5)

function Invoke-Api {
    param([string]$Uri, [hashtable]$Headers = @{}, [string]$Body = "{}")
    $h = @{ "Content-Type" = "application/json" }
    foreach ($k in $Headers.Keys) { $h[$k] = $Headers[$k] }
    try {
        $raw = Invoke-WebRequest -Uri $Uri -Method POST -Headers $h -Body $Body -UseBasicParsing -TimeoutSec 30
        return ($raw.Content | ConvertFrom-Json)
    } catch {
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $errBody = $reader.ReadToEnd()
            $reader.Close()
            try { return ($errBody | ConvertFrom-Json) } catch { throw $errBody }
        }
        throw
    }
}

Write-Host "Fetching captcha..."
$captcha = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/captcha" -Body '{"type":"image","use":"login"}'
if ($captcha.code -ne 200) { throw "captcha fail $($captcha.msg)" }
$CodeId = $captcha.data.codeID
Write-Host "codeID=$CodeId captcha=$CaptchaCode"

$loginObj = @{
    loginType  = "pwd"
    tenantCode = "default"
    account    = $Account
    password   = $PasswordMd5
    pwdType    = 2
    code       = $CaptchaCode
    codeID     = $CodeId
}
$login = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/login" -Body ($loginObj | ConvertTo-Json -Compress)
if ($login.code -ne 200) {
    Write-Host "LOGIN FAIL code=$($login.code) msg=$($login.msg)" -ForegroundColor Red
    exit 1
}
$Token = $login.data.token.accessToken
Write-Host "LOGIN OK userID=$($login.data.info.userID)" -ForegroundColor Green

$auth = @{
    "ithings-token"      = $Token
    "ithings-project-id" = $ProjectId
    "ithings-app-code"   = $AppCode
}

Write-Host ""
Write-Host "=== NOTIFY CONFIG ruleScene ===" -ForegroundColor Cyan
$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
if ($cfg.data.list -and $cfg.data.list.Count -gt 0) {
    $c = $cfg.data.list[0]
    Write-Host "supportTypes=$($c.supportTypes -join ',') enableTypes=$($c.enableTypes -join ',') isRecord=$($c.isRecord)"
}

Write-Host ""
Write-Host "=== TEMPLATE ruleScene+message ===" -ForegroundColor Cyan
$tpl = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/index" -Headers $auth -Body '{"notifyCode":"ruleScene","type":"message","page":{"page":1,"size":10}}'
$tplId = $null
$tplCount = 0
if ($tpl.data.list) { $tplCount = $tpl.data.list.Count }
Write-Host "count=$tplCount"
if ($tplCount -eq 0) {
    $createBody = '{"name":"ruleScene message","notifyCode":"ruleScene","type":"message","code":"ruleScene_message","subject":"{{.title}}","body":"{{.body}}","desc":"debug"}'
    $created = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/create" -Headers $auth -Body $createBody
    if ($created.code -eq 200) {
        $tplId = $created.data.id
        Write-Host "created template id=$tplId" -ForegroundColor Green
    } else {
        Write-Host "create fail code=$($created.code) msg=$($created.msg)" -ForegroundColor Red
    }
} else {
    $tplId = $tpl.data.list[0].id
    Write-Host "existing id=$tplId subject=$($tpl.data.list[0].subject)"
}

if ($tplId) {
    $bindBody = '{"notifyCode":"ruleScene","type":"message","templateID":' + $tplId + '}'
    $bind = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/template/update" -Headers $auth -Body $bindBody
    Write-Host "bind code=$($bind.code) msg=$($bind.msg)"
}

Write-Host ""
Write-Host "=== MESSAGES BEFORE ===" -ForegroundColor Cyan
$groupJson = '{"page":{"page":1,"size":20},"group":"' + $MsgGroup + '"}'
$mBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$totalBefore = 0
if ($mBefore.data.total) { $totalBefore = [int]$mBefore.data.total }
Write-Host "total=$totalBefore"

Write-Host ""
Write-Host "=== SCENE TRIGGER ===" -ForegroundColor Cyan
$scenes = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/index" -Headers $auth -Body '{"page":{"page":1,"size":30}}'
$targetScene = $null
if ($scenes.data.list) {
    foreach ($s in $scenes.data.list) {
        $thenJson = $s.then | ConvertTo-Json -Depth 20 -Compress
        if ($thenJson -match 'notify' -and $thenJson -match 'ruleScene') {
            $targetScene = $s
            break
        }
    }
}
if ($null -eq $targetScene) {
    Write-Host "No scene with notify/ruleScene" -ForegroundColor Yellow
    if ($scenes.data.list) {
        Write-Host "Available scenes:"
        foreach ($s in $scenes.data.list) {
            Write-Host "  id=$($s.id) name=$($s.name)"
        }
    }
} else {
    Write-Host "scene id=$($targetScene.id) name=$($targetScene.name)"
    $triggerBody = '{"id":"' + $targetScene.id + '"}'
    $tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body $triggerBody
    Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
    Start-Sleep -Seconds 3
    $mAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
    $totalAfter = 0
    if ($mAfter.data.total) { $totalAfter = [int]$mAfter.data.total }
    Write-Host "total after=$totalAfter delta=$($totalAfter - $totalBefore)"
    if ($totalAfter -gt $totalBefore) {
        Write-Host "SUCCESS - message recorded" -ForegroundColor Green
        if ($mAfter.data.list -and $mAfter.data.list.Count -gt 0) {
            Write-Host "latest subject=$($mAfter.data.list[0].subject)"
            Write-Host "latest body=$($mAfter.data.list[0].body)"
        }
    } else {
        Write-Host "No new message - check scene then.actions notify config" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "Token=$Token"
