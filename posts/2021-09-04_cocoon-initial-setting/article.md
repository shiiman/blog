---
id: 217
title: 【Cocoon】WordPressでテーマ導入してみよう / 初期設定もバッチリ紹介！
slug: cocoon-initial-setting
status: publish
excerpt: WordPressを立ち上げて、ある程度セキュリティーの設定をしたら、サイトの見た目を整えたくなりますよね。 今回はCocoonというテーマを使って、サイトデザインを行っていきます。この記事を読むことでCocoonの初期 \[…\]
categories:
    - 4
    - 11
    - 3
tags:
    - 6
    - 23
    - 31
featured_media: 218
date: 2021-09-04T19:30:00
modified: 2021-09-24T21:43:39
---

WordPressを立ち上げて、ある程度セキュリティーの設定をしたら、

サイトの見た目を整えたくなりますよね。

今回はCocoonというテーマを使って、サイトデザインを行っていきます。

この記事を読むことでCocoonの初期設定が分かるようになります。

では早速みていきましょう！

## テーマの反映

まずはじめにテーマを反映させる前(デフォルトのテーマ)の状態はこのような感じです。

だいぶシンプルで、これではブログに成り立っていないですよね。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/12659675-ba0e6df566ec04d496eb332322900971-800x635.png)](https://shiimanblog.com/wp-content/uploads/2021/08/12659675-ba0e6df566ec04d496eb332322900971.png)

そこでWordPressの画面から 外観 -> テーマ と選択し、このページからテーマを反映させていきます。

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でテーマCocoonを選択した場合既にCocoonがインストールされています。

Cocoonがインストールされていない場合は画面上部の新規追加から画面移動し、

そのページの上部にあるテーマのアップロードからCocooのテーマをアップロードします。

テーマ自体は下記の公式ページから無料でダウンロード可能です。

[https://wp-cocoon.com/](https://wp-cocoon.com/)

[![](https://shiimanblog.com/wp-content/uploads/2021/09/1-1-1-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/1-1-1.jpg)

テーマには親テーマと子テーマがあり、Cocoonに関わらず、子テーマの方を反映しましょう。

今回でいうとCocoon Child という方を選択肢有効化してください。

子テーマの方を有効化しないとテーマがアップデートされた際に設定が上書きされてしまい、設定した内容が消えてしまう場合があります。

Cocooを有効化すると下記のようなデザインになります。

グッとよくあるwebサイトの見た目に近づきましたね。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/2-2-800x425.png)](https://shiimanblog.com/wp-content/uploads/2021/09/2-2.png)

## スキンを選択

テーマが有効化できたら、左メニューにCocoon 設定の項目が追加されます。

これを選択して、まずは「スキン」タブを開きましょう。

すると既に色々な方が作成したスキンの一覧が表示されます。

お好みのスキンを探して保存してみましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/3-2-800x1561.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/3-2.jpg)

私は [tecurio moon](https://tecurio.com/) _\[作者: [風塵(ふーじん)](https://tecurio.com/)\]_ というスキンを選択してみました。

こちらを反映するとこのような表示になります。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/4-800x560.png)](https://shiimanblog.com/wp-content/uploads/2021/09/4.png)

スキンを設定しただけでCocooの多くの設定がいい感じに調整されます。

一度時間がある時にどんな設定があるのか確認してみても面白いと思います。

## ホームイメージ画像の設定

スキンを設定したら基本的には設定する項目はなく使用可能です。

ただこのままだとSNSなどにご自身のURLを共有した際にホームイメージというものが表示されることがあるのですが、それがCocooの画像が表示されることになってしまいます。

こちらはSlackというチャットツールに展開してみた例になりますが、自身のサイトなのにCocooのイメージが表示されているのが分かります。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/61d9a17b4171923f9e96d960538dabe7.png)](https://shiimanblog.com/wp-content/uploads/2021/09/61d9a17b4171923f9e96d960538dabe7.png)

流石ににこれは変更しておきましょう。

OGPタブを選択すると一番下のホームイメージという箇所があります。

そこから画像をアップロードして保存することでホームイメージの画像を変更することができます。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/42527942e011382bc15c751bc923ed2a-800x430.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/42527942e011382bc15c751bc923ed2a.jpg)

実際に設定が反映されると下記のように展開されるようになります。

必ずご自身の設定したイメージを表示されるように設定しておきましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/a3deaa49ae8c0fd3f6738d0095e7f41c.png)](https://shiimanblog.com/wp-content/uploads/2021/09/a3deaa49ae8c0fd3f6738d0095e7f41c.png)

## 高速化

次はデザイン周りではなく高速化の設定を行っていきます。

こちらの設定は必須ではありません。

またお使いの環境によっては不具合が生じることがございますのでご注意を

WordPressの左メニューから Cocoon設定 -> 高速化 を選択します。

そして以下5つの項目にチェックを入れましょう。

- ブラウザキャッシュ
- HTML縮小化
- CSS縮小化
- JavaScript縮小化
- Lazy Loadの設定 遅延読み込み

最初の方はブログ記事も少ないので問題になりませんが、記事が増えてくるとwebサイトの表示が遅く感じることがあります。その際にこの設定が効いてきます。

ブラウザキャッシュを設定すると更新されたページが反映されない(反映されるのに時間がかかる)場合があるので注意してください。

縮小化系はお使いのスキンによっては不具合を起こす可能性があります。

私の場合はCSS縮小化を行うと画像が表示されなくなるという不具合が発生しましたので、CSS縮小化のみ無効にしています。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/6-800x127.png)](https://shiimanblog.com/wp-content/uploads/2021/09/6.png)[![](https://shiimanblog.com/wp-content/uploads/2021/09/7-2-800x445.png)](https://shiimanblog.com/wp-content/uploads/2021/09/7-2.png)[![](https://shiimanblog.com/wp-content/uploads/2021/09/8-800x265.png)](https://shiimanblog.com/wp-content/uploads/2021/09/8.png)

## まとめ

今回はWordPressの無料で使用できるテーマであるCocoonを導入し、最初にやるべき最低限の設定を行いました。簡単な設定だけで見た目がきれいなwebサイトが構築できましたね。

またCocoonは細かく設定が分かれているので、色々と自分好みに調整することができます。

慣れてきたらカスタマイズに挑戦してもいいかもしれませんね。