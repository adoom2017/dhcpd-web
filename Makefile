# 定义构建变量
CC=go
BUILD_FLAGS=-v -ldflags "-s -w" -trimpath

# 默认目标
all: build

# 构建目标
build:
	$(CC) build -o ./build/dhcpd-web  $(BUILD_FLAGS) .
	cp -rf ./web ./build/

# 测试目标
test:
	$(CC) test $(BUILD_FLAGS) ./...

# 格式化代码
fmt:
	$(CC) fmt ./...

# 清理所有构建文件
clean:
	rm -rf ./build

# 伪目标，确保在最后执行
.PHONY: all build test fmt clean
