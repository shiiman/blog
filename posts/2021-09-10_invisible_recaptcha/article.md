---
id: 215
title: 【WordPress】おすすめのプラグイン Part 4 〜 Invisible reCaptcha 〜
slug: invisible_recaptcha
status: publish
date: 2021-09-10T19:30:00
modified: 2021-09-08T02:16:28
excerpt: WordPressのBot攻撃対策として「Invisible reCaptcha」プラグインの導入・設定方法を初心者向けに解説します。
categories:
    - 10
    - 3
tags:
    - 6
    - 23
    - 38
featured_media: 478
---

WordPressのコメントやお問い合わせなどのフォームでは、セキュリティー対策を何もしないとBotにより攻撃を受けやすいです。不正なコメントやお問い合わせが大量に送られてきたり、誹謗中傷を受けたりする場合もあります。

そこで今回紹介する Invisible reCaptcha を導入することにより、そのようなBotによる攻撃を **完全にシャットアウト**させることができます。

ぜひお問い合わせフォームなどを設置した際には、こちらの記事を参考にお早めに導入することをおすすめします！！

## Invisible reCaptchaとは

 [![](https://ps.w.org/invisible-recaptcha/assets/banner-772x250.png?rev=1560060)\
\
Invisible reCaptcha for WordPress\
\
Invisible reCaptcha for WordPress プラグインはあなたのサイトをGoogleの新しい Invisible reCaptcha によって悪意のあるスパムボットから守ります。\
\
![](https://www.google.com/s2/favicons?domain=https://ja.wordpress.org/plugins/invisible-recaptcha/)\
\
ja.wordpress.org](https://ja.wordpress.org/plugins/invisible-recaptcha/ "Invisible reCaptcha for WordPress")

Invisible reCaptchaとは、WordPressの各種フォームにGoogleが提供する認証システムを導入できるプラグインです。

こちら「コメント」だけでなく **Contact Form 7で作成した**「お問い合わせ」にも対応しているので、多くの方が使用しているBot対策プラグインです。

WordPressを導入すると最初から「Akismet Anti-Spam（アキスメット アンチスパム）」というプラグインがインストールされています。こちらもスパム対策のプラグインです。

違いとしては Akismet Anti-Spamが受け取ったスパムコメントやメールを自動判別して振り分けるのに対して、Invisible reCaptchaはそもそもスパムコメントやメールが送られないようすることが可能です。

ですので、私自身は Invisible reCaptcha の方を好んで使用しています。

### 認証システム用のキー取得

では早速設定していきましょう。

Invisible reCaptchaの導入にはGoogleの認証システムを使用するので、Googleでサイトを登録してキーを取得していきます。こちらはGoogleアカウントを持っていれば無料で取得可能です。

まずは [reCAPCHAの作成画面](https://www.google.com/recaptcha/admin/create) からキーを発行します。

以下の項目を入力しましょう。

- ラベル : 分かりやすいようにご自身のサイト名などにしておくと良いでしょう
- reCAPTCHAタイプ: reCAPTCHA v3 を選択
- ドメイン : ご自身のサイトのドメインを設定
- オーナー : メールアドレス
- 利用条件に同意にチェック

[![reCaptcha - 登録](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting.png)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting.png)

入力が完了したら、「 **送信**」をしましょう。

[![reCaptcha - 登録2](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting5-800x515.png)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting5.png)

送信が完了すると reCAPTCHAの **サイトキー** と **シークレットキー** が表示されます。

こちらをこの後Invisible reCaptchaに設定していきますので、控えておきましょう。

### プラグインのインストール

プラグインの新規追加から「Invisible reCaptcha」を検索します。

検索の結果から「今すぐインストール」を選択してインストールをして、インストールが完了したら有効化をしましょう。

[![Invisible reCaptcha - インストール](https://shiimanblog.com/wp-content/uploads/2021/09/plugin_reCAPTCHA-800x511.png)](https://shiimanblog.com/wp-content/uploads/2021/09/plugin_reCAPTCHA.png)

### Invisible reCaptchaの設定

有効化するとWordPressの左メニューの「設定」の中にInvisible reCaptchaという項目が追加されます。こちらから設定してきましょう。

先程GoogleのreCAPTCHA登録時に控えていたサイトキーとシークレットキーをそれぞれ登録します。

バッチ位置を今回はインラインとします。

入力が終わったら「変更を保存」しましょう。

[![Invisible reCaptcha - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting1-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting1.jpg)

次にInvisible reCaptchaの左メニューから WordPress を選択します。

こちらでどのフォームを対象にするかを選択できるので、問題なければ全てにチェックをいれます。

[![Invisible reCaptcha - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting2-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting2.jpg)

最後にお問い合わせフォームの設定です。

Contact Form 7のプラグインに対応しておりますので、こちらにチェックを入れて保存しましょう。

[![Invisible reCaptcha - 設定3](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting3-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting3.jpg)

ここまでで設定は完了になります。

### 設定が正しいか動作確認

正しく設定が反映されているかみていきましょう。

コメント欄のしたにreCAPTHAが表示されましたね。

[![reCAPTHA - 確認1](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting4-800x555.png)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting4.png)

Contact Form 7で作成したお問い合わせフォームの画面も確認しましょう。

送信ボタンの下にreCAPTHAが表示されているのが確認できるかと思います。

[![reCAPTHA - 確認2](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting5-1.png)](https://shiimanblog.com/wp-content/uploads/2021/09/reCAPTCHA_setting5-1.png)

## まとめ

今回はInvisible reCaptchaというプラグインを導入し、コメントやお問い合わせなどのフォームからBotが不正なメッセージを送れないようにする対応をしました。

こちらの対応をするだけでかなりの攻撃を防ぐことが可能です。

下記の画像は私のブログのリクエスト数を表しています。

サイトを登録したのが2021/08/30 なので登録したばかりなのですが、リクエストが何件かあるのが確認できます。

[![BOTと判定されたリクエスト](https://shiimanblog.com/wp-content/uploads/2021/09/attach_request-800x523.png)](https://shiimanblog.com/wp-content/uploads/2021/09/attach_request.png)

こちらのグラフはGoogleの方に登録した [reCaptchaの管理画面](https://www.google.com/recaptcha/admin) から確認可能です。

これからアクセス数が伸びればもっとBot判定されるリクエストは増えていくでしょう。

このようなリクエストから Invisible reCaptcha は守ってくれますので、ぜひ導入をしてみてください。