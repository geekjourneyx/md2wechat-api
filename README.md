# MD2WeChat API

一个将 Markdown 内容转换为微信公众号草稿的 OpenAPI 服务。

## 功能特性

- ✅ 将 Markdown 转换为微信公众号兼容的 HTML 格式
- ✅ 自动创建微信公众号草稿
- ✅ 支持多种主题和字体大小
- ✅ RESTful API 设计
- ✅ Swagger 文档自动生成
- ✅ 结构化日志记录
- ✅ 参数验证
- ✅ 错误处理
- ✅ Docker 容器化部署

## 技术栈

- **语言**: Go 1.24+
- **Web 框架**: Gin
- **微信 SDK**: github.com/silenceper/wechat/v2
- **日志**: Zap
- **配置管理**: Viper
- **文档**: Swagger
- **容器化**: Docker

## 项目结构

```
md2wechat-api/
├── cmd/
│   └── server/           # 应用入口
├── internal/
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP 处理器
│   ├── model/           # 数据模型
│   └── service/         # 业务逻辑
├── pkg/
│   └── logger/          # 日志工具
├── configs/             # 配置文件
├── docs/                # Swagger 文档
├── Dockerfile           # Docker 构建文件
├── Makefile            # 构建脚本
└── README.md           # 项目文档
```

## 快速开始

### 环境要求

- Go 1.24+
- 微信公众号（已认证）
- MD2WeChat API Key（可选）

### 安装依赖

```bash
make deps
```

### 配置

1. 复制配置文件模板：
```bash
cp configs/config.yaml configs/config.local.yaml
```

2. 编辑配置文件 `configs/config.local.yaml`：
```yaml
server:
  port: 8080
  mode: debug

log:
  level: info

md2wechat:
  base_url: https://www.md2wechat.cn
  api_key: "your_api_key_here"  # 如果需要
```

### 运行

#### 开发模式
```bash
make dev
```

#### 生产模式
```bash
make build
make run
```

#### Docker 运行
```bash
make docker
docker run -p 8080:8080 md2wechat-api:latest
```

## API 文档

启动服务后，访问 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) 查看完整的 API 文档。

### 核心 API

#### POST /api/v1/convert-and-draft

将 Markdown 内容转换为 HTML 并创建微信草稿。

**请求示例:**

```bash
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -H "Wechat-Appid: wx1234567890abcdef" \
  -H "Wechat-App-Secret: your_app_secret_here" \
  -H "Md2wechat-API-Key: wme_your_api_key_here" \
  -d '{
    "markdown": "# 标题\n\n这是一个**加粗**文本的例子。",
    "theme": "default",
    "fontSize": "medium"
  }'
```

**响应示例:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "media_id": "media_id_123456",
    "html": "<section>...</section>",
    "theme": "default",
    "fontSize": "medium",
    "wordCount": 156,
    "estimatedReadTime": 1
  },
  "timestamp": 1640995200
}
```

## 配置说明

### 环境变量

可以通过环境变量覆盖配置文件中的设置：

```bash
export SERVER_PORT=8080
export SERVER_MODE=release
export LOG_LEVEL=info
export MD2WECHAT_BASE_URL=https://www.md2wechat.cn
export MD2WECHAT_API_KEY=your_api_key
```

### Header 参数配置

调用 API 时需要在请求头中设置以下参数：

**必填参数：**
- `Wechat-Appid`: 微信公众号的 AppID（从微信公众号后台获取）
- `Wechat-App-Secret`: 微信公众号的 AppSecret（从微信公众号后台获取）

**可选参数：**
- `Md2wechat-API-Key`: MD2WeChat 服务的 API 密钥（如果服务需要认证）

### MD2WeChat 配置

- `base_url`: MD2WeChat 服务地址（默认：https://www.md2wechat.cn）

## 开发指南

### 代码格式化
```bash
make fmt
```

### 代码检查
```bash
make lint
```

### 生成文档
```bash
make swagger
```

### 运行测试
```bash
make test
```

### 清理构建文件
```bash
make clean
```

## 部署

### Docker 部署

1. 构建镜像：
```bash
make docker
```

2. 运行容器：
```bash
docker run -d \
  --name md2wechat-api \
  -p 8080:8080 \
  -e SERVER_MODE=release \
  -e LOG_LEVEL=info \
  md2wechat-api:latest
```

### Kubernetes 部署

可以使用以下 YAML 文件进行部署：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: md2wechat-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: md2wechat-api
  template:
    metadata:
      labels:
        app: md2wechat-api
    spec:
      containers:
      - name: md2wechat-api
        image: md2wechat-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_MODE
          value: "release"
        - name: LOG_LEVEL
          value: "info"
---
apiVersion: v1
kind: Service
metadata:
  name: md2wechat-api-service
spec:
  selector:
    app: md2wechat-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

## 错误处理

API 使用统一的错误响应格式：

```json
{
  "code": 400,
  "message": "参数验证失败",
  "details": "markdown 字段不能为空",
  "timestamp": 1640995200
}
```

常见错误代码：
- `400`: 请求参数错误
- `500`: 服务器内部错误

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

本项目采用 Apache 2.0 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件至：support@example.com

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持 Markdown 转换和微信草稿创建
- 完整的 API 文档和部署指南 