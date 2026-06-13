# Core 仓库

> SaaS 中台核心模块，负责用户、角色、项目、应用、通知、字典等核心能力。

## 仓库信息

| 项目 | 值 |
|------|-----|
| 远程 | `git@codeup.aliyun.com:642f7dca8b47795dae985084/ee/core.git` |
| 当前分支 | `dev`（最新） |
| 模块 | `gitee.com/unitedrhino/core` |
| Go 版本 | 1.24.4 |

## 关键目录

```
core/
├── service/
│   ├── apisvr/          # API 网关（HTTP 入口）
│   ├── syssvr/          # 系统管理 RPC（用户、角色、项目、菜单）
│   ├── datasvr/         # 数据分析与动态查询
│   └── timed/           # 定时任务相关服务
│       ├── timedjobsvr/
│       └── timedschedulersvr/
└── Makefile
```

## 常用操作

```bash
# 编译全部
cd core && make build

# 仅编译 API 网关
cd core && make build.api

# 编译并上传到 47.94.112.109
cd core && make packback

# 查看远程
git remote -v
```

## 依赖升级

升级 share 依赖后执行 `go mod tidy`，然后打 tag：
```bash
git tag v1.5.x
git push origin v1.5.x
git push gitee v1.5.x
```

## 部署踩坑记录

### 部署 coresvr 禁止覆盖远端 sys.yaml（2026-06-11）

- 远端 `47.94.112.109:/root/run/core/etc/sys.yaml` 含手工维护的 `UniPush` 配置（`Enabled: true` + Secret，与云函数 `YK_PUSH_HTTP_SECRET` 一致）；仓库默认值是 `Enabled: false`、`Secret: ""`
- 2026-06-11 14:38 一次部署把远端 sys.yaml 覆盖成仓库默认值，导致系统通知（systemNotice → uni-push）全部静默不发，且无任何报错日志（`shouldSendAppPush` 直接返回 false）
- `scripts/deploy-apisvr-remote.mjs` 默认不上传 sys.yaml；除非明确要更新配置，否则不要设置 `DEPLOY_UPLOAD_SYSYAML=1`
- 误覆盖后的恢复：运行 `scripts/fix-unipush-config.mjs`（写回 Enabled/Secret 并重启 coresvr）
- 排查口诀：场景执行了、消息中心有记录、但通知栏没推送 → 先查远端 sys.yaml 的 `UniPush.Enabled`
