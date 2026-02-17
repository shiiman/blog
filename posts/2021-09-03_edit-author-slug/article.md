---
id: 152
title: 【WordPress】おすすめのプラグイン Part 3 〜 Edit Author Slug 〜
slug: edit-author-slug
status: publish
date: 2021-09-03T19:30:00
modified: 2021-09-02T14:23:19
excerpt: WordPressのセキュリティ対策として、ユーザー名の漏洩を防ぐプラグイン「Edit Author Slug」の導入と設定方法を紹介します。
categories:
    - 10
    - 3
tags:
    - 6
    - 23
    - 26
featured_media: 145
---

WordPressを初めてセキュリティー対策としてプラグインを追加すると思いますが、

今回は「Edit Author Slug」というものを紹介します。

皆さんWordPressでブログやwebサイトを運用していて、

悪い人が一番手軽に悪さをする方法は何だと思いますか？

それは管理画面にログインすることです。

管理画面にログインさえできてしまえば、情報の抽出だったり、サイト破壊ということも容易にできてしまいます。

ではログインさせないようにするにはどうしたら良いでしょうか。

簡単なのはログイン時に使用するユーザ名とパスワードがバレなければ大丈夫です。

しかしWordPressを立ち上げてすぐの場合、デフォルト設定だとユーザ名は第三者にバレてしまいます。そこで役に立つのが今回紹介する「Edit Author Slug」というプラグインです。

## Edit Author Slugとは

WordPressを立ち上げて直ぐにブラウザの入力欄に

`https://xxx.com/?author=1` と入力すると

`https://xxx.com/authhor/[ユーザ名]`

というURLにリダイレクトされてしまいます。

するとリダイレクトされたURLにユーザ名が入っているので、

第三者でも簡単にユーザ名が分かってしまうというわけです。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/2-1.png)](https://shiimanblog.com/wp-content/uploads/2021/09/2-1.png)

↓↓ リダイレクトされる

[![](https://shiimanblog.com/wp-content/uploads/2021/09/3-1.png)](https://shiimanblog.com/wp-content/uploads/2021/09/3-1.png)

というわけで、上記URLを入力されてもユーザ名がバレないように設定する必要があります。

それが簡単に出来るのがEdit Author Slugというプラグインです。

## Edit Author Slugの設定手順

では早速プラグインを導入してみましょう！！

### プラグインのインストール

インストールは簡単でプラグインの新規追加入力画面から「Edit Author Slug」を検索すると画像のように見つかりますので、そこからインストールをしてください。

インストールが終わったら忘れずに有効化しましょう！

[![](https://shiimanblog.com/wp-content/uploads/2021/09/1.png)](https://shiimanblog.com/wp-content/uploads/2021/09/1.png)

### Edit Author Slugの設定

プラグインの導入が終わったら ユーザ -\> プロフィール の画面に「Edit Author Slug」の項目が追加されます。そのカスタム設定に変更したいユーザ名に変更してください。

入力が終わったらプロフィールを更新ボタンを押して設定を反映させましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-1024x341.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot.png)

### 設定が正しいか動作確認

これで設定が完了したので、動作が期待通りか確認しましょう。

再度 https://xxx.com/?author=1 にアクセスします。

するとリダイレクトされたURLが変更されたのが確認できます。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/7-1.png)](https://shiimanblog.com/wp-content/uploads/2021/09/7-1.png)

## もう一つ設定する項目

Edit Author Slug プラグインの設定は以上になります。

しかし、もう一箇所設定しないとまずそうな箇所が残っています。

先程確認しましたが、リダイレクト後のURLにユーザ名が表示されてしまう問題は解決しました。

しかしブラウザのタブの方を確認してみましょう。

こちらにもユーザ名が表示されています。

こちらの表示の修正はプラグインを使用せずに修正が可能ですので、一緒に設定しておくと良いでしょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/5-2.png)](https://shiimanblog.com/wp-content/uploads/2021/09/5-2.png)

こちらの設定変更も簡単です。

先程と同じで ユーザー -\> プロフィール 画面を開きます。

デフォルトだとユーザ名がタブに表示されてしまいますが、変更はできません。

なのでニックネームを設定してあげましょう。

そして設定したニックネームを「ブログ上の表示名」に設定すると設定完了です。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-1-1024x394.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-1.png)

## まとめ

今回は Edit Author Slug というプラグインを使用してユーザ名がバレないように設定しました。

設定自体は簡単ですが、色々とセキュリティーを意識しないといけないなぁと感じた方もいらっしゃるのではないでしょうか。

セキュリティーといっても何をやったら分からないという方が大半だとは思いますが、今回のように少しずつ対策をしていけばそうそう問題になることはありません。

最初の導入時に意識して設定し、後は定期的に見直すようにしていきましょう。