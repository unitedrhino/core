# RPC Logic 层开发规范

## 1. 概述

Logic 层是 RPC 服务的核心业务逻辑处理层，负责实现具体的业务功能，包括参数校验、事务处理、数据访问、错误处理等。Logic 层通过依赖注入获取所需的服务和资源，与数据访问层（Repo）和缓存层交互，完成业务逻辑的处理。

## 2. 文件结构与命名规范

### 2.1 目录结构

Logic 层采用模块化设计，按照业务领域划分为不同的子目录，每个模块包含相关的业务逻辑和数据组装函数：

```
internal/logic/
├── accessmanage/      # 权限管理模块
├── appmanage/         # 应用管理模块
├── areamanage/        # 区域管理模块
├── common/            # 通用功能模块
├── datamanage/        # 数据管理模块
├── departmentmanage/  # 部门管理模块
├── dictmanage/        # 字典管理模块
```

### 2.2 命名规范

| 类型         | 命名格式           | 示例                                  |
| ------------ | ------------------ | ------------------------------------- |
| Logic 文件   | 小驼峰+Logic 后缀  | `userInfoCreateLogic.go`            |
| 数据组装文件 | 统一命名           | `assemble.go`                       |
| 通用逻辑文件 | 功能命名           | `auth.go`、`menu.go`、`user.go` |
| Logic 结构体 | 大驼峰+Logic 后缀  | `UserInfoCreateLogic`               |
| 构造函数     | New+结构体名       | `NewUserInfoCreateLogic`            |
| 核心方法     | RPC 接口名         | `UserInfoCreate`                    |
| 辅助方法     | 小驼峰             | `getPwd`、`GetUserInfo`           |
| 数据转换函数 | 源类型+To+目标类型 | `UserInfoToPb`                      |
| 通用校验函数 | Check+功能         | `CheckUserName`、`CheckPwd`       |

## 3. 代码结构规范

### 3.1 Logic 结构体定义

```go
// userInfoCreateLogic.go
type UserInfoCreateLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext  // 由 gozero 自动生成，包含所有依赖服务
    logx.Logger
    UiDB *relationDB.UserInfoRepo  // 数据访问层实例
}
```

**规范说明**：

- 必须包含 `ctx`、`svcCtx` 和 `logx.Logger` 字段
- `svcCtx` 由 gozero 自动生成，包含数据库、缓存、OSS客户端等依赖服务
- 数据访问层实例（Repo）作为结构体字段直接声明，便于依赖注入
- 字段命名采用小驼峰命名法，结构体名采用大驼峰命名法

### 3.2 构造函数实现

```go
// userInfoCreateLogic.go
func NewUserInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoCreateLogic {
    return &UserInfoCreateLogic{
        ctx:    ctx,
        svcCtx: svcCtx,  // 直接使用 gozero 自动生成的服务上下文
        Logger: logx.WithContext(ctx),
        UiDB:   relationDB.NewUserInfoRepo(ctx),
    }
}
```

**规范说明**：

- 构造函数接收 `ctx` 和 `svcCtx` 参数
- `svcCtx` 由 gozero 自动生成，直接传入使用
- 使用 `logx.WithContext(ctx)` 创建带上下文的日志实例
- 初始化数据访问层实例并注入到结构体中

### 3.3 核心方法实现

```go
// userInfoCreateLogic.go
func (l *UserInfoCreateLogic) UserInfoCreate(in *sys.UserInfoCreateReq) (*sys.UserCreateResp, error) {
    l.Infof("%s req=%+v", utils.FuncName(), in)  // 记录请求日志
  
    // 业务逻辑实现
  
    return &sys.UserCreateResp{UserID: userID}, nil
}
```

**规范说明**：

- 方法名与 RPC 接口名保持一致
- 入参为对应的请求消息结构体，出参为响应消息结构体和错误
- 使用 `l.Infof` 记录请求日志，便于问题追踪
- 使用 `utils.FuncName()` 获取当前函数名

## 4. 业务逻辑实现规范

### 4.1 参数校验

**规范说明**：

- 对所有外部输入参数进行严格校验
- 使用正则表达式、格式检查等方式确保参数合法性
- 校验失败时返回 `errors.Parameter` 类型错误

**示例代码**：

```go
// 用户名格式校验
if err := logic.CheckUserName(in.UserName); err != nil {
    return nil, err
}

// 手机号格式校验
if !utils.IsPhone(in.Phone) {
    return nil, errors.Parameter.AddMsgf("手机号格式错误")
}
```

### 4.2 事务管理

**规范说明**：

- 涉及多个数据操作的业务逻辑必须使用事务
- 使用 `stores.GetCommonConn(l.ctx).Transaction` 或 `stores.GetTenantConn(l.ctx).Transaction` 封装事务
- 事务内的所有数据操作必须使用同一个事务对象
- 事务结束后根据错误决定提交或回滚

