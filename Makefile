# Blog Project Makefile
# tools/wp-cli の Makefile に委譲

WP_CLI_DIR := tools/wp-cli

.PHONY: all build test lint clean install-hooks

## all: ビルド、テスト、リントを全て実行
all: lint test build

## build: wp-cli をビルド
build:
	$(MAKE) -C $(WP_CLI_DIR) build

## test: テストを実行
test:
	$(MAKE) -C $(WP_CLI_DIR) test

## lint: golangci-lint を実行
lint:
	$(MAKE) -C $(WP_CLI_DIR) lint

## clean: ビルド成果物を削除
clean:
	$(MAKE) -C $(WP_CLI_DIR) clean

## install-hooks: pre-commit フックをインストール
install-hooks:
	cp .githooks/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "pre-commit フックをインストールしました"
