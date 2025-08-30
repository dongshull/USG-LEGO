# USG-LEGO

USG-LEGO 是一个基于 Go 和 Gin 框架开发的单文件托管服务，专门为 iOS 代理工具（如 Surge、Shadowrocket 等）提供规则和配置文件托管。

## 功能特性

- **文件托管服务**: 为 iOS 代理工具提供规则和配置文件托管
- **权限控制**: 支持公共和私有路径，私有文件需要 API 密钥访问
- **API 密钥认证**: 基于 API 密钥的访问控制机制
- **安全防护**: 防止路径遍历攻击，确保文件访问安全
- **多种响应格式**: 支持文件流和 JSON 元数据响应
- **环境变量配置**: 通过环境变量进行灵活配置
- **日志记录**: 使用 zerolog 记录结构化日志
- **Docker 支持**: 提供多架构 Docker 镜像

## 技术栈

- **后端**: Go + Gin 框架
- **数据库**: SQLite（用于存储 API 密钥）
- **日志库**: zerolog
- **前端**: 嵌入式 SPA（单页应用）

## 目录结构

```
.
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置管理
│   ├── auth/
│   │   ├── auth.go          # 认证逻辑
│   │   └── apikey/
│   │       └── apikey.go    # API 密钥管理
│   ├── files/
│   │   ├── handler.go       # 文件处理接口
│   │   └── service.go       # 文件服务逻辑
│   ├── routes/
│   │   └── routes.go        # 路由注册
│   └── logger/
│       └── logger.go        # 日志处理
├── pkg/
│   └── utils/
│       └── utils.go         # 工具函数
├── web/
│   └── dist/                # 前端静态文件
├── configs/
│   └── .usg-lego.yml        # 路径规则配置
├── Dockerfile               # Docker 构建文件
└── docs/
    └── openapi.yaml         # API 文档
```

## 环境变量配置

| 变量名       | 描述             | 默认值    | 必需 |
|--------------|------------------|-----------|------|
| ROOT_DIR     | 文件根目录       | 无        | 是   |
| LISTEN       | 监听地址         | :8080     | 否   |
| LOG_LEVEL    | 日志级别         | info      | 否   |
| JWT_SECRET   | JWT 密钥         | 无        | 否   |

## 路径规则配置

在 `ROOT_DIR/.usg-lego.yml` 文件中配置公共和私有路径规则：

```yaml
public:
  - "/public/**"
  - "/assets/**"

private:
  - "/private/**"
  - "/conf/**"
```

## API 接口

### 健康检查

```
GET /health
```

检查服务运行状态，返回 "ok"。

### 文件访问

```
GET /api/files?path=<文件路径>&api=<API密钥>
```

获取指定路径的文件内容。

**参数:**
- `path` (必需): 文件路径
- `api` (可选): API 密钥，访问私有文件时必需

**响应:**
- 默认返回文件流
- 如果请求头包含 `Accept: application/json`，返回文件元数据

### 原始文件访问

```
GET /raw/<文件路径>?api=<API密钥>
```

直接访问文件内容。

### API 密钥管理

```
GET /api/keys
```

列出所有 API 密钥（仅管理员）。

```
POST /api/keys
```

生成新的 API 密钥。

```
DELETE /api/keys/:id
```

删除指定的 API 密钥。

## 权限控制

### 公共路径
- 以 `/public` 开头的路径为公共路径
- 公共路径下的文件可无需认证直接访问

### 私有路径
- 除公共路径外的所有路径均为私有路径
- 访问私有路径需要提供有效的 API 密钥

## 安全机制

1. **路径遍历防护**: 所有路径参数都经过 `filepath.Clean` 处理，并检查不能跳出 ROOT_DIR
2. **API 密钥验证**: 私有文件访问需要提供有效的 32 字节随机字符串形式的 API 密钥
3. **访问控制**: 写操作和私有路径读操作都必须携带有效的 API 密钥

## 日志格式

使用 zerolog 输出 JSON 格式日志，包含以下字段：
- `time`: 时间戳
- `level`: 日志级别
- `msg`: 日志消息
- `path`: 请求路径
- `ip`: 客户端 IP

## 快速开始

### 使用 Go 直接运行

```bash
# 克隆项目
git clone <项目地址>
cd USG-LEGO

# 设置环境变量
export ROOT_DIR=./configs
export LISTEN=:8080

# 运行服务
go run cmd/server/main.go
```

### 使用 Docker 运行

```bash
# 构建镜像
docker build -t usg-lego .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v /path/to/your/files:/app/configs \
  -e ROOT_DIR=/app/configs \
  usg-lego
```

## 构建

```bash
# 构建二进制文件
go build -o usg-lego cmd/server/main.go
```

## 测试

```bash
# 运行测试
go test ./...
```

## Docker 支持

Dockerfile 支持多阶段构建，生成轻量级的 Alpine 镜像，支持 linux/amd64 和 linux/arm64 架构。

## 开发规范

1. 函数长度不超过 60 行
2. 错误尽早 return
3. 统一使用 ctx 传递 request-scoped 值
4. 所有外部依赖需加入 go.sum 校验
5. 使用 golangci-lint 进行代码检查

## 许可证

[MIT License](LICENSE)