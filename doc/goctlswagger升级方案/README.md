# goctl Swagger 升级方案：2.0 → OpenAPI 3.0

## 背景

当前 goctl 的 `api swagger` 命令生成 Swagger 2.0 格式文档，许多现代 API 工具（如 Apifox、Swagger UI 4.x、Redoc 等）和前端代码生成器对 OpenAPI 3.0 支持更好。本方案将输出格式从 Swagger 2.0 完全升级为 OpenAPI 3.0。

## 方案概述

采用**后置转换法**：保留现有 Swagger 2.0 生成逻辑作为内部中间步骤，使用 `kin-openapi` 库自动转换为 OpenAPI 3.0 输出。

## 文档导航

| 文档 | 内容 |
|------|------|
| [01-现状分析](01-现状分析.md) | 当前代码架构、两个命令版本的区别、Swagger 2.0 的局限性 |
| [02-方案设计](02-方案设计.md) | 后置转换法技术细节、kin-openapi 库、2.0→3.0 映射关系 |
| [03-实施步骤](03-实施步骤.md) | 逐步实施指南，包含完整代码示例 |
| [04-验证方案](04-验证方案.md) | 测试方法、验证清单、回归测试 |

## 快速开始

升级后命令行接口不变，直接输出 OpenAPI 3.0：

```shell
goctl api swagger -filename swagger.json -api http/api.api -dir ./http
```

## 改动范围

- 新增 1 个依赖：`github.com/getkin/kin-openapi`
- 新增 1 个文件：`convert.go`（~30 行核心转换函数）
- 修改 2 个文件：`cmd.go`、`command.go`（各加 ~15 行转换调用）