**示例代码**：

```go
err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
    uidb := relationDB.NewUserInfoRepo(tx)  // 使用事务对象创建Repo
  
    // 1. 检查用户是否已存在
    ui, err := uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}})
    if err != nil && !errors.Cmp(err, errors.NotFind) {
        return err
    }
  
    // 2. 创建用户或更新信息（根据业务逻辑）
    // ... 业务逻辑实现 ...
  
    // 3. 关联角色信息
    return relationDB.NewUserRoleRepo(tx).MultiUpdate(l.ctx, ui.UserID, in.RoleIDs)
})
```

### 4.3 错误处理

**规范说明**：

- 使用 `errors` 包进行错误分类和处理
- 错误类型包括：`errors.Parameter`（参数错误）、`errors.Permissions`（权限错误）、`errors.NotFind`（未找到）、`errors.System`（系统错误）等
- 使用 `errors.AddMsgf` 或 `errors.AddDetail` 添加错误详情
- 使用 `errors.Cmp` 比较错误类型

**示例代码**：

```go
// 参数错误
if !utils.IsPhone(in.Phone) {
    return nil, errors.Parameter.AddMsgf("手机号格式错误")
}

// 权限错误
if uc == nil {
    return nil, errors.Permissions.WithMsg("无操作权限")
}

// 系统错误
if err != nil {
    return nil, errors.System.AddDetail(err)
}
```

### 4.4 缓存使用

**规范说明**：

- 合理使用缓存提高系统性能
- 通过 `l.svcCtx` 使用 gozero 自动注入的缓存实例
- 缓存键名使用统一命名规范：`sys:user:wxak:{prefix}:{code}`
- 设置合理的缓存过期时间

**示例代码**：

```go
// 获取租户配置缓存（通过 svcCtx 使用 gozero 自动注入的缓存实例）
tc, err := l.svcCtx.TenantConfigCache.GetData(l.ctx, uc.TenantCode)
if err != nil {
    return nil, err
}
```

### 4.5 事件发布

**规范说明**：

- 使用事件驱动架构实现模块解耦
- 通过 `l.svcCtx.FastEvent` 发布事件（gozero 自动注入）
- 事件主题定义在 `topics` 包中

**示例代码**：

```go
// 用户创建事件发布
e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{userID}})
if e != nil {
    l.Errorf("Publish event err:%v", e)
}
```

### 4.6 数据组装

**规范说明**：

- 系统级通用数据组装放在根目录下的 `assemble.go` 文件中
- 模块内数据组装放在各模块目录下的 `assemble.go` 文件中
- 命名采用「源类型+To+目标类型」格式，如 `UserInfoToPb`、`ToPageInfo`
- 使用 `utils.Copy` 进行结构体复制，提高代码复用性
- 处理特殊字段，如文件URL签名、坐标转换等

**示例代码**：

```go
// UserInfoToPb 用户信息转换为PB结构体
func UserInfoToPb(ctx context.Context, ui *relationDB.SysUserInfo, svcCtx *svc.ServiceContext) *sys.UserInfo {
    // 处理头像URL签名
    if ui.Avatar != "" {
        ui.Avatar, _ = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, ui.Avatar, 24*60*60, common.OptionKv{})
    }
    return utils.Copy[sys.UserInfo](ui)
}
```

### 4.7 通用逻辑复用

**规范说明**：

- 跨模块的通用逻辑封装在根目录下的通用逻辑文件中，如 `user.go`、`auth.go`、`menu.go`
- 通用逻辑函数命名采用「Check+功能」或「Get+功能」格式
- 通用逻辑需考虑可复用性，避免耦合特定业务逻辑

**示例代码**：

```go
// CheckUserName 校验用户名格式（通用逻辑）
func CheckUserName(userName string) error {
    if ret, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]{6,19}$", userName); !ret {
        return errors.UsernameFormatErr.AddDetail("账号格式错误")
    }
    return nil
}
```

## 5. 依赖注入规范

**规范说明**：

- 通过 `svcCtx` 传递所有依赖服务
- 依赖服务包括：数据库连接、缓存实例、OSS客户端、事件总线等
- 避免在Logic层直接创建依赖实例

**示例代码**：

```go
// 使用svcCtx中的服务（gozero自动注入）
userID := l.svcCtx.UserID.GetSnowflakeId()  // 雪花ID生成器
e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{userID}})  // 事件发布
tc, err := l.svcCtx.TenantCache.GetData(l.ctx, uc.TenantCode)  // 获取租户缓存
```

## 6. 日志记录规范

**规范说明**：

