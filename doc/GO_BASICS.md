# Go 基础知识要点

## 1. 函数返回结构体指针 vs Java 返回对象

### 核心差异：值类型 vs 引用类型

**Java：对象天生是引用**
```java
// Java
public Config load(String path) {
    Config cfg = new Config();  // cfg 本身就是引用
    return cfg;  // 返回的是引用，不是对象本身
}
```

**Go：结构体是值类型**
```go
// Go
func Load(path string) (*Config, error) {
    cfg := Config{}  // cfg 是值，存储实际数据
    return &cfg      // 返回指针（地址）
}
```

### Go 返回指针的原因

| 原因 | 说明 |
|------|------|
| 避免复制开销 | 返回指针只复制 8 字节，返回值会复制整个结构体 |
| 允许修改 | 返回指针，调用者可以直接修改原对象 |
| 支持 nil 语义 | 可以返回 `nil` 表示"无结果"或错误 |
| 方法调用一致性 | 与接收者为指针的方法配合使用 |

### 对比总结

| 特性 | Go | Java |
|------|-----|------|
| 结构体/类 | 值类型 | 引用类型 |
| 赋值/传参 | 复制整个数据 | 复制引用（8字节） |
| 返回对象 | `*Config`（指针） | `Config`（引用） |
| 空值表示 | `nil` | `null` |
| 修改影响 | 指针才影响原对象 | 始终影响原对象 |

---

## 2. 配置加载模式：默认值 + 环境变量覆盖

### 代码示例

```go
cfgPath := "config/config.yaml"
if v := os.Getenv("CONFIG_PATH"); v != "" {
    cfgPath = v
}
cfg, err := config.Load(cfgPath)
if err != nil {
    fmt.Fprintf(os.Stderr, "load config: %v\n", err)
    os.Exit(1)
}
```

### 关键语法：if 初始化语句

```go
if v := os.Getenv("CONFIG_PATH"); v != "" {
    cfgPath = v
}
```

这是 Go 特有的 `if` 语句写法：

```
if 初始化语句; 条件判断 { ... }
```

| 部分 | 代码 | 作用 |
|------|------|------|
| 初始化 | `v := os.Getenv("CONFIG_PATH")` | 读取环境变量 |
| 条件 | `v != ""` | 判断是否非空 |
| 执行体 | `cfgPath = v` | 用环境变量覆盖默认值 |

### 执行流程

```
默认路径: config/config.yaml
        ↓
检查环境变量 CONFIG_PATH
存在且非空? → 覆盖 cfgPath
        ↓
config.Load(cfgPath) 加载配置
        ↓
   成功: cfg      失败: err != nil
        ↓              ↓
   继续执行     打印错误 + os.Exit(1)
```

### 相关 API 对照

| Go | Java |
|-----|------|
| `os.Getenv("KEY")` | `System.getenv("KEY")` |
| `os.Exit(1)` | `System.exit(1)` |
| `fmt.Fprintf(os.Stderr, ...)` | `System.err.println(...)` |
| `if err != nil` | `try-catch` |

### 设计模式：12-Factor App

这是 **12-Factor App** 推荐的配置管理方式：

| 环境 | 配置方式 |
|------|----------|
| 本地开发 | 使用默认 `config/config.yaml` |
| 测试环境 | `CONFIG_PATH=/app/config/test.yaml` |
| 生产环境 | `CONFIG_PATH=/app/config/prod.yaml` |

**运行示例：**
```bash
# 本地开发（用默认路径）
go run cmd/server/main.go

# 生产环境（指定配置）
CONFIG_PATH=/etc/app/prod.yaml go run cmd/server/main.go

# Docker 环境
docker run -e CONFIG_PATH=/app/config.yaml myapp
```

---

## 3. Context 上下文

### 什么是 Context？

Context 是 Go 中用于**传递请求范围数据、取消信号、超时控制**的标准机制。

```
┌─────────────────────────────────────────────────────┐
│                    Context 树                        │
│                                                      │
│              context.Background()  ← 根节点          │
│                    │                                 │
│          ┌────────┼────────┐                        │
│          ↓        ↓        ↓                         │
│      WithTimeout  WithCancel  WithValue             │
│          │         │          │                      │
│          ↓         ↓          ↓                      │
│       子请求    子请求     携带数据                    │
└─────────────────────────────────────────────────────┘
```

### Context 类型

| 类型 | 创建方式 | 特点 |
|------|----------|------|
| Background | `context.Background()` | 空的根 Context，永不取消 |
| TODO | `context.TODO()` | 占位符，暂不确定用哪个 |
| WithTimeout | `context.WithTimeout(parent, d)` | 超时自动取消 |
| WithCancel | `context.WithCancel(parent)` | 手动取消 |
| WithValue | `context.WithValue(parent, k, v)` | 携带数据 |

