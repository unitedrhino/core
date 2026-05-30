# 联调：用户场景「通知」- propertyReport + message notify
param(
    [string]$BaseUrl = "https://new.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Account = "18059688688",
    [string]$PasswordMd5 = "053a15b5f356b5f6e0e4c9a7b65e1b15",
    [string]$AppCode = "client-app-android",
    [string]$CaptchaCode = "",
    [string]$CodeId = "",
    [string]$Token = "",
    [string]$SceneName = "通知"
)

$ErrorActionPreference = "Stop"

function Invoke-Api {
    param([string]$Uri, [hashtable]$Headers = @{}, [string]$Body = "{}")
    $h = @{ "Content-Type" = "application/json" }
    foreach ($k in $Headers.Keys) { $h[$k] = $Headers[$k] }
    try {
        $raw = Invoke-WebRequest -Uri $Uri -Method POST -Headers $h -Body $Body -UseBasicParsing -TimeoutSec 45
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

Write-Host "========== ENV ==========" -ForegroundColor Cyan
Write-Host "BaseUrl=$BaseUrl ProjectId=$ProjectId Account=$Account SceneName=$SceneName"

if ($Token -eq "") {
    $captcha = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/captcha" -Body '{"type":"image","use":"login"}'
    if ($captcha.code -ne 200) { throw "captcha fail $($captcha.msg)" }
    $useCodeId = $captcha.data.codeID
    if ($CodeId -ne "") { $useCodeId = $CodeId }
    if ($CaptchaCode -eq "") {
        $b64 = $captcha.data.url -replace '^data:image/\w+;base64,', ''
        $imgPath = Join-Path $env:TEMP "lianxi-captcha-user.png"
        [IO.File]::WriteAllBytes($imgPath, [Convert]::FromBase64String($b64))
        Set-Content (Join-Path $env:TEMP "lianxi-captcha-codeid.txt") $useCodeId -NoNewline
        Start-Process $imgPath
        Write-Host "Captcha saved: $imgPath codeID=$useCodeId"
        throw "Need -CaptchaCode (image opened)"
    }
    $loginObj = @{
        loginType  = "pwd"
        tenantCode = "default"
        account    = $Account
        password   = $PasswordMd5
        pwdType    = 2
        code       = $CaptchaCode
        codeID     = $useCodeId
    }
    $login = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/login" -Body ($loginObj | ConvertTo-Json -Compress)
    if ($login.code -ne 200) { throw "login fail code=$($login.code) msg=$($login.msg)" }
    $Token = $login.data.token.accessToken
    Write-Host "LOGIN OK userID=$($login.data.info.userID)" -ForegroundColor Green
}

$auth = @{
    "ithings-token"      = $Token
    "ithings-project-id" = $ProjectId
    "ithings-app-code"   = $AppCode
}

Write-Host ""
Write-Host "========== NOTIFY CONFIG ruleScene ==========" -ForegroundColor Cyan
$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
if ($cfg.data.list -and $cfg.data.list.Count -gt 0) {
    $c = $cfg.data.list[0]
    Write-Host "supportTypes=$($c.supportTypes -join ',') enableTypes=$($c.enableTypes -join ',') isRecord=$($c.isRecord)"
    if ($c.supportTypes -notcontains "message") {
        Write-Host "[WARN] supportTypes missing message" -ForegroundColor Yellow
        if ($c.id) {
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
            Write-Host "update supportTypes code=$($upd.code) msg=$($upd.msg)"
        }
    }
}

Write-Host ""
Write-Host "========== TEMPLATE ruleScene+message ==========" -ForegroundColor Cyan
$tpl = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/index" -Headers $auth -Body '{"notifyCode":"ruleScene","type":"message","page":{"page":1,"size":10}}'
$tplId = $null
$tplCount = 0
if ($tpl.data.list) { $tplCount = $tpl.data.list.Count }
Write-Host "template count=$tplCount"
if ($tplCount -eq 0) {
    $createBody = '{"name":"scene notify message","notifyCode":"ruleScene","type":"message","code":"ruleScene_message","subject":"{{.title}}","body":"{{.body}}","desc":"debug user scene"}'
    $created = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/create" -Headers $auth -Body $createBody
    Write-Host "create template code=$($created.code) msg=$($created.msg) id=$($created.data.id)"
    if ($created.code -eq 200) { $tplId = $created.data.id }
} else {
    $tplId = $tpl.data.list[0].id
    Write-Host "existing template id=$tplId"
}
if ($tplId) {
    $bindBody = '{"notifyCode":"ruleScene","type":"message","templateID":' + $tplId + '}'
    $bind = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/template/update" -Headers $auth -Body $bindBody
    Write-Host "bind template code=$($bind.code) msg=$($bind.msg)"
}

Write-Host ""
Write-Host "========== FIND SCENE [$SceneName] ==========" -ForegroundColor Cyan
$scenes = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/index" -Headers $auth -Body ('{"page":{"page":1,"size":50},"name":"' + $SceneName + '"}')
$targetScene = $null
if ($scenes.data.list) {
    foreach ($s in $scenes.data.list) {
        if ($s.name -eq $SceneName) {
            $targetScene = $s
            break
        }
    }
    if ($null -eq $targetScene -and $scenes.data.list.Count -gt 0) {
        $targetScene = $scenes.data.list[0]
    }
}
if ($null -eq $targetScene) {
    Write-Host "Scene not found, skip trigger" -ForegroundColor Yellow
} else {
    $sid = $targetScene.id
    Write-Host "scene id=$sid name=$($targetScene.name) status=$($targetScene.status) type=$($targetScene.type)"
    $read = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/read" -Headers $auth -Body ('{"id":' + $sid + '}')
    if ($read.code -eq 200) {
        $thenJson = $read.data.then | ConvertTo-Json -Depth 15 -Compress
        Write-Host "then=$thenJson"
        $ifJson = $read.data.if
        if ($ifJson) {
            $ifStr = $ifJson | ConvertTo-Json -Depth 15 -Compress
            if ($ifStr.Length -gt 500) { $ifStr = $ifStr.Substring(0, 500) + "..." }
            Write-Host "if=$ifStr"
        }
    }

    Write-Host ""
    Write-Host "========== MESSAGE COUNT BEFORE ==========" -ForegroundColor Cyan
    $mBefore = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5},"notifyCode":"ruleScene"}'
    $totalBefore = 0
    if ($mBefore.data.total) { $totalBefore = [int]$mBefore.data.total }
    Write-Host "ruleScene messages total=$totalBefore"

    Write-Host ""
    Write-Host "========== MANUAL TRIGGER scene id=$sid ==========" -ForegroundColor Cyan
    $tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body ('{"id":' + $sid + '}')
    Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
    Start-Sleep -Seconds 4

    $mAfter = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5},"notifyCode":"ruleScene"}'
    $totalAfter = 0
    if ($mAfter.data.total) { $totalAfter = [int]$mAfter.data.total }
    Write-Host "ruleScene messages after=$totalAfter delta=$($totalAfter - $totalBefore)"
    if ($mAfter.data.list -and $mAfter.data.list.Count -gt 0) {
        $latest = $mAfter.data.list[0]
        Write-Host "latest subject=$($latest.subject)"
        Write-Host "latest body=$($latest.body)"
    }
    if ($totalAfter -gt $totalBefore) {
        Write-Host "SUCCESS: message recorded" -ForegroundColor Green
    } else {
        Write-Host "FAIL: no new ruleScene message after manual trigger" -ForegroundColor Yellow
        Write-Host "Possible: scene notify execute path issue OR notify.type not message in DB"
    }
}

Write-Host ""
Write-Host "========== DONE ==========" -ForegroundColor Cyan
Write-Host "Token=$Token"
