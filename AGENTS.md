# gtoken — Agent 指南

## 项目概述
基于 GoFrame v2.x 的 Token 认证库，提供服务端 Token 管理，支持 gcache、gredis 和文件存储。

## 模块结构
- **根模块**: `github.com/goflyfox/gtoken/v2` — 核心库在 `gtoken/`
- **示例模块**: `gtoken-demo` 在 `example/` — 独立演示应用，带 `replace` 指令指向根模块
- **扩展模块**: `contrib/jwt/` — JWT Token 扩展（独立模块）

## 命令
```bash
# 单元测试 (gtoken 包)
cd gtoken && go test ./...

# 集成测试 (示例应用)
cd example && go test ./...

# 运行示例服务器
cd example && go run .
```

## 关键约定
- **中间件顺序**: CORS → `Auth` → 业务处理器 → 响应处理器
- **Token 获取**: `Authorization: Bearer <token>` 请求头 或 `token` 表单字段
- **认证错误码**: `gcode.CodeBusinessValidationFailed` (默认)
- **上下文键**: `KeyUserKey` 存储已认证用户标识

## 存储模式
| CacheMode | 存储方式 | 使用场景 |
|-----------|---------|---------|
| 1 | gcache | 单机测试 |
| 2 | gredis | 生产集群 |
| 3 | file | 个人项目 |

## 测试说明
- 集成测试 (`example/api_test.go`) 通过 `TestMain` 启动/停止服务器
- 测试默认使用端口 `8083` (`TestURL` 常量)
- 单元测试使用内存 gcache；集成测试使用配置驱动的存储

## 配置
示例应用从 `config.yaml` 加载 `gToken` 配置。核心库接受 `Options{}` 编程式配置。

## 版本兼容性
- GoFrame v2.x → gtoken v2.x (当前)
- GoFrame v1.x → gtoken v1.4.x (旧版)
