# 系统通知（systemNotice）手动测试文档

> **API 域名**：`https://new.ykhl.vip`（以实际环境为准）  
> **App**：含 uni-push 2.0 的自定义基座或正式包  
> **说明**：全程通过 **管理端 UI、App 操作、Postman、SSH 日志、SQL** 完成，**不使用任何 Node 脚本**。

---

## 一、测试目标

| 序号 | 能力 | 通过标准 |
|------|------|----------|
| 1 | 通知授权时机 | 安装/登录后**不**弹系统通知授权；仅在场景里勾选「系统通知」并点确定后弹出 |
| 2 | 管理端配置 | `ruleScene` 已启用 **系统推送（systemNotice）** 并绑定模板 |
| 3 | cid 上报 | 用户允许通知后，服务端能查到 `push_client_id` |
| 4 | API 直推 | 测试下发接口能在手机通知栏收到推送 |
| 5 | 场景自动化 | 触发场景后通知栏收到推送，点击可跳转 App |
| 6 | 触发用户信息 | 消息中心返回 `triggerUser`；推送 payload 含触发人相关字段 |

---

## 二、测试前准备

### 2.1 账号与设备

| 项 | 要求 | □ |
|----|------|---|
| 测试手机 | Android 真机（建议 13+），可选 iOS | |
| App | 最新自定义基座/正式包 | |
| 账号 A | 设备主人，手机号与场景通知账号一致 | |
| 账号 B | （可选）被分享设备账号 | |
| 管理端 | 可登录 my-things-admin | |
| 测试设备 | 已绑定项目，属性可上报触发场景 | |

### 2.2 服务端（人工确认）

| 项 | 期望 | □ |
|----|------|---|
| coresvr | 进程正常，API 可访问 | |
| thingsEEsvr | 进程正常 | |
| UniPush | `Enabled: true`，Secret 与云函数一致 | |
| 云函数 yk-push-send | 已部署，鉴权密钥已配置 | |

### 2.3 工具

- 浏览器（管理端 + F12 复制 Token）
- Postman / Apifox
- （可选）SSH、数据库客户端

---

## 三、测试流程概览

```
阶段 1 → App 安装与隐私（不应提前要通知权限）
阶段 2 → 管理端消息配置
阶段 3 → App 授权 + cid 上报
阶段 4 → API 直推（隔离场景链路）
阶段 5 → 场景自动化端到端
阶段 6 → 触发用户信息 triggerUser
阶段 7 → 回归与异常
```

按 **1→2→3→4→5→6→7** 顺序执行；**4 通过后再做 5**。

---

## 四、阶段 1 — App 安装与隐私

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 1.1 | 卸载旧包，安装最新 App | 安装成功 | |
| 1.2 | 冷启动 | 出现隐私协议 | |
| 1.3 | 同意隐私前 | **不出现**系统「允许通知」弹窗 | |
| 1.4 | 同意隐私、进入登录 | 仍 **不出现**通知授权 | |
| 1.5 | 登录进入首页 Tab | 仍 **不出现**通知授权 | |

---

## 五、阶段 2 — 管理端消息配置

**入口**：客户管理 → 消息通知

### 5.1 消息模板

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 2.1 | 打开消息模板列表 | 存在 `ruleScene` + 通道 **系统推送 / systemNotice** | |
| 2.2 | 查看模板内容 | 主题/正文含 `{{.title}}`、`{{.body}}` | |

### 5.2 消息配置

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 2.3 | 消息配置 → **场景联动通知** | **系统推送** 已勾选 | |
| 2.4 | 绑定模板 | 已绑定上述 systemNotice 模板 | |
| 2.5 | 保存 | 成功无报错 | |

### 5.3 数据库核对（可选）

```sql
SELECT code, enable_types FROM sys_notify_config WHERE code = 'ruleScene';

SELECT notify_code, type, template_id
FROM sys_notify_config_template
WHERE notify_code = 'ruleScene';
```

