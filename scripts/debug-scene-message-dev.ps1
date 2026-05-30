# Scene linkage in-app message integration debug (demo dev: new.ykhl.vip)
# Usage: powershell -ExecutionPolicy Bypass -File .\scripts\debug-scene-message-dev.ps1

param(
    [string]$BaseUrl = "https://new.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Account = "18059688688",
    [string]$PasswordMd5 = "053a15b5f356b5f6e0e4c9a7b65e1b15",
    [string]$AppCode = "client-app-android",
    [string]$CaptchaCode = "",
    [string]$CodeId = "",
    [string]$Token = ""
)

$ErrorActionPreference = "Stop"
$MsgGroup = [string]::Concat([char]0x573A, [char]0x666F, [char]0x8054, [char]0x52A8, [char]0x901A, [char]0x77E5)

function Invoke-Api {
    param(
        [string]$Uri,
        [hashtable]$Headers = @{},
        [string]$Body = "{}"
    )
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

function New-AuthHeaders {
    param([string]$AccessToken)
    return @{
        "ithings-token"      = $AccessToken
        "ithings-project-id" = $ProjectId
        "ithings-app-code"   = $AppCode
    }
}

Write-Host ""
Write-Host "========== ENV ==========" -ForegroundColor Cyan
Write-Host "  BaseUrl   = $BaseUrl"
Write-Host "  ProjectId = $ProjectId"
Write-Host "  AppCode   = $AppCode"
Write-Host "  Account   = $Account"

Write-Host ""
Write-Host "========== 1. CAPTCHA ==========" -ForegroundColor Cyan
$captcha = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/captcha" -Body '{"type":"image","use":"login"}'
if ($captcha.code -ne 200) {
    Write-Host "  FAIL code=$($captcha.code) msg=$($captcha.msg)" -ForegroundColor Red
    exit 1
}
Write-Host "  OK codeID=$($captcha.data.codeID)" -ForegroundColor Green

if ($Token -eq "") {
    $useCodeId = $captcha.data.codeID
    if ($CodeId -ne "" -and $CodeId -ne "skip") { $useCodeId = $CodeId }
    if ($CaptchaCode -eq "") {
        $b64 = $captcha.data.url -replace '^data:image/\w+;base64,', ''
        $imgPath = Join-Path $env:TEMP "lianxi-captcha-dev.png"
        [IO.File]::WriteAllBytes($imgPath, [Convert]::FromBase64String($b64))
        Start-Process $imgPath
        Write-Host "  Captcha image: $imgPath"
        $CaptchaCode = Read-Host "Enter captcha"
    }
    Write-Host ""
    Write-Host "========== 2. LOGIN ==========" -ForegroundColor Cyan
    $loginObj = @{
        loginType = "pwd"
        account   = $Account
        password  = $PasswordMd5
        pwdType   = 2
        code      = $CaptchaCode
        codeID    = $useCodeId
    }
    $login = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/login" -Body ($loginObj | ConvertTo-Json -Compress)
    if ($login.code -ne 200) {
        Write-Host "  FAIL code=$($login.code) msg=$($login.msg)" -ForegroundColor Red
        exit 1
    }
    $Token = $login.data.token.accessToken
    Write-Host "  OK userID=$($login.data.info.userID)" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "========== 2. USE EXISTING TOKEN ==========" -ForegroundColor Cyan
}

$auth = New-AuthHeaders -AccessToken $Token

Write-Host ""
Write-Host "========== 3. USER ==========" -ForegroundColor Cyan
$me = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/read" -Headers $auth -Body "{}"
Write-Host "  phone=$($me.data.phone) userID=$($me.data.userID)"

Write-Host ""
Write-Host "========== 4. NOTIFY CONFIG ruleScene ==========" -ForegroundColor Cyan
$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
if ($cfg.data.list -and $cfg.data.list.Count -gt 0) {
    $c = $cfg.data.list[0]
    Write-Host "  name=$($c.name) isRecord=$($c.isRecord)"
    Write-Host "  supportTypes=$($c.supportTypes -join ',')"
    Write-Host "  enableTypes=$($c.enableTypes -join ',')"
} else {
    Write-Host "  ruleScene config not found" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========== 5. TEMPLATE ruleScene+message ==========" -ForegroundColor Cyan
$tpl = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/index" -Headers $auth -Body '{"notifyCode":"ruleScene","type":"message","page":{"page":1,"size":10}}'
$tplCount = 0
$tplId = $null
if ($tpl.data.list) { $tplCount = $tpl.data.list.Count }
Write-Host "  template count = $tplCount"
if ($tplCount -eq 0) {
    Write-Host "  creating template..." -ForegroundColor Yellow
    $subjectTpl = '{{.title}}'
    $bodyTpl = '{{.body}}'
    $createBody = '{"name":"ruleScene message","notifyCode":"ruleScene","type":"message","code":"ruleScene_message","subject":"' + $subjectTpl + '","body":"' + $bodyTpl + '","desc":"debug auto create"}'
    $created = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/create" -Headers $auth -Body $createBody
    if ($created.code -eq 200) {
        $tplId = $created.data.id
        Write-Host "  created id=$tplId" -ForegroundColor Green
    } else {
        Write-Host "  create FAIL code=$($created.code) msg=$($created.msg)" -ForegroundColor Red
    }
} else {
    $tplId = $tpl.data.list[0].id
    Write-Host "  existing id=$tplId"
}

if ($tplId) {
    Write-Host ""
    Write-Host "========== 6. BIND TEMPLATE ==========" -ForegroundColor Cyan
    $bindBody = '{"notifyCode":"ruleScene","type":"message","templateID":' + $tplId + '}'
    $bind = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/template/update" -Headers $auth -Body $bindBody
    if ($bind.code -eq 200) {
        Write-Host "  bind OK" -ForegroundColor Green
    } else {
        Write-Host "  bind code=$($bind.code) msg=$($bind.msg)" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "========== 7. MESSAGES BEFORE ==========" -ForegroundColor Cyan
$groupJson = '{"page":{"page":1,"size":20},"group":"' + $MsgGroup + '"}'
$mBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
$totalBefore = 0
if ($mBefore.data.total) { $totalBefore = [int]$mBefore.data.total }
Write-Host "  total = $totalBefore"

Write-Host ""
Write-Host "========== 8. SCENE TRIGGER ==========" -ForegroundColor Cyan
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
    Write-Host "  no scene with notify/ruleScene - create one in APP first" -ForegroundColor Yellow
} else {
    Write-Host "  scene id=$($targetScene.id) name=$($targetScene.name)"
    $triggerBody = '{"id":"' + $targetScene.id + '"}'
    $tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body $triggerBody
    if ($tr.code -eq 200) {
        Write-Host "  trigger OK" -ForegroundColor Green
    } else {
        Write-Host "  trigger code=$($tr.code) msg=$($tr.msg)" -ForegroundColor Yellow
    }
    Start-Sleep -Seconds 3
    $mAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body $groupJson
    $totalAfter = 0
    if ($mAfter.data.total) { $totalAfter = [int]$mAfter.data.total }
    Write-Host "  total after = $totalAfter (delta = $($totalAfter - $totalBefore))"
    if ($totalAfter -gt $totalBefore) {
        Write-Host "  SUCCESS: message recorded" -ForegroundColor Green
    } else {
        Write-Host "  WARN: no new message" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "========== DONE ==========" -ForegroundColor Cyan
Write-Host ('Token=' + $Token)
