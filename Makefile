.PHONY: build install clean test

# 构建
build:
	go build -o bin/gopherflush ./cmd/gopherflush

# 安装
install:
	go install ./cmd/gopherflush

# 清理
clean:
	rm -rf bin/

# 测试
test:
	go test -v ./...

# 下载依赖
deps:
	go mod tidy
	go mod download

# 运行
run:
	go run ./cmd/gopherflush