### context.Background() 的特点

| 特性 | 值 |
|------|-----|
| 可取消 | ❌ 永不取消 |
| 超时 | ❌ 无超时 |
| 截止时间 | ❌ 无截止时间 |
| 携带数据 | ❌ 无数据 |
| 用途 | 作为根 Context |

### 使用场景

**1. 程序入口（main 函数）**

```go
func main() {
    ctx := context.Background()  // 根 Context，作为起点
    
    // 数据库连接
    database, err := db.New(ctx, cfg.Database.URL)
    
    // 服务器优雅关闭
    srv.Shutdown(ctx)
}
```

**2. HTTP 请求处理**

```go
// ctx 由 HTTP 框架自动创建，随请求结束而取消
func (h *Handler) GetUser(ctx context.Context, input *struct{}) {
    userID := middleware.GetUserID(ctx)  // 从 ctx 获取数据
    user, err := h.svc.GetByID(ctx, userID)
}
```

**3. 超时控制**

```go
func (s *Service) FetchData(ctx context.Context) error {
    // 派生一个 3 秒超时的 ctx
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()  // 重要！函数返回时释放资源
    
    result, err := db.Query(ctx, "SELECT ...")
    // 3 秒后 ctx 自动取消
    return err
}
```

**4. 监听取消信号**

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():  // Context 被取消
            fmt.Println("收到取消信号:", ctx.Err())
            return
        default:
            doWork()
        }
    }
}
```

### Context 生命周期

**核心原则：Context 的生命周期由创建者控制，随请求/任务结束而结束。**

```
创建 ──→ 使用 ──→ 取消/超时 ──→ 结束 ──→ GC 回收
 │        │         │
 │        │         └── cancel() 或 超时自动触发
 │        └── 传递给下游函数
 └── Background() / WithTimeout / WithCancel
```

**派生 Context 的生命周期：子 Context 随父 Context 取消而取消**

```
Background() ──────────────────────────────────────→ 永不取消
    │
    └── WithTimeout(3s) ─────→ 3秒后取消
            │
            └── WithValue ────→ 随父取消
                    │
                    └── WithCancel ─→ 随父取消 或 手动取消
```

### 最佳实践

```go
// 1. 始终 defer cancel（防止内存泄漏）
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// 2. Context 作为函数第一个参数
func DoSomething(ctx context.Context, arg string) {  // ✅
}

// 3. 不要把 Context 存储在结构体中
type Service struct {
    ctx context.Context  // ❌ 不要这样做
}

// 4. 不要传递 nil Context
func bad(ctx context.Context) {
    if ctx == nil {
        ctx = context.Background()  // 替代方案
    }
}
```

### 与 Java 的类比

| Go | Java 类比 |
|-----|----------|
| `context.Background()` | 空的初始状态 / main 方法入口 |
| `context.WithValue()` | `ThreadLocal.set()` 或 `request.setAttribute()` |
| `ctx.Value()` | `ThreadLocal.get()` 或 `request.getAttribute()` |
| `context.WithTimeout()` | `Future.get(timeout)` 或 `ExecutorService` 超时 |

**关键区别**：Go 的 Context 是**显式传递**的，Java 通常用 **ThreadLocal 或请求对象**隐式传递。

```java
// Java - ThreadLocal 隐式传递
public class RequestContext {
    private static final ThreadLocal<Long> USER_ID = new ThreadLocal<>();
    
    public static void setUserId(Long id) { USER_ID.set(id); }
    public static Long getUserId() { return USER_ID.get(); }
}
```

```go
// Go - Context 显式传递
ctx := context.WithValue(context.Background(), "user_id", 123)
userId := ctx.Value("user_id")
```

### main 中 ctx 的典型用途

```go
func main() {
    ctx := context.Background()
    
    // 用途1：数据库连接（控制连接建立超时）
    database, err := db.New(ctx, cfg.Database.URL)
    
    // 用途2：服务器优雅关闭（控制等待请求完成的时间）
    <-quit  // 等待退出信号
    srv.Shutdown(ctx)
}
```

**改进建议：为不同操作设置超时**

```go
// 数据库连接：10秒超时
dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
database, err := db.New(dbCtx, cfg.Database.URL)
dbCancel()

// 服务器关闭：30秒超时
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
defer shutdownCancel()
srv.Shutdown(shutdownCtx)
```

---

## 4. 结构体字面量与返回指针

### 代码示例

```go
type DB struct {
    Pool *pgxpool.Pool
}

