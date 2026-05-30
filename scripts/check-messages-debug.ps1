param(
    [string]$BaseUrl = "https://app.ykhl.vip",
    [string]$ProjectId = "1802965102490136576",
    [string]$Token = "",
    [string]$AppCode = "client-app-android"
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

Write-Host "=== ALL MESSAGES ==="
$m1 = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/index" -Headers $auth -Body '{"page":{"page":1,"size":20}}'
Write-Host "total=$($m1.data.total)"
if ($m1.data.list) {
    foreach ($item in $m1.data.list) {
        Write-Host "group=$($item.group) subject=$($item.subject) body=$($item.body)"
    }
}

Write-Host ""
Write-Host "=== STATISTICS ==="
$st = Invoke-Api -Uri "$BaseUrl/api/v1/system/user/self/message/statistics" -Headers $auth -Body "{}"
Write-Host ($st | ConvertTo-Json -Depth 5 -Compress)

Write-Host ""
Write-Host "=== SCENE 5551 READ ==="
$scene = Invoke-Api -Uri "$BaseUrl/api/v1/things/rule/scene/info/read" -Headers $auth -Body '{"id":5551}'
Write-Host ($scene.data.then | ConvertTo-Json -Depth 10 -Compress)