`enable_types` 须含 `"systemNotice"`。

---

## 六、阶段 3 — App 授权与 cid 上报

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 3.1 | App 登录账号 A | 成功 | |
| 3.2 | 场景联动 → 新建/编辑 → 执行动作 → **通知** | 弹出通知方式选择 | |
| 3.3 | 勾选 **系统通知**，点 **确定** | 弹出系统通知授权框 | |
| 3.4 | 点 **允许** | App 正常，无崩溃 | |
| 3.5 | 杀进程重开 | **不再**重复弹授权 | |

### cid 验收（SQL，手机号改实际值）

```sql
SELECT user_id, push_client_id, platform, updated_at
FROM sys_user_push_client
WHERE user_id = (
  SELECT user_id FROM sys_user_info
  WHERE phone = '你的测试手机号' LIMIT 1
)
ORDER BY updated_at DESC;
```

期望：`push_client_id` 非空。

### cid 验收（Postman，可选）

```http
POST https://new.ykhl.vip/api/v1/system/user/self/push-client/report
iThings-set-token: <App 登录 token>
Content-Type: application/json

{
  "pushClientId": "<cid>",
  "platform": "android",
  "appId": "__UNI__F82AD01",
  "appVersion": "2.3.7"
}
```

---

## 七、阶段 4 — API 直推

**目的**：不经过场景，验证推送链路。

### 7.1 准备

1. 管理端 F12 → 复制 **`iThings-set-token`**
2. 查测试用户 **`user_id`**（用户管理或 SQL）

### 7.2 请求

```http
POST https://new.ykhl.vip/api/v1/system/notify/config/test/send
Content-Type: application/json
iThings-set-token: <管理员 token>

{
  "notifyCode": "ruleScene",
  "type": "systemNotice",
  "userIDs": ["<user_id>"],
  "params": {
    "title": "联调-API直推",
    "body": "systemNotice 手动测试",
    "triggerType": "manual",
    "triggerUserId": "<user_id>"
  }
}
```

### 7.3 验收

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 4.1 | 发送请求 | HTTP 200 | |
| 4.2 | 手机桌面/锁屏 | 通知栏有标题和正文 | |
| 4.3 | 点击通知 | App 打开，跳转场景 Tab | |

### 失败看日志（SSH）

```bash
grep -iE 'unipush|pushAppNotify|systemNotice' /tmp/coresvr.log | tail -20
```

| 日志 | 含义 |
|------|------|
| `no active push client ids` | 回到阶段 3，cid 未上报 |
| `NotEnable` | 管理端未启用 systemNotice |
| `pushAppNotify err` | Secret / 云函数 / 网络 |

---

## 八、阶段 5 — 场景自动化

### 8.1 创建场景（App 或管理端）

| 字段 | 示例 |
|------|------|
| 名称 | 系统通知测试 |
| 类型 | auto 或 manual |
| 触发 | 设备属性变化（与物模型一致） |
| 执行 | 通知 → 勾选 **系统通知**（可同时勾 **消息中心**） |
| notifyCode | ruleScene |
| 通知账号 | 测试账号手机号 |

保存并 **启用**。

### 8.2 执行验收

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 5.1 | 确认阶段 3、4 已通过 | — | |
| 5.2 | 手动执行 或 设备触发 | 数秒内通知栏收到推送 | |
| 5.3 | 点击通知 | 跳转 App 场景页 | |
| 5.4 | 仅勾 systemNotice | 消息中心无新记录（正常） | |
| 5.5 | 仅勾 message | 有站内信，无通知栏 | |

### 失败排查

```bash
# 场景是否执行
grep -iE 'SceneExec|NotifyConfigSend' /root/run/things-ee/thingsEEsvr.log | tail -30

# 是否推送
grep -iE 'pushAppNotify|NotEnable' /tmp/coresvr.log | tail -30
```

```sql
SELECT id, name, LEFT(`then`, 400) FROM uds_scene_info
WHERE name LIKE '%系统通知测试%' ORDER BY id DESC LIMIT 1;
```

