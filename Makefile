# NetKit Makefile

# 变量
BINARY    = netkit
VERSION   = $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")
COMMIT    = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE      = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS   = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)
GO        = go
GOFLAGS   = -trimpath

# 平台
GOOS      = $(shell go env GOOS)
GOARCH    = $(shell go env GOARCH)

.PHONY: all build clean test lint fmt vet run install release cross-build

all: lint test build

## build: 编译当前平台
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) .

## run: 快速运行
run:
	$(GO) run . $(ARGS)

## install: 安装到 GOPATH/bin
install:
	$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS)" .

## test: 运行测试
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

## lint: 运行 golangci-lint
lint:
	golangci-lint run ./...

## fmt: 格式化代码
fmt:
	gofmt -w .
	$(GO) mod tidy

## vet: 静态分析
vet:
	$(GO) vet ./...

## clean: 清理编译产物
clean:
	rm -f $(BINARY)
	rm -f coverage.out
	rm -rf dist/

## cross-build: 交叉编译所有平台
cross-build:
	GOOS=linux   GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-arm64.exe .
	GOOS=windows GOARCH=386   $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-386.exe .

## release: 使用 GoReleaser 发布
release:
	goreleaser release --rm-dist

## snapshot: 使用 GoReleaser 生成快照
snapshot:
	goreleaser release --snapshot --rm-dist

## help: 显示帮助
help:
	@echo "NetKit Makefile 命令:"
	@grep -E '^## ' Makefile | sed 's/## //'