- 使用 `logx.WithContext(ctx)` 创建带上下文的日志实例
- 记录请求参数、关键操作和错误信息
- 使用 `utils.FuncName()` 获取当前函数名

**示例代码**：

```go
// 记录请求日志
l.Infof("%s req=%+v", utils.FuncName(), in)

// 记录错误日志
if err != nil {
    l.Errorf("%s err:%v", utils.FuncName(), err)
}
```

## 7. 安全规范

### 7.1 密码处理

**规范说明**：

- 密码必须加密存储
- 使用 `utils.MakePwd` 进行密码加密
- 支持不同的密码类型和加密方式

**示例代码**：

```go
// 密码加密
password := utils.MakePwd(in.Password, userID, false)
```

### 7.2 权限控制

**规范说明**：

- 基于租户和角色的权限控制
- 使用上下文获取当前用户信息
- 验证用户是否有权限执行操作

**示例代码**：

```go
// 获取用户上下文
uc := ctxs.GetUserCtx(l.ctx)
if uc == nil {
    return 0, errors.Permissions.WithMsg("无租户号")
}

// 验证租户应用状态
ta, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(ctxs.CommonWithDefault(l.ctx),
    relationDB.TenantAppFilter{AppID: l.uc.AppID, WithApp: true})
if err != nil {
    return nil, err
}
if ta.Status != app.TenantStatusNormal {
    return nil, errors.AppCantUse
}
```

## 8. 最佳实践

### 8.1 代码复用

- 将通用逻辑提取到单独的函数或文件中
- 数据转换函数统一放在 `assemble.go` 文件中
- 使用工具函数库 `utils` 处理通用操作

### 8.2 性能优化

- 合理使用缓存减少数据库查询
- 批量操作减少网络开销
- 避免不必要的数据加载

### 8.3 可维护性

- 代码结构清晰，职责单一
- 充分的日志记录便于问题追踪
- 错误处理统一规范

## 9. 代码示例

### 9.1 用户创建逻辑

```go
func (l *UserInfoCreateLogic) UserInfoCreate(in *sys.UserInfoCreateReq) (*sys.UserCreateResp, error) {
    l.Infof("%s req=%+v", utils.FuncName(), in)
  
    // 参数校验
    if err := l.validateParams(in); err != nil {
        return nil, err
    }
  
    // 事务处理创建用户
    var userID int64
    err := stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
        var err error
        userID, err = l.createUser(in, tx)
        return err
    })
  
    if err != nil {
        return nil, err
    }
  
    // 发布用户创建事件
    if e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{userID}}); e != nil {
        l.Errorf("%s publish event err:%v", utils.FuncName(), e)
    }
  
    return &sys.UserCreateResp{UserID: userID}, nil
}
```

### 9.2 用户登录逻辑

```go
// 登录逻辑结构体定义
type LoginLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext  // 由 gozero 自动生成，包含所有依赖服务
    logx.Logger
    UiDB *relationDB.UserInfoRepo  // 用户信息数据访问层
}

// 构造函数
func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
    return &LoginLogic{
        ctx:    ctx,
        svcCtx: svcCtx,  // 使用 gozero 自动生成的服务上下文
        Logger: logx.WithContext(ctx),
        UiDB:   relationDB.NewUserInfoRepo(ctx),
    }
}

// 登录主函数
func (l *LoginLogic) UserLogin(in *sys.UserLoginReq) (*sys.UserLoginResp, error) {
    l.Infof("%s req=%+v", utils.FuncName(), in)
  
    // 参数校验
    if err := l.validateLoginParams(in); err != nil {
        return nil, err
    }
  
    // 获取用户信息
    userInfo, err := l.getUserInfo(in)
    if err != nil {
        return nil, err
    }
  
    // 生成登录响应（使用svcCtx中的JWT生成器）
    return l.genLoginResponse(userInfo)
}

// 验证登录参数
func (l *LoginLogic) validateLoginParams(in *sys.UserLoginReq) error {
    // 登录参数校验逻辑
    return nil
}

// 获取用户信息
func (l *LoginLogic) getUserInfo(in *sys.UserLoginReq) (*relationDB.SysUserInfo, error) {
    // 根据登录方式获取用户信息
    return l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}})
}

// 生成登录响应
func (l *LoginLogic) genLoginResponse(userInfo *relationDB.SysUserInfo) (*sys.UserLoginResp, error) {
    // 生成JWT令牌并返回响应
    return &sys.UserLoginResp{
        Info: UserInfoToPb(l.ctx, userInfo, l.svcCtx),
        // JWT令牌生成逻辑
    }, nil
}
```

## 10. 总结

Logic 层是 RPC 服务的核心业务处理层，遵循上述规范可以提高代码的可读性、可维护性和可扩展性。在实际开发中，应根据具体业务需求灵活应用这些规范，并不断优化和改进代码质量。
