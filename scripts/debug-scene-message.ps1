# 场景消息为空 - 联犀 API 排查脚本
# 用法: 在 PowerShell 中执行:  .\scripts\debug-scene-message.ps1

$ErrorActionPreference = "Stop"
$BASE = "https://app.ykhl.vip"
$PROJECT = "1802965102490136576"
$ACCOUNT = "18059688688"
$PASSWORD_MD5 = "053a15b5f356b5f6e0e4c9a7b65e1b15"

function Invoke-Api {
    param(
        [string]$Uri,
        [hashtable]$Headers = @{},
        [string]$Body = "{}"
    )
    $h = @{ "Content-Type" = "application/json" }
    foreach ($k in $Headers.Keys) { $h[$k] = $Headers[$k] }
    try {
        $raw = Invoke-WebRequest -Uri $Uri -Method POST -Headers $h -Body $Body -UseBasicParsing
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
    param([string]$Token)
    return @{
        "ithings-token"      = $Token
        "ithings-project-id" = $PROJECT
        "ithings-app-code"   = "core"
    }
}

Write-Host "`n========== 1. 获取图形验证码 ==========" -ForegroundColor Cyan
$captchaBody = '{"type":"image","use":"login"}'
$captcha = Invoke-Api -Uri "$BASE/api/v1/system/user/self/captcha" -Body $captchaBody
if ($captcha.code -ne 200) {
    Write-Host "失败: $($captcha.msg)" -ForegroundColor Red
    exit 1
}
$codeID = $captcha.data.codeID
$b64 = $captcha.data.url -replace '^data:image/\w+;base64,', ''
$imgPath = Join-Path $env:TEMP "lianxi-captcha.png"
[IO.File]::WriteAllBytes($imgPath, [Convert]::FromBase64String($b64))
Start-Process $imgPath
Write-Host "验证码图片已打开: $imgPath"
Write-Host "codeID: $codeID"
$captchaCode = Read-Host "请输入图片中的验证码"

Write-Host "`n========== 2. 登录 ==========" -ForegroundColor Cyan
$loginObj = @{
    loginType = "pwd"
    account   = $ACCOUNT
    password  = $PASSWORD_MD5
    pwdType   = 2
    code      = $captchaCode
    codeID    = $codeID
}
$login = Invoke-Api -Uri "$BASE/api/v1/system/user/self/login" -Body ($loginObj | ConvertTo-Json)
if ($login.code -ne 200) {
    Write-Host "登录失败 code=$($login.code) msg=$($login.msg)" -ForegroundColor Red
    exit 1
}
$TOKEN = $login.data.token.accessToken
Write-Host "登录成功 user=$($login.data.info.userName) userID=$($login.data.info.userID)" -ForegroundColor Green

$auth = New-AuthHeaders -Token $TOKEN

Write-Host "`n========== 3. 当前用户 ==========" -ForegroundColor Cyan
$me = Invoke-Api -Uri "$BASE/api/v1/system/user/self/read" -Headers $auth -Body "{}"
Write-Host "  phone    = $($me.data.phone)"
Write-Host "  userName = $($me.data.userName)"
Write-Host "  userID   = $($me.data.userID)"

Write-Host "`n========== 4. 消息列表 ==========" -ForegroundColor Cyan
$m1 = Invoke-Api -Uri "$BASE/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":20}}'
Write-Host "  total(无时间) = $($m1.data.total)"
$m2 = Invoke-Api -Uri "$BASE/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":20},"createdTime":{"start":1779455441,"end":1780060241}}'
Write-Host "  total(有时间) = $($m2.data.total)"

Write-Host "`n========== 5. ruleScene + message 模板 ==========" -ForegroundColor Cyan
$tpl = Invoke-Api -Uri "$BASE/api/v1/system/notify/template/index" -Headers $auth -Body '{"notifyCode":"ruleScene","type":"message","page":{"page":1,"size":10}}'
$tplCount = 0
if ($tpl.data.list) { $tplCount = $tpl.data.list.Count }
Write-Host "  模板数量 = $tplCount"
if ($tplCount -eq 0) {
    Write-Host "  [!] 未配置站内信模板，场景触发后可能无法写入消息中心" -ForegroundColor Yellow
}

Write-Host "`n========== 6. 场景列表 & 手动触发 ==========" -ForegroundColor Cyan
$scenes = Invoke-Api -Uri "$BASE/api/v1/things/rule/scene/info/index" -Headers $auth -Body '{"page":{"page":1,"size":20},"name":"消息通知"}'
if (-not $scenes.data.list -or $scenes.data.list.Count -eq 0) {
    Write-Host "  未找到名为「消息通知」的场景" -ForegroundColor Yellow
} else {
    $scene = $scenes.data.list[0]
    Write-Host "  场景 id=$($scene.id) name=$($scene.name) status=$($scene.status)"
    $triggerBody = "{`"id`":`"$($scene.id)`"}"
    $tr = Invoke-Api -Uri "$BASE/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body $triggerBody
    if ($tr.code -eq 200) {
        Write-Host "  手动触发: 成功" -ForegroundColor Green
    } else {
        Write-Host "  手动触发: code=$($tr.code) msg=$($tr.msg)" -ForegroundColor Yellow
    }
    Start-Sleep -Seconds 3
    $m3 = Invoke-Api -Uri "$BASE/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":20}}'
    Write-Host "  触发后 total = $($m3.data.total)"
}

Write-Host "`n========== 完成 ==========" -ForegroundColor Cyan
Write-Host "Token 已保存在变量 `$TOKEN（本会话有效，约10分钟）"
