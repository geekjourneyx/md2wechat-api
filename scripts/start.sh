#!/bin/bash

# 启动脚本

set -e

echo "🚀 启动 MD2WeChat API 服务..."

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ Go 环境未安装，请先安装 Go 1.24+"
    exit 1
fi

# 检查项目依赖
echo "📦 检查项目依赖..."
go mod tidy

# 生成 Swagger 文档
echo "📚 生成 API 文档..."
if command -v swag &> /dev/null; then
    swag init -g cmd/server/main.go -o docs
else
    echo "⚠️  swag 工具未安装，跳过文档生成"
    echo "💡 可以运行以下命令安装: go install github.com/swaggo/swag/cmd/swag@latest"
fi

# 设置环境变量（如果配置文件存在）
if [ -f "configs/config.local.yaml" ]; then
    echo "📄 使用本地配置文件: configs/config.local.yaml"
elif [ -f ".env" ]; then
    echo "📄 加载环境变量文件: .env"
    set -a
    source .env
    set +a
else
    echo "⚠️  未找到配置文件，使用默认配置"
    echo "💡 您可以复制 configs/config.yaml 到 configs/config.local.yaml 进行自定义配置"
fi

# 构建并运行
echo "🔨 构建应用..."
go build -o bin/md2wechat-api cmd/server/main.go

echo "🎯 启动服务..."
./bin/md2wechat-api

echo "✅ 服务已启动！"
echo "📖 API 文档: http://localhost:8080/swagger/index.html"
echo "💚 健康检查: http://localhost:8080/health" 