func New(ctx context.Context, connStr string) (*DB, error) {
    pool, err := pgxpool.New(ctx, connStr)
    if err != nil {
        return nil, err
    }
    return &DB{Pool: pool}, nil
}
```

### 语法拆解

```go
return &DB{Pool: pool}, nil
```

| 部分 | 含义 |
|------|------|
| `DB{Pool: pool}` | 结构体字面量，创建 DB 实例，Pool 字段赋值为 pool 变量 |
| `&DB{...}` | 取地址，得到 `*DB` 类型的指针 |
| `..., nil` | 多返回值：指针 + error |

### 完整流程

```
┌─────────────────────────────────────────────────────┐
│  return &DB{Pool: pool}, nil                        │
│                                                      │
│  ┌─────────────┐   ┌─────────────┐   ┌───────────┐ │
│  │ DB{Pool:    │ → │ &DB{...}    │ → │ return    │ │
│  │   pool}     │   │ 取地址      │   │ 多返回值  │ │
│  │ 创建结构体  │   │ 得到指针    │   │           │ │
│  └─────────────┘   └─────────────┘   └───────────┘ │
│                                                      │
│  结果: 返回 *DB 类型的指针 + nil error              │
└─────────────────────────────────────────────────────┘
```

### 结构体字面量的几种写法

```go
// 写法1：字段名:值（推荐，可读性好）
return &DB{Pool: pool}, nil

// 写法2：分步创建
db := DB{Pool: pool}
return &db, nil

// 写法3：new + 赋值
db := new(DB)
db.Pool = pool
return db, nil

// 写法4：按顺序（不推荐，可读性差）
return &DB{pool}, nil
```

### 与 Java 对比

**Go：**
```go
type DB struct {
    Pool *pgxpool.Pool
}

return &DB{Pool: pool}, nil  // 直接创建并返回指针
```

**Java：**
```java
public class DB {
    private Pool pool;
    
    public DB(Pool pool) {
        this.pool = pool;
    }
}

return new DB(pool);  // new 返回的就是引用
```

| Go | Java |
|-----|------|
| `&DB{Pool: pool}` | `new DB(pool)` |
| 显式取地址 `&` | `new` 天生返回引用 |
| 字段名:值 初始化 | 构造函数初始化 |

### 为什么返回指针？

| 原因 | 说明 |
|------|------|
| 避免复制 | 结构体较大时，返回指针只复制 8 字节 |
| 一致性 | 方法接收者 `func (db *DB)` 是指针类型 |
| 可修改 | 调用者可以修改结构体的状态 |
| nil 语义 | 可以返回 nil 表示失败或无结果 |

### 一句话总结

`return &DB{Pool: pool}, nil` = 创建结构体 + 取地址 + 多返回值，是 Go 中返回结构体指针的标准写法。

---

## 5. 方法接收者

### 语法结构

```go
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
```

```
func (接收者) 方法名(参数) 返回值
     ↓
func (s *UserService) Register(ctx context.Context, ...) (*dto.AuthResponse, error)
```

| 部分 | 含义 |
|------|------|
| `s` | 接收者变量名（类似 `this` / `self`） |
| `*UserService` | 接收者类型（指针类型） |
| `Register` | 方法名 |

### 作用：将方法绑定到类型

```go
type UserService struct {
    db    *pgxpool.Pool
    redis *redis.Client
    jwt   *jwt.Manager
}

// Register 方法绑定到 *UserService
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    // s 可以访问 UserService 的字段
    s.db.Query(ctx, "INSERT INTO users...")
    s.redis.Get(ctx, "key")
    s.jwt.Generate(userID)
}
```

### 调用方式

```go
// 创建实例
svc := &UserService{db: pool, redis: rdb, jwt: jwtMgr}

// 通过实例调用方法
resp, err := svc.Register(ctx, req)  // s = svc
```

### 与普通函数的区别

```go
// 普通函数：独立存在
func Register(svc *UserService, ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    svc.db.Query(...)
}
// 调用: Register(svc, ctx, req)

// 方法：绑定到类型
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    s.db.Query(...)
}
// 调用: svc.Register(ctx, req)
```

### 指针接收者 vs 值接收者

```go
// 指针接收者（常用）
func (s *UserService) Register(...) { }  // s 是指针，可以修改结构体

// 值接收者
func (s UserService) Register(...) { }   // s 是副本，修改不影响原对象
```

| 类型 | 特点 | 适用场景 |
|------|------|----------|
| 指针接收者 `*T` | 可以修改结构体、避免复制 | 结构体较大、需要修改状态 |
| 值接收者 `T` | 操作副本、不修改原对象 | 小结构体、只读操作 |

### 与 Java 对比

**Go：**
```go
type UserService struct {
    db *pgxpool.Pool
}

func (s *UserService) Register(ctx context.Context, req *Request) (*Response, error) {
    s.db.Query(...)  // s 类似 this
}
```

**Java：**
```java
public class UserService {
    private Pool db;
    
