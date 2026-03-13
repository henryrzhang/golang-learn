# golang-learn

Golang 单体应用示例，技术栈：huma+chi、sqlc+pgx、Redis、JWT、PostgreSQL。

## 功能

- 用户注册、登录、获取用户信息
- 剧集 CRUD（创建、查询、列表、更新、软删除）
- JWT 鉴权
- 分级日志（debug/info/warn/error）
- 本地配置文件 + 环境变量覆盖

## 快速开始

```bash
# 1. 创建数据库
createdb golang_learn

# 2. 执行建表
make migrate

# 3. 修改 config/config.yaml 中的数据库、Redis 连接

# 4. 运行
go run ./cmd/server
```

## API 路径

| 模块 | 路径前缀 |
|------|----------|
| 健康检查 | `/health`、`/ready` |
| 认证 | `/api/auth` |
| 用户 | `/api/user` |
| 剧集 | `/api/drama` |

启动后访问 http://localhost:8080/docs 查看 OpenAPI 文档。

## 项目结构

见 [doc/PROJECT_STRUCTURE.md](doc/PROJECT_STRUCTURE.md)。
