.PHONY: help build run dev test clean docker swagger deps fmt lint

# 默认目标
help:
	@echo "Available commands:"
	@echo "  deps     - 下载依赖"
	@echo "  fmt      - 格式化代码"
	@echo "  lint     - 代码检查"
	@echo "  swagger  - 生成 Swagger 文档"
	@echo "  build    - 构建应用"
	@echo "  run      - 运行应用"
	@echo "  dev      - 开发模式运行"
	@echo "  test     - 运行测试"
	@echo "  clean    - 清理构建文件"
	@echo "  docker   - 构建 Docker 镜像"

# 下载依赖
deps:
	go mod tidy
	go mod download

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 生成 Swagger 文档
swagger:
	swag init -g cmd/server/main.go -o docs

# 构建应用
build: swagger
	go build -o bin/md2wechat-api cmd/server/main.go

# 运行应用
run: build
	./bin/md2wechat-api

# 开发模式运行
dev: swagger
	go run cmd/server/main.go

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf docs/

# 构建 Docker 镜像
docker:
	docker build -t md2wechat-api:latest . 