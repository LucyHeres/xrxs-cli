GO ?= go
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

LDFLAGS = -s -w \
	-X 'github.com/LucyHeres/xrxs-cli/internal/app.version=$(VERSION)' \
	-X 'github.com/LucyHeres/xrxs-cli/internal/app.buildTime=$(BUILD_TIME)' \
	-X 'github.com/LucyHeres/xrxs-cli/internal/app.gitCommit=$(GIT_COMMIT)'

.PHONY: all help build test lint fmt install release release-dry-run release-local

all: fmt build test

help:
	@printf "Available targets:\n"
	@printf "  make build            - 编译当前平台二进制\n"
	@printf "  make test             - 运行测试\n"
	@printf "  make lint             - 代码检查\n"
	@printf "  make fmt              - 格式化代码\n"
	@printf "  make install          - 安装到 ~/go/bin\n"
	@printf "  make release-dry-run  - 预览发布产物（本地测试）\n"
	@printf "  make release-local    - 本地构建所有平台二进制\n"
	@printf "  make release          - 正式发布 (需要 GitHub Token)\n"

build:
	@mkdir -p bin
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/xrxs ./cmd

test:
	$(GO) test -v -count=1 ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed"

fmt:
	@find . -name '*.go' -not -path './vendor/*' -print0 2>/dev/null | xargs -0r gofmt -w

install: build
	@mkdir -p $(shell go env GOPATH)/bin
	cp bin/xrxs $(shell go env GOPATH)/bin/xrxs
	@echo "Installed to $(shell go env GOPATH)/bin/xrxs"

# 本地构建所有平台（不需要 goreleaser）
release-local:
	@echo "Building for all platforms..."
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/xrxs_linux_amd64/xrxs ./cmd
	GOOS=linux   GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/xrxs_linux_arm64/xrxs ./cmd
	GOOS=darwin  GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/xrxs_darwin_amd64/xrxs ./cmd
	GOOS=darwin  GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/xrxs_darwin_arm64/xrxs ./cmd
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/xrxs_windows_amd64/xrxs.exe ./cmd
	GOOS=windows GOARCH=386   $(GO) build -ldflags "$(LDFLAGS)" -o dist/xrxs_windows_386/xrxs.exe ./cmd
	@echo "Packaging..."
	cd dist && for d in */; do tar -czf "xrxs_$(VERSION)_$${d%/}.tar.gz" -C "$$d" . 2>/dev/null; done
	cd dist && for d in */; do zip -qr "xrxs_$(VERSION)_$${d%/}.zip" "$$d" 2>/dev/null; done
	cd dist && sha256sum *.tar.gz *.zip > checksums.txt 2>/dev/null || shasum -a 256 *.tar.gz *.zip > checksums.txt
	@echo "Done. Artifacts in dist/"

# 使用 goreleaser 预览（不真正发布）
release-dry-run:
	goreleaser release --snapshot --clean --skip=publish

# 使用 goreleaser 正式发布
release:
	goreleaser release --clean
