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
