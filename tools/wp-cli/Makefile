# wp-cli Makefile

BINARY_NAME := wp-cli
GO := env -u GOROOT go
GOLANGCI_LINT := env -u GOROOT golangci-lint
LINT_CONFIG := ../../.golangci.yml
VERSION ?= dev

.PHONY: all build test lint clean coverage fmt vet

## all: ビルド、テスト、リントを全て実行
all: lint test build

## build: バイナリをビルド（-ldflags でバージョン注入）
build:
	$(GO) build -ldflags "-X github.com/shiimanblog/wp-cli/cmd.version=$(VERSION)" -o $(BINARY_NAME) .

## test: テストを実行
test:
	$(GO) test ./...

## coverage: カバレッジ付きテストを実行
coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

## lint: golangci-lint を実行
lint:
	$(GOLANGCI_LINT) run --config $(LINT_CONFIG) ./...

## fmt: コードフォーマットを実行
fmt:
	$(GO) fmt ./...

## vet: 静的解析を実行
vet:
	$(GO) vet ./...

## clean: ビルド成果物を削除
clean:
	rm -f $(BINARY_NAME) coverage.out
