---
id: 144
title: 【Cocoon】サイトマップをヘッダー・フッターに追加する方法 / 初心者でも簡単！5分でサイトマップ完成
slug: sitemap
status: publish
date: 2021-09-05T19:30:00
modified: 2021-09-02T18:15:20
excerpt: Cocoonを使用してWordPressにサイトマップを作成し、ヘッダーやフッターに設置する方法を初心者向けに解説します。
categories: [4, 11, 3]
tags: [6, 27, 31]
featured_media: 247
---

WordPressでブログを書いていると最初のうちは記事数も少ないのでいいですが、

だんだんページ数が多くなってくると目的のページを見つけるのが大変になってきます。

またGoogleなどの検索エンジンのクローラーに対してサイト内のページを循環して正しく伝えられるようにする効果もあります。

クローラーとは検索エンジンが検索順位を決めるための要素を、サイトから収集するロボットのことです。

そこで今回はWordPressのテーマであるCocoonでサイトマップをあっという間に作成する方法を紹介致します。

## サイトマップの作成

サイトマップはWordPressの固定ページに作っていきます。固定ページとは普段執筆する投稿とは違い、ヘッダーやフッターなどに常に表示されているページを作成する際に使用します。

固定ページの上部の新規追加から固定ページを作っていきます。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/5-3-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/5-3.jpg)

タイトルに「サイトマップ」

本文に下記ショートコードを入力

```
[sitemap]
```

それだけでcocoonを使用している場合いい感じのサイトマップを作成してくれます。

凝ったものを作成したい場合を除きこちらで十分です。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-2-800x138.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-2.png)

それでは画面を保存をして画面を確認しましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/sitemap.png)](https://shiimanblog.com/wp-content/uploads/2021/09/sitemap.png)

こんな感じで「固定ページ」「投稿一覧」「カテゴリー」と分類して表示してくれます。

もちろん記事やカテゴリーが増えれば自動で追加されます。

## ヘッダーの作り方

サイトマップの固定ページができたのでこれを画面トップのヘッダーと画面下のフッターに追加していきましょう。ヘッターとフッターを追加するとグッとサイトの使い勝手がよくなります。

外観 -\> メニュー から新しくメニューを作っていきます。

今回はヘッダーを作るので

メニュー名: ヘッダーメニュー

メニューの種類：ヘッダーメニュー

として保存します。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/2-3-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/2-3.jpg)

メニューが作成できたら項目を追加していきましょう。

固定ページの「すべて表示タブ」から必要なページを選択してメニューに追加します。

すると「メニュー構造」に追加されるので、好きな順番に並べかえます。そして「メニュー設定」でヘッダーメニューを選択し保存します。

参考画像ではホームと今回作成したサイトマップを追加しています。

ホームはサイトのトップページに遷移するメニューになります。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-1-1-800x486.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-1-1.png)

設定ができたら実際に画面で確認してみましょう。

サイトの上部にヘッダーメニューが追加されていることが分かります。

今回はサイトマップしか追加しておりませんが、必要な項目があったら今回作成したメニューに追加していきましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/header-800x367.png)](https://shiimanblog.com/wp-content/uploads/2021/09/header.png)

## フッターの作り方

次にヘッターと同じように 外観 -\> メニューからフッターメニューを追加していきましょう。

メニュー名: フッターメニュー

メニュー設定: フッターメニュー

を選択して保存します。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/1-1-2-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/1-1-2.jpg)

そして固定ページから対象のページを選択しメニュー構造でフッターメニューを選択します。

忘れずに保存を行いましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-3-800x483.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-3.png)

するとサイト下部にフッターメニューが表示されているのが確認できます。

フッターの位置はデフォルトでは右端に表示されていると思います。

参考画像と同じように中央に設定したい場合はCocoon 設定のフッタータブから変更可能です。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/footer-800x395.png)](https://shiimanblog.com/wp-content/uploads/2021/09/footer.png)

## まとめ

今回はWordPressの人気無料テーマであるCocooでサイトマップを作成してみました。

また、作成したサイトマップをヘッターとフッターに表示させるところまで行いました。

ヘッターとフッターを作ることで、かなりサイトっぽい見た目になってきましたね。

ヘッターとフッターに表示する内容は好みで決めていただいて問題ありません。

ただしサイトマップの他にお問い合わせやプライバシーポリシーは鉄板で入れておきたい固定ページになりますので、ぜひ追加してみてください。