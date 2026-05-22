# Google JWKS Fetcher

这个腾讯云 SCF 函数部署在新加坡，用来从 Google 官方 JWKS 地址获取公钥，并通过函数 URL 给国内生产服务缓存使用。

目标：

- 国内生产服务不直接访问 Google。
- 函数 URL 返回标准 JWKS 顶层 `keys`，并附带 `_meta` 调试信息。
- 后端遇到未知 `kid` 时应拒绝登录并告警，不要信任客户端传来的 `sub`。

部署区域：`ap-singapore`

函数名：`ykhl-google-jwks-fetcher`
