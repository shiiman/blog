---
id: 659
title: 【WordPress】おすすめプラグイン part 6 〜 Copy Anything to Clipboard 〜
slug: copy-anything-to-clipboard
status: publish
excerpt: こんばんは、しーまんです！ 皆さんwebサイトを見ていて「こちらのコードをコピー」とか「こちらを貼り付けてください」という箇所をご覧になったことはないでしょうか。 クーポンなどをコピーして貼り付けさせることで特典がもらえ \[…\]
categories:
    - 10
    - 3
tags:
    - 6
    - 45
featured_media: 660
date: 2021-09-18T19:30:00
modified: 2021-09-14T15:11:34
---

こんばんは、しーまんです！

皆さんwebサイトを見ていて「こちらのコードをコピー」とか「こちらを貼り付けてください」という箇所をご覧になったことはないでしょうか。

クーポンなどをコピーして貼り付けさせることで特典がもらえたりする場合がありますね。

また、技術記事などではプログラミングコードをよく「コピーして使用してください」という場合が多々あります。

こういった場合にいちいちドラッグして選択するのってめんどくさいですよね。

今回はこの文章中の文字列をワンクリックでコピーさせる事ができる「Copy Anything to Clipboard」というプラグインを紹介します！

## Copy Anything to Clipboard

今回紹介する「Copy Anything to Clipboard」は必須プラグインではありませんが、私のように技術記事を書いていこうというブロガー様は導入を検討してもよいかもしれません。

それは記事中にプログラムコードが多く登場しますので、ユーザビリティの観点からつけておいた方がよい機能だからです。

### インストール

では早速インストールをしていきましょう！

WordPressのメニュー「プラグイン」-> 「新規追加」から「Copy Anything to Clipboard」を検索します。すると下記のプラグインが表示されますので、今すぐインストール、有効化していきましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard-800x366.png)](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard.png)

### 確認

では実際にコードブロックを作って確認してみましょう。

コードブロックを作る方法は「“\`」とうって改行する方法と、ブロック挿入ツールから選択する方法があります。なれた方は「“\`」と打つ方が早いのでおすすめですが、分からない方はブロック挿入ツールでブロックを作ってみましょう。

[![Copy Anything to Clipboard - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting6.png)](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting6.png)

実際にコードを入力してプレビューで確認すると以下のようにコードブロックの右に「コピー」ボタンがついているのが分かると思います。こちらを押して貰えればコピーがされるようになります。

[![Copy Anything to Clipboard - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting4-800x59.png)](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting4.png)

しかし、このままだとなんだかダサいですよね。。

ご安心ください。「Copy Anything to Clipboard」にはデフォルトでコピーボタンではなくコピーアイコンに変更することが可能です。

### コピーアイコンに変更

ではその設定をみていきましょう！

まず、WordPressの設定に「Copy to Clipboard」というメニューが追加されていると思いますので、こちらを選択します。すると一つだけ「pre」を書かれた項目があると思います。こちらを選択します。

[![Copy Anything to Clipboard - 設定3](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting-800x371.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting.jpg)

すると編集画面が開きますので一箇所だけ変更します。

名前が「style」の部分の値が「button」になっていると思いますので、これを「svg-icon」に変更します。変更できたら「保存」を押しましょう。

[![Copy Anything to Clipboard - 設定4](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting2-800x397.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting2.jpg)

以上で、コピーアイコンの設定は完了です。

### コピーアイコンの確認

設定ができたら改めてプレビューで確認します。

このようにコピーボタンからコピーアイコンに変わっていることが確認できるかと思います。

こちらの方がスタイリッシュですよね！好みの問題もありますのでお好きな方を使用してください。

[![Copy Anything to Clipboard - 設定5](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting5-800x57.png)](https://shiimanblog.com/wp-content/uploads/2021/09/copy_anything_to_clipboard_setting5.png)

## まとめ

今回は記事中のテキストをワンクリックでコピーできるようにする「Copy Anything to Clipboard」というプラグインを紹介しました。

このような小さな配慮がユーザビリティに関わってきますので、少しでもみてくれる方にとって便利な記事になればと思いおすすめさせていただきました。

私と同じように技術記事などをあげている方は導入を検討してみてはいかがでしょうか。