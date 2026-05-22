# Google JWKS Fetcher

这个腾讯云 SCF 函数部署在新加坡，用来从 Google 官方 JWKS 地址获取公钥，并通过自定义域名给国内生产服务缓存使用。

目标：

- 国内生产服务不直接访问 Google。
- 函数自定义域名保持公网访问，但必须携带 `X-YKHL-JWKS-TOKEN` 共享密钥。
- 函数自定义域名返回标准 JWKS 顶层 `keys`，并附带 `_meta` 调试信息。
- 后端遇到未知 `kid` 时应拒绝登录并告警，不要信任客户端传来的 `sub`。

部署区域：`ap-singapore`

函数名：`ykhl-google-jwks-fetcher`

自定义域名：`https://sgtick.ykhl.vip`

运行环境变量：

- `YKHL_JWKS_ACCESS_TOKEN`：函数自定义域名访问密钥，不提交仓库。

后端配置：

- `GOOGLE_JWKS_URL=https://sgtick.ykhl.vip`
- `GOOGLE_JWKS_HEADER_NAME=X-YKHL-JWKS-TOKEN`
- `GOOGLE_JWKS_ACCESS_TOKEN=<与 SCF 环境变量一致>`

监控告警：

- 通知模板：`notice-vrs6sz8n`（杨磊告警）
- 调用次数策略：`policy-hvkh55p9`，5 分钟 `Invocation > 1000` 触发。
- 外网出流量策略：`policy-xwvl0y2s`，5 分钟 `OutFlow > 10240 Kb` 触发。
- 绑定对象：`sg/default/ykhl-google-jwks-fetcher/$LATEST`。
