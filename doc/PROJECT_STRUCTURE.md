# 项目结构说明

## 技术栈

| 层级 | 技术选型 | 说明 |
|------|----------|------|
| HTTP 框架 | huma + chi | REST API，OpenAPI 3.1 文档 |
| 数据库 | PostgreSQL + sqlc + pgx | 类型安全 SQL，高性能驱动 |
| 缓存 | Redis | 会话、热点数据缓存 |
| 鉴权 | JWT | 无状态 Token 鉴权 |
| 配置 | 本地 YAML 文件 | 环境变量可覆盖 |

## 目录结构

```
golang-learn/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── config/
│   ├── config.go               # 配置加载与管理
│   └── config.yaml             # 默认配置文件
├── doc/
│   ├── PROJECT_STRUCTURE.md    # 本文件
│   └── schema.sql              # PostgreSQL 建表语句
├── internal/
│   ├── infra/                  # 基础设施层（纯基建：连接、工具）
│   │   ├── db/                 # 数据库连接池
│   │   ├── redis/              # Redis 客户端
│   │   ├── jwt/                # JWT 工具
│   │   └── logger/             # 日志
│   ├── repository/             # 数据访问层
│   │   └── sqlc/               # sqlc 生成（模型、查询）
│   ├── router/                 # 路由注册（健康检查、/api 分组、中间件）
│   │   └── router.go
│   ├── handler/                # HTTP 处理器
│   │   ├── user.go             # 注册、登录、用户信息
│   │   └── drama.go            # 剧集 CRUD
│   ├── middleware/             # 中间件
│   │   ├── auth.go             # JWT 鉴权
│   │   └── logger.go           # 请求日志
│   ├── dto/                    # API 请求/响应 DTO
│   │   ├── user.go
│   │   └── drama.go
│   └── service/                # 业务逻辑层
│       ├── user.go             # 注册、登录、用户信息（含 Redis 缓存）
│       └── drama.go
├── sqlc/                       # sqlc 配置与 SQL
│   ├── sqlc.yaml
│   ├── schema/
│   │   └── schema.sql
│   └── queries/
│       ├── user.sql
│       └── drama.sql
├── go.mod
├── go.sum
└── Makefile
```

## 分层说明

- **router**: 路由分组（/health、/ready 无鉴权；/api/* 统一中间件）、中间件注册
- **handler**: 解析请求、校验参数、调用 service、返回响应
- **service**: 业务逻辑，编排 db/redis/jwt 等
- **dto**: API 请求/响应 DTO，与 OpenAPI 绑定
- **middleware**: 鉴权、日志、限流等横切逻辑
- **infra**: 基础设施层（db 连接池、redis、jwt、logger），纯基建无业务
- **repository**: 数据访问层（sqlc 模型与查询），service 通过 infra/db.Pool 调用

## 运行说明

```bash
# 安装依赖
go mod tidy

# 生成 sqlc 代码
make sqlc

# 运行（需配置 PostgreSQL、Redis）
go run ./cmd/server
```
