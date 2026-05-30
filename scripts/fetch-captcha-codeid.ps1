param([string]$BaseUrl = "https://new.ykhl.vip")
$ErrorActionPreference = "Stop"
$raw = Invoke-WebRequest -Uri "$BaseUrl/api/v1/system/user/self/captcha" -Method POST -ContentType "application/json" -Body '{"type":"image","use":"login"}' -UseBasicParsing -TimeoutSec 30
$j = $raw.Content | ConvertFrom-Json
if ($j.code -ne 200) { Write-Error $j.msg; exit 1 }
$codeId = $j.data.codeID
$path = Join-Path $env:TEMP "lianxi-captcha-codeid.txt"
Set-Content -Path $path -Value $codeId -NoNewline
$b64 = $j.data.url -replace '^data:image/\w+;base64,', ''
$imgPath = Join-Path $env:TEMP "lianxi-captcha-dev.png"
[IO.File]::WriteAllBytes($imgPath, [Convert]::FromBase64String($b64))
Write-Output $codeId
