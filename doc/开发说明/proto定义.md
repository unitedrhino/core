# Protobuf定义规范

## 1. 基本结构

### 1.1 文件头部
```proto
syntax = "proto3";
option go_package = "pb/sys";
import "google/protobuf/wrappers.proto";

package sys;
```

## 2. Service定义规范

### 2.1 命名规范
- 采用大驼峰命名法（PascalCase）
- 服务名以`Manage`结尾，表示管理服务，如：`UserManage`、`RoleManage`、`TenantManage`
- 通用服务直接命名为`Common`
- 日志服务直接命名为`Log`

### 2.2 注释规范
- 每个服务必须添加注释，说明服务的功能和职责
- 注释应包含服务处理的主要业务范围
- 示例：
  ```proto
  // UserManage 用户管理服务
  // 处理用户CRUD、认证、授权等所有用户相关操作
  service UserManage {
      // ...
  }
  ```

### 2.3 服务分组
- 按业务功能模块划分服务，如用户管理、角色管理、租户管理等
- 每个服务负责一个明确的业务领域

### 2.4 示例
```proto
// RoleManage 角色管理服务
// 处理角色CRUD、角色权限分配、角色菜单权限管理等所有角色相关操作
service RoleManage {
    // ...
}

// TenantManage 租户管理服务
// 处理租户信息CRUD、租户配置、租户权限管理、租户开放接口等所有租户相关操作
service TenantManage {
    // ...
}
```

## 3. RPC接口定义规范

### 3.1 命名规范
- 采用小驼峰命名法（camelCase）
- 命名格式：`资源+操作`，如：`userInfoCreate`、`userInfoUpdate`、`userInfoRead`
- 资源名使用名词复数或单数形式，根据具体场景选择
- 操作名使用动词形式，如：`Create`、`Index`、`Update`、`Read`、`Delete`、`Login`等

### 3.2 请求/响应消息命名
- 请求消息名：`资源+操作+Req`，如：`UserInfoCreateReq`、`UserInfoReadReq`
- 响应消息名：`资源+操作+Resp`，如：`UserCreateResp`、`UserInfoIndexResp`
- 通用响应使用已定义的通用消息，如：`Empty`、`WithID`

### 3.3 注释规范
- 每个RPC方法必须添加注释，说明方法的功能
- 注释应简洁明了，使用中文
- 示例：`rpc userInfoCreate(UserInfoCreateReq) returns(UserCreateResp);//创建用户信息`

### 3.4 接口设计原则
- 每个接口应只完成一个明确的功能
- 避免设计过大的接口，拆分复杂功能为多个简单接口
- 遵循RESTful设计理念，使用合适的动词表示操作类型

### 3.5 示例
```proto
service UserManage {
  rpc userInfoCreate(UserInfoCreateReq) returns(UserCreateResp);//创建用户信息
  rpc userInfoIndex(UserInfoIndexReq) returns(UserInfoIndexResp);//获取用户列表（支持分页和过滤）
  rpc userInfoUpdate(userInfoUpdateReq) returns(Empty);//更新用户基本数据
  rpc userInfoRead(UserInfoReadReq) returns(UserInfo);//获取用户详细信息
  rpc userInfoDelete(UserInfoDeleteReq) returns(Empty);//删除用户
  rpc userLogin(UserLoginReq) returns(UserLoginResp);//用户登录
}
```

## 4. Message定义规范

### 4.1 命名规范
- 采用大驼峰命名法（PascalCase）
- 请求消息：`XXXReq`，如：`UserInfoCreateReq`、`UserLoginReq`
- 响应消息：`XXXResp`，如：`UserCreateResp`、`UserInfoIndexResp`
- 通用数据消息：直接使用实体名称，如：`UserInfo`、`RoleInfo`、`TenantInfo`
- 通用工具消息：使用描述性名称，如：`IDList`、`Attachment`、`SendOption`

### 4.2 结构设计
- 消息结构应与业务数据模型保持一致
- 使用嵌套结构体表达复杂数据关系
- 避免冗余字段，保持消息简洁

### 4.3 注释规范
- 每个消息必须添加注释，说明消息的用途
- 注释应包含消息的业务含义
- 示例：
  ```proto
  // UserInfoCreateReq 用户创建请求
  message UserInfoCreateReq {
      // ...
  }
  ```

