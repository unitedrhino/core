param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android",
    [int64]$SceneId = 5551
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

function Get-SceneNotifyCount {
    $st = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
    foreach ($item in $st.data.list) {
        if ($item.notifyCode -eq "ruleScene" -or $item.group -like "*ruleScene*") {
            return [int]$item.count
        }
    }
    # fallback: second group often scene notify after direct send
    foreach ($item in $st.data.list) {
        if ($item.count -lt 100 -and $item.group -ne "feedBack") {
            Write-Host "stat group=$($item.group) count=$($item.count)"
            return [int]$item.count
        }
    }
    return 0
}

$before = Get-SceneNotifyCount
Write-Host "scene-notify-like count before=$before"
$tr = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/manually-trigger" -Headers $auth -Body ('{"id":' + $SceneId + '}')
Write-Host "trigger code=$($tr.code) msg=$($tr.msg)"
Start-Sleep -Seconds 3
$after = Get-SceneNotifyCount
Write-Host "count after=$after delta=$($after-$before)"

$recent = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":5},"notifyCode":"ruleScene"}'
Write-Host "recent ruleScene messages:"
if ($recent.data.list) {
    foreach ($item in $recent.data.list) {
        Write-Host "  subject=$($item.subject) body=$($item.body)"
    }
}
