# API文件定义规范

## 1. 目录结构

API文件应按照业务模块进行组织，遵循以下目录结构：

```
service/apisvr/http/
├── 模块名称/          # 如 system
│   ├── 子模块名称/    # 如 user
│   │   ├── 功能.api   # 如 info.api
│   │   └── ...
│   └── 模块.api       # 如 user.api
├── common.api         # 公共结构体定义
└── 其他模块/          # 如 device、product 等
```

**示例**：
```
service/apisvr/http/
├── system/
│   ├── user/
│   │   ├── info.api    # 用户信息管理
│   │   ├── role.api    # 角色管理
│   │   └── dept.api    # 部门管理
│   ├── tenant/
│   │   ├── tenant.api  # 租户管理
│   │   └── app.api     # 应用管理
│   └── system.api      # 系统模块入口
└── common.api          # 公共结构体定义
```

## 2. 文件命名规范

- **模块入口文件**：使用模块名称命名，如 `user.api`
- **功能文件**：使用功能名称命名，如 `info.api`、`role.api`
- **公共文件**：命名为 `common.api`，存放公共结构体定义

## 3. API服务定义

### 3.1 服务元信息

每个API文件必须包含服务元信息定义：

```api
info(
    title: "模块名称"            # 模块名称
    desc: "模块描述信息"          # 模块详细描述
    author: "作者"               # 作者名称
    email: "作者邮箱"             # 作者邮箱
    version: "版本号"            # 版本号，如 v0.1.0
)
```

### 3.2 服务配置

使用 `@server` 注解配置服务：

```api
@server (
    group: 模块/子模块/功能        # 如 system/user/info
    prefix: /api/v1/模块/子模块/功能  # 如 /api/v1/system/user/info
    accessCodePrefix: "权限前缀"    # 如 "systemUserManage"
    accessNamePrefix: "权限名称前缀"  # 如 "用户信息"
    accessGroup: "权限组名称"       # 如 "用户管理"
    defaultAuthType: "认证类型"      # 如 "admin"
    defaultNeedAuth: "是否需要认证"   # 如 "true"
    middleware: 中间件1,中间件2     # 如 CheckTokenWare,InitCtxsWare
)
```

### 3.3 接口定义

使用 `service` 定义API接口：

```api
service api {
    @doc "接口描述"          # 接口文档
    @handler 处理器名称        # 如 create
    请求方法 /路径 (请求参数) returns (返回参数)
}
```

**示例**：
```api
service api {
    @doc "创建用户信息"
    @handler create
    post /create (UserInfoCreateReq) returns (UserCreateResp)

    @doc "查询用户信息列表"
    @handler Index
    post /index (UserInfoIndexReq) returns (UserInfoIndexResp)

    @doc "获取用户信息"
    @handler read
    post /read (UserInfoReadReq) returns (UserInfo)
}
```

### 3.4 文件上传接口定义

当需要上传文件时，可以使用 `@doc` 注解的 `injectFormdataParam` 属性来支持 formdata 格式上传：

```api
service api {
    @doc(
        summary: "文件直传"
        injectFormdataParam: "file"
    )
    @handler uploadFile
    post /upload-file(Empty) returns (UploadFileResp)
}
```

**说明**：
- `injectFormdataParam: "file"` 表示该接口支持 formdata 格式上传，文件字段名为 `file`
- 请求参数通常使用 `Empty` 结构体，文件数据通过 formdata 传递
- 服务器端可以通过 `context` 获取上传的文件数据

## 4. 结构体定义

### 4.1 定义格式

使用 `type` 关键字定义结构体：

```api
type (
    结构体名称 {
        字段名 类型 `json:"json字段名,属性1,属性2"` // 字段描述
    }
)
```

### 4.2 字段定义规范

#### 4.2.1 命名规范

- **结构体名称**：使用大驼峰命名法，如 `UserInfo`、`UserInfoCreateReq`
- **字段名称**：使用小驼峰命名法，与Proto文件保持一致
- **JSON字段名**：使用小驼峰命名法，与字段名保持一致

#### 4.2.2 类型规范

- **基本类型**：`int64`、`string`、`bool` 等
- **指针类型**：用于可选字段，如 `*string`
- **数组类型**：如 `[]int64`、`[]string`
- **Map类型**：如 `map[string]string`
- **嵌套结构体**：如 `*UserInfo`

#### 4.2.3 属性规范

| 属性 | 说明 | 示例 |
|------|------|------|
| optional | 可选字段 | `json:"userName,optional"` |
| string | 序列化时转为字符串 | `json:"userID,string"` |
| omitempty | 空值时不序列化 | `json:"password,optional,omitempty"` |
| range | 数值范围限制 | `json:"roleIDs,optional,range=(0:120]"` |