---

## 九、阶段 6 — 触发用户信息（triggerUser）

### 9.1 手动触发 + 消息中心

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 6.1 | 场景勾选 **消息中心**（或双选） | — | |
| 6.2 | 账号 A **手动执行**场景 | — | |
| 6.3 | App → 消息中心 → **通知** Tab | 副标题：**由 xxx 手动触发** | |

**Postman 查消息列表**：

```http
POST https://new.ykhl.vip/api/v1/system/user/self/message/index
iThings-set-token: <App token>
Content-Type: application/json

{
  "page": { "page": 1, "size": 20 },
  "group": "场景联动通知"
}
```

期望（新消息）：

```json
"triggerUser": {
  "userId": "...",
  "nickName": "张三",
  "account": "13800138000",
  "triggerType": "manual"
}
```

### 9.2 自动触发

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 6.4 | 设备/定时自动触发（勾 message） | 副标题：**系统自动触发** | |
| 6.5 | API 响应 | `triggerType": "auto"` | |

### 9.3 仅 systemNotice

| 步骤 | 操作 | 期望 | □ |
|------|------|------|---|
| 6.6 | 只勾系统通知并触发 | 消息中心无记录；推送 payload 含 trigger 字段 | |

**SQL（可选）**：

```sql
SELECT id, trigger_user_id, trigger_user_nick, trigger_type, LEFT(subject, 40)
FROM sys_message_info
WHERE notify_code = 'ruleScene'
ORDER BY id DESC LIMIT 5;
```

---

## 十、阶段 7 — 回归与异常

| 编号 | 场景 | 期望 | □ |
|------|------|------|---|
| 7.1 | 双选 message + systemNotice | 通知栏 + 消息中心都有 | |
| 7.2 | 拒绝通知权限 | 不崩溃，场景可保存 | |
| 7.3 | 杀 App 后设备触发 | 可能仍收到离线推送（视机型） | |
| 7.4 | App 前台触发 | 可收到推送/本地补通知 | |
| 7.5 | 分享账号 B 手动执行 | triggerUser 为 B；按策略通知主人 | |
| 7.6 | 升级前旧消息 | 无 triggerUser 时不显示副标题、不报错 | |

---

## 十一、测试结果记录

| 阶段 | 执行人 | 日期 | 结果 | 备注 |
|------|--------|------|------|------|
| 1 隐私 | | | □通过 □失败 | |
| 2 管理端 | | | □通过 □失败 | |
| 3 cid | | | □通过 □失败 | |
| 4 API 直推 | | | □通过 □失败 | |
| 5 场景 | | | □通过 □失败 | |
| 6 triggerUser | | | □通过 □失败 | |
| 7 回归 | | | □通过 □失败 | |

失败请记录：机型、系统版本、App 版本、user_id、场景 ID、Postman 响应、日志片段。

---

## 十二、常见问题

| 现象 | 检查 |
|------|------|
| 安装就弹通知权限 | 是否旧包；是否未用「仅场景触发」版本 |
| API 200 无通知 | cid 是否上报；UniPush 是否 Enabled |
| 场景不执行 | 触发条件/物模型；边缘是否已触发过 |
| 场景执行无推送 | 是否勾 systemNotice；enable_types 配置 |
| 无 triggerUser | 是否新消息；是否勾 message 类型 |
| 仅 systemNotice 无站内信 | 设计如此，非缺陷 |

---

## 十三、接口速查

| 用途 | 方法 | 路径 |
|------|------|------|
| 测试下发 | POST | `/api/v1/system/notify/config/test/send` |
| 消息列表 | POST | `/api/v1/system/user/self/message/index` |
| 上报 cid | POST | `/api/v1/system/user/self/push-client/report` |

鉴权头：`iThings-set-token: <token>`

---

**版本**：2026-06-03 · 含 systemNotice 全链路 + triggerUser + App 延迟授权
