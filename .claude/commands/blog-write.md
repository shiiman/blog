# Blog Write Command

技術ブログ記事を執筆します。

## 使い方

```bash
/blog-write
/blog-write <テーマ>
```

## 実行例

```bash
# 対話的に記事を作成
/blog-write

# テーマを指定して作成
/blog-write ConoHa VPSでDockerをセットアップする方法
```

## Claudeへの指示

1. テーマがない場合はユーザーにテーマを聞く
2. WebSearchでSEOキーワードをリサーチ
3. 見出し構造を提案して承認を得る
4. `.claude/agents/blog-writer.md` の形式に従って記事を執筆
5. `drafts/YYYY-MM-DD_slug/article.md` に保存
6. 次のステップ（/blog-publish）を案内

$ARGUMENTS