**示例**：
```api
type (
    UserInfo {
        UserID      int64            `json:"userID,string,optional,omitempty"`        // 用户id
        UserName    string           `json:"userName,optional,omitempty"`             // 登录用户名
        NickName    string           `json:"nickName,optional,omitempty"`             // 用户的昵称
        Password    string           `json:"password,optional,omitempty"`             // 登录密码
        Email       *string          `json:"email,optional,omitempty"`                 // 邮箱
        Phone       *string          `json:"phone,optional,omitempty"`                 // 手机号
        Avatar      string           `json:"avatar,optional,omitempty"`                // 头像
        Tenants     []*UserTenant    `json:"tenants,optional,omitempty"`               // 租户信息列表
        Thirds      []*UserThird     `json:"thirds,optional,omitempty"`                // 第三方绑定信息
    }
)
```

### 4.3 与Proto文件的一致性

API接口定义的结构体字段必须与对应的Proto文件保持一致，包括：

- 字段名称
- 字段类型
- 字段是否必填
- 字段顺序（建议保持一致，便于维护）

**示例**：

Proto文件定义：
```proto
message UserInfo {
    int64 userID = 1;
    string userName = 2;
    string nickName = 3;
    string password = 4;
    string email = 5;
    string phone = 6;
    string avatar = 7;
}
```

API文件定义：
```api
UserInfo {
    UserID      int64    `json:"userID,string,optional,omitempty"`        // 用户id
    UserName    string   `json:"userName,optional,omitempty"`             // 登录用户名
    NickName    string   `json:"nickName,optional,omitempty"`             // 用户的昵称
    Password    string   `json:"password,optional,omitempty"`             // 登录密码
    Email       *string  `json:"email,optional,omitempty"`                 // 邮箱
    Phone       *string  `json:"phone,optional,omitempty"`                 // 手机号
    Avatar      string   `json:"avatar,optional,omitempty"`                // 头像
}
```

## 5. WithXxx 模式

使用 `WithXxx` 模式实现按需加载关联数据，提高接口性能。

### 5.1 定义方式

在请求结构体中添加 `WithXxx` 类型的布尔字段，用于控制是否返回关联数据。

**示例**：
```api
UserInfoIndexReq {
    Page *PageInfo `json:"page,optional"` //分页信息
    UserName    string `json:"userName,optional"`          //用户名
    Phone       string `json:"phone,optional"`             // 手机号
    Email       string `json:"email,optional"`             // 邮箱
    WithRoles   bool   `json:"withRoles,optional"`         // 同时返回角色列表
    WithDepts   bool   `json:"withDepts,optional"`         // 同时返回部门列表
}

UserInfoReadReq {
    UserID      int64  `json:"userID,string,optional"`     // 用户id
    WithRoles   bool   `json:"withRoles,optional"`         // 同时返回角色列表
    WithTenant  bool   `json:"withTenant,optional"`        // 同时返回租户信息
    WithDepts   bool   `json:"withDepts,optional"`         // 同时返回部门列表
}
```

### 5.2 使用场景

- 当接口需要返回主数据及可选的关联数据时
- 当关联数据较大或查询耗时较长时
- 当不同客户端对返回数据有不同需求时

### 5.3 实现要求

- 在API层解析 `WithXxx` 参数
- 在业务逻辑层根据 `WithXxx` 参数决定是否加载关联数据
- 在数据组装层根据 `WithXxx` 参数决定是否返回关联数据

## 6. 公共结构体

公共结构体应定义在 `common.api` 文件中，供所有API文件共享使用。

**示例**：
```api
type (
    PageInfo {
        Page     int64 `json:"page,optional"`         // 页码
        Size     int64 `json:"size,optional"`         // 每页大小
        Orders   []*OrderBy `json:"orders,optional"`  //排序
    }
    PageResp {
        Page     int64 `json:"page,optional"`         // 页码
        PageSize int64 `json:"pageSize,optional"`     // 每页大小
        Total    int64 `json:"total"`                  // 总条数
    }
    OrderBy {
        Field string `json:"field,optional"`          // 排序字段名
        Sort  int64  `json:"sort,optional"`           // 排序方式：0 从小到大, 1 从大到小
    }
    WithID {
        ID int64 `json:"id,optional"`                  // id
    }
    WithCode {
        Code string `json:"code,optional"`             // 编码
    }
    Empty {
        // 空结构体，用于无请求参数或无返回参数的接口
    }
)
```

## 7. 服务新增流程

### 7.1 创建RPC服务

```shell
goctl rpc new 服务名称  --style=goZero -m
```

**示例**：
```shell
goctl rpc new opssvr  --style=goZero -m
```

