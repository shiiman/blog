# カテゴリ・タグ slug リファレンス

記事の Front Matter で使用するカテゴリ/タグは **文字列 slug** で指定します（旧 WordPress の数値IDは使いません）。

## 一覧の参照元

slug の一覧は以下の JSON で管理されています。各エントリの `slug`（英語 slug）を Front Matter に指定してください。

- カテゴリ: [`data/categories.json`](../data/categories.json)
- タグ: [`data/tags.json`](../data/tags.json)

```bash
# 例: カテゴリの name と slug を一覧する
node -e "const c=require('./data/categories.json'); for(const k in c) console.log(c[k].slug, '-', c[k].name)"

# タグも同様
node -e "const t=require('./data/tags.json'); for(const k in t) console.log(t[k].slug, '-', t[k].name)"
```

## Front Matter での使い方

```yaml
categories: [savings, fire]   # 文字列 slug の配列
tags: [mail, freelance]        # 文字列 slug の配列
```

> カテゴリ/タグを新規追加・変更した場合は、`npm run build:redirects` を実行して `public/_redirects` を再生成・コミットしてください（旧URLのリダイレクト維持のため）。
