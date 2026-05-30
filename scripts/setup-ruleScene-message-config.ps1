# 在指定环境配置 ruleScene 站内信（notify 配置 + 模板 + 绑定）
param(
    [string]$BaseUrl = "https://new.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Account = "18059688688",
    [string]$PasswordMd5 = "053a15b5f356b5f6e0e4c9a7b65e1b15",
    [string]$AppCode = "client-app-android",
    [string]$TenantCode = "default",
    [Parameter(Mandatory = $true)][string]$CaptchaCode,
    [string]$CodeId = ""
)

$ErrorActionPreference = "Stop"

function Invoke-Api {
    param([string]$Uri, [hashtable]$Headers = @{}, [string]$Body = "{}")
    $h = @{ "Content-Type" = "application/json" }
    foreach ($k in $Headers.Keys) { $h[$k] = $Headers[$k] }
    $raw = Invoke-WebRequest -Uri $Uri -Method POST -Headers $h -Body $Body -UseBasicParsing -TimeoutSec 45
    return ($raw.Content | ConvertFrom-Json)
}

Write-Host "========== ENV $BaseUrl ==========" -ForegroundColor Cyan

$captcha = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/captcha" -Body '{"type":"image","use":"login"}'
if ($captcha.code -ne 200) { throw "captcha fail $($captcha.msg)" }
$useCodeId = $captcha.data.codeID
if ($CodeId -ne "") { $useCodeId = $CodeId }

$loginObj = @{
    loginType  = "pwd"
    tenantCode = $TenantCode
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

$auth = @{
    "ithings-token"      = $Token
    "ithings-project-id" = $ProjectId
    "ithings-app-code"   = $AppCode
}

Write-Host ""
Write-Host "========== NOTIFY CONFIG ruleScene ==========" -ForegroundColor Cyan
$cfg = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
if (-not $cfg.data.list -or $cfg.data.list.Count -eq 0) {
    throw "ruleScene notify config not found"
}
$c = $cfg.data.list[0]
Write-Host "id=$($c.id) supportTypes=$($c.supportTypes -join ',') enableTypes=$($c.enableTypes -join ',') isRecord=$($c.isRecord)"

$needUpdate = ($c.supportTypes -notcontains "message") -or ($c.enableTypes -notcontains "message") -or ($c.isRecord -ne 1)
if ($needUpdate) {
    $updBody = @{
        id           = $c.id
        code         = $c.code
        name         = $c.name
        group        = $c.group
        isRecord     = 1
        supportTypes = @("message", "sms", "email", "phoneCall", "dingTalk", "dingWebhook", "wxEWebHook", "wxMini")
        enableTypes  = @("message")
    } | ConvertTo-Json -Compress
    $upd = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/update" -Headers $auth -Body $updBody
    Write-Host "update config code=$($upd.code) msg=$($upd.msg)"
}

Write-Host ""
Write-Host "========== TEMPLATE ruleScene+message ==========" -ForegroundColor Cyan
$tpl = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/index" -Headers $auth -Body '{"notifyCode":"ruleScene","type":"message","page":{"page":1,"size":10}}'
$tplId = $null
if ($tpl.data.list -and $tpl.data.list.Count -gt 0) {
    $tplId = $tpl.data.list[0].id
    Write-Host "existing template id=$tplId subject=$($tpl.data.list[0].subject)"
} else {
    $createBody = '{"name":"场景联动站内信","notifyCode":"ruleScene","type":"message","code":"ruleScene_message","subject":"{{.title}}","body":"{{.body}}","desc":"场景联动站内信默认模板"}'
    $created = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/template/create" -Headers $auth -Body $createBody
    Write-Host "create template code=$($created.code) msg=$($created.msg) id=$($created.data.id)"
    if ($created.code -eq 200) { $tplId = $created.data.id }
}

if ($tplId) {
    $bindBody = '{"notifyCode":"ruleScene","type":"message","templateID":' + $tplId + '}'
    $bind = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/template/update" -Headers $auth -Body $bindBody
    Write-Host "bind template code=$($bind.code) msg=$($bind.msg)"
}

Write-Host ""
Write-Host "========== VERIFY ==========" -ForegroundColor Cyan
$cfg2 = Invoke-Api -Uri "$BaseUrl/api/v1/system/notify/config/index" -Headers $auth -Body '{"code":"ruleScene","page":{"page":1,"size":5}}'
$c2 = $cfg2.data.list[0]
Write-Host "supportTypes=$($c2.supportTypes -join ',') enableTypes=$($c2.enableTypes -join ',') isRecord=$($c2.isRecord)"
Write-Host "DONE" -ForegroundColor Green