### 4.4 通用消息设计
- 定义通用的响应消息，如：`Empty`、`WithID`、`WithIDCode`
- 通用消息应具有广泛的适用性

### 4.5 示例
```proto
// UserInfoCreateReq 用户创建请求
message UserInfoCreateReq {
    string userName = 1; // 登录用户名（必填）
    string password = 2; // 密码（必填）
    string phone = 3; // 手机号（可选）
    google.protobuf.StringValue email = 4; // 邮箱（可选）
}

// UserCreateResp 用户创建响应
message UserCreateResp {
    int64 userID = 1; // 用户ID
    string userName = 2; // 用户名
}

// UserInfo 用户信息
message UserInfo {
    int64 userID = 1; // 用户ID
    string userName = 2; // 用户名
    string nickName = 3; // 昵称
    string phone = 4; // 手机号
    google.protobuf.StringValue email = 5; // 邮箱
    int64 status = 6; // 状态（1:启用 2:禁用）
}

// IDList ID列表
// 用于批量操作的ID集合
message IDList{
    repeated int64 ids=1; //ID集合
}
```

## 5. 字段定义规范

### 5.1 命名规范
- 采用小驼峰命名法（camelCase）
- 字段名应准确表达字段的含义
- 避免使用缩写，除非是广为人知的缩写

### 5.2 字段编号
- 字段编号从1开始，依次递增
- 字段编号一旦分配，不应轻易修改
- 保留字段编号，避免重复使用

### 5.3 字段类型
- 使用proto3基础类型：`int64`、`string`、`bool`等
- 可选字段使用`google/protobuf/wrappers.proto`类型，如：`google.protobuf.StringValue`
- 重复字段使用`repeated`关键字
- 映射字段使用`map`关键字

### 5.4 字段注释
- 每个字段必须添加注释，说明字段含义、用途和约束条件
- 注释使用`//`开头，紧跟字段定义
- 注释应清晰明了，包含字段的业务含义
- 对于必填字段，应在注释中明确说明

### 5.5 布尔字段表示
- 使用`int64`类型表示布尔值
- `1`表示`true`，`2`表示`false`
- 示例：`int64 isGlobal = 2;//是否全局消息（1:是，2:否）`

### 5.6 示例
```proto
// NotifyConfig 通知配置信息
// 定义通知的基本配置，包括支持的类型、参数和模板
message NotifyConfig{
    int64 id =1; // 配置ID编号
    string group =2; //通知分组
    string code =3; // 通知类型编码
    string name =4; //通知名称
    repeated string supportTypes =5; //支持的通知类型（如sms、email、dingtalk等）
    string desc =6; // 配置备注说明
    int64 isRecord =7; //是否记录消息（1:是，2:否），是的情况下会将消息存一份到消息中心
    map<string,string> params =12; //通知变量属性（key是参数名，value是参数描述）
    repeated NotifyConfigTemplate templates =13;//绑定的模板信息列表
}
```

## 6. 最佳实践

### 6.1 版本兼容性
- 新增字段时，使用新的字段编号
- 不要删除字段，而是将其标记为`reserved`
- 不要修改现有字段的类型

### 6.2 性能考虑
- 合理使用字段类型，避免使用过大的类型
- 对于可选字段，使用`google/protobuf/wrappers.proto`类型
- 避免定义过大的消息结构

### 6.3 可读性
- 使用清晰的命名和注释
- 按逻辑顺序组织字段
- 使用空行分隔相关的字段组

### 6.4 一致性
- 保持命名风格的一致性
- 保持注释格式的一致性
- 保持服务、接口、消息设计的一致性

## 7. 通用消息定义

### 7.1 空响应
```proto
message Empty {
}
```

### 7.2 带ID响应
```proto
message WithID {
    int64 id = 1; // ID
}
```

### 7.3 带ID和Code响应
```proto
message WithIDCode {
    int64 id = 1; // ID
    string code = 2; // 编码
}
```

### 7.4 分页信息
```proto
message PageInfo {
    int64 page = 1; // 页码
    int64 page_size = 2; // 每页条数
    int64 total = 3; // 总条数
}
```
