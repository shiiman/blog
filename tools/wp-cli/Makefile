# wp-cli Makefile

BINARY_NAME := wp-cli
GO := go
GOLANGCI_LINT := golangci-lint
LINT_CONFIG := ../../.golangci.yml

.PHONY: all build test lint clean

## all: ビルド、テスト、リントを全て実行
all: lint test build

## build: バイナリをビルド
build:
	$(GO) build -o $(BINARY_NAME) .

## test: テストを実行
test:
	$(GO) test ./...

## lint: golangci-lint を実行
lint:
	$(GOLANGCI_LINT) run --config $(LINT_CONFIG) ./...

## clean: ビルド成果物を削除
clean:
	rm -f $(BINARY_NAME)