### 7.2 创建API服务

```shell
goctl api new 服务名称  --style=goZero 
```

**示例**：
```shell
goctl api new viewsvr  --style=goZero 
```

## 8. 代码生成

使用以下命令生成API代码：

```shell
cd 服务目录 && goctl api go -api http/接口文件  -dir ./  --style=goZero -ws
```

**示例**：
```shell
cd apisvr && goctl api go -api http/api.api  -dir ./  --style=goZero -ws
```

## 9. 最佳实践

### 9.1 接口设计原则

- **单一职责**：每个接口只负责一个功能
- **按需返回**：使用 `WithXxx` 模式实现按需加载
- **版本控制**：在URL中包含版本号，如 `/api/v1/`
- **统一命名**：接口名称、参数名称、返回字段名称保持一致的命名规范

### 9.2 性能优化

- 避免在接口中返回大量数据
- 使用分页查询处理列表数据
- 合理使用缓存减少数据库查询
- 避免N+1查询问题

### 9.3 安全性

- 对敏感接口添加认证和授权
- 对输入参数进行校验
- 避免返回敏感信息
- 使用HTTPS协议传输数据

### 9.4 文件上传

- 使用 `@doc` 注解的 `injectFormdataParam` 属性实现文件上传
- 文件上传接口的请求参数通常使用 `Empty` 结构体
- 服务器端需对上传文件大小、类型进行限制
- 上传成功后返回文件路径或访问URL

## 10. 完整示例

### 10.1 用户信息API文件

```api
info(
    title: "用户管理模块"
    desc: "用户管理相关接口，包括创建账号，登录，获取验证码，获取用户列表，获取单个用户信息，更新用户信息，删除用户"
    author: "L"
    email: "174805676@qq.com"
    version: "v0.1.0"
)

@server (
    group: system/user/info
    prefix: /api/v1/system/user/info
    accessCodePrefix: "systemUserManage"
    accessNamePrefix: "用户信息"
    accessGroup: "用户管理"
    defaultAuthType: "admin"
    defaultNeedAuth: "true"
    middleware:  CheckTokenWare,InitCtxsWare
)

service api {
    @doc "创建用户信息"
    @handler create
    post /create (UserInfoCreateReq) returns (UserCreateResp)

    @doc "查询用户信息列表"
    @handler Index
    post /index (UserInfoIndexReq) returns (UserInfoIndexResp)

    @doc "获取用户信息"
    @handler read
    post /read (UserInfoReadReq) returns (UserInfo)
}

type (
    UserInfo {
        UserID      int64            `json:"userID,string,optional,omitempty"`        // 用户id
        UserName    string           `json:"userName,optional,omitempty"`             // 登录用户名
        NickName    string           `json:"nickName,optional,omitempty"`             // 用户的昵称
        Password    string           `json:"password,optional,omitempty"`             // 登录密码
        Email       *string          `json:"email,optional,omitempty"`                 // 邮箱
        Phone       *string          `json:"phone,optional,omitempty"`                 // 手机号
        Avatar      string           `json:"avatar,optional,omitempty"`                // 头像
        Tenants     []*UserTenant    `json:"tenants,optional,omitempty"`               // 租户信息列表
    }

    UserInfoCreateReq {
        Info    *UserInfo `json:"info"`
        RoleIDs []int64   `json:"roleIDs,optional,range=(0:120]"` //角色编号列表
    }

    UserInfoIndexReq {
        Page     *PageInfo `json:"page,optional"` //分页信息
        UserName string     `json:"userName,optional"`          //用户名
        Phone    string     `json:"phone,optional"`             // 手机号
        Email    string     `json:"email,optional"`             // 邮箱
        WithRoles bool      `json:"withRoles,optional"`         // 同时返回角色列表
        WithDepts bool      `json:"withDepts,optional"`         // 同时返回部门列表
    }

    UserInfoIndexResp {
        List []*UserInfo `json:"list,omitempty"`           //用户信息列表
        PageResp
    }

    UserInfoReadReq {
        UserID    int64 `json:"userID,string,optional"`     // 用户id
        WithRoles bool  `json:"withRoles,optional"`         // 同时返回角色列表
        WithTenant bool `json:"withTenant,optional"`        // 同时返回租户信息
    }

    UserCreateResp {
        UserID int64 `json:"userID,string,optional"`        // 用户id
    }
)
```

## 11. 注意事项

- API文件的字段定义必须与Proto文件保持一致
- 使用 `WithXxx` 模式实现按需加载，提高接口性能
- 公共结构体应放在 `common.api` 文件中
- 接口命名应清晰、简洁，表达接口的功能
- 为每个接口添加详细的文档说明
- 遵循统一的命名规范和目录结构