    public Response register(Request req) {
        this.db.query(...);  // this 是隐式的
    }
}
```

| Go | Java |
|-----|------|
| `(s *UserService)` 显式定义接收者 | `this` 隐式存在 |
| `s.db` 访问字段 | `this.db` 或直接 `db` |
| 指针/值接收者可选 | 只有引用 |

---

## 6. 数据库查询与 Scan 方法

### sqlc 生成的查询代码

```go
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
    row := q.db.QueryRow(ctx, getUserByEmail, email)
    var i User
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Email,
        &i.Phone,
        &i.Password,
        &i.Status,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}
```

### 逐行解析

| 行 | 代码 | 含义 |
|-----|------|------|
| 1 | `func (q *Queries) GetUserByEmail(...)` | 方法接收者，q 是 Queries 实例 |
| 2 | `row := q.db.QueryRow(ctx, getUserByEmail, email)` | 执行 SQL 查询，返回单行 |
| 3 | `var i User` | 声明空结构体，用于接收结果 |
| 4-12 | `row.Scan(&i.ID, ...)` | 将查询结果映射到结构体字段 |
| 13 | `return i, err` | 返回结果和错误 |

### 完整流程

```
┌─────────────────────────────────────────────────────────┐
│                    GetUserByEmail                        │
│                                                          │
│  输入: email = "user@example.com"                        │
│         ↓                                                │
│  ┌─────────────────────────────────────┐                │
│  │ QueryRow(ctx, SQL, email)           │                │
│  │ 执行: SELECT * FROM users           │                │
│  │       WHERE email = 'user@...'      │                │
│  └─────────────────┬───────────────────┘                │
│                    ↓                                     │
│  ┌─────────────────────────────────────┐                │
│  │ row.Scan(&i.ID, &i.Name, ...)       │                │
│  │ 将查询结果映射到 User 结构体        │                │
│  └─────────────────┬───────────────────┘                │
│                    ↓                                     │
│  输出: User{ID:1, Name:"张三", ...}, nil                │
└─────────────────────────────────────────────────────────┘
```

### 为什么 Scan 需要传指针 `&i.ID`

**核心原因：Scan 需要修改变量的值**

```go
// Scan 的函数签名
func (r *Row) Scan(dest ...any) error
```

`dest` 参数需要传入**指针**，这样 Scan 才能修改原始变量。

```go
var i User

// ❌ 错误：传值，Scan 无法修改 i.ID
row.Scan(i.ID)  // 传入的是 ID 的副本，Scan 修改的是副本

// ✅ 正确：传指针，Scan 可以修改 i.ID
row.Scan(&i.ID) // 传入的是 ID 的地址，Scan 通过地址修改值
```

### 图解：传值 vs 传指针

```
传值 i.ID:
┌─────────────────────────────────────────────┐
│  i.ID = 0          Scan 内部                │
│  ┌─────┐           ┌─────┐                  │
│  │  0  │ ──复制──→ │  0  │ ← Scan 修改副本  │
│  └─────┘           └─────┘                  │
│    ↑                                         │
│    └── 原值不变，仍然是 0                    │
└─────────────────────────────────────────────┘

传指针 &i.ID:
┌─────────────────────────────────────────────┐
│  i.ID = 0          Scan 内部                │
│  ┌─────┐           ┌─────────┐              │
│  │  0  │ ←──────── │ 地址 0x1 │ ← 通过地址  │
│  └─────┘           └─────────┘   修改原值   │
│    ↑                                         │
│    └── 原值被修改，变成查询结果              │
└─────────────────────────────────────────────┘
```

### 类比理解

```go
func setValue(val int) {
    val = 100  // 修改副本，外部变量不变
}

func setValueByPointer(val *int) {
    *val = 100  // 通过指针修改，外部变量改变
}

func main() {
    x := 0
    setValue(x)
    fmt.Println(x)  // 输出: 0（没变）
    
    setValueByPointer(&x)
    fmt.Println(x)  // 输出: 100（变了）
}
```

### 与 Java 对比

**Go：显式传指针**
```go
var id int64
row.Scan(&id)  // 必须传指针
```

**Java：ResultSet 返回值**
```java
long id = rs.getLong("id");  // 返回值，直接赋值
```

| Go | Java |
|-----|------|
| `Scan(&i.ID)` 传指针 | `rs.getLong()` 返回值 |
| Scan 写入传入的变量 | ResultSet 返回值，自己赋值 |

### 一句话总结

- **方法接收者** `(s *UserService)`：将方法绑定到类型，`s` 类似 `this`
- **Scan 传指针**：`&i.ID` 传递地址，让 Scan 能够修改原变量的值
