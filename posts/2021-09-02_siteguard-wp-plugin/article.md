---
id: 120
title: 【WordPress】おすすめのプラグイン Part 2 〜 SiteGuard WP Plugin 〜
slug: siteguard-wp-plugin
status: publish
date: 2021-09-02T19:30:00
modified: 2021-09-24T21:42:36
excerpt: WordPress管理画面のセキュリティを強化するプラグイン「SiteGuard WP Plugin」の導入手順と設定方法を解説します。
categories:
    - 10
    - 3
tags:
    - 6
    - 23
    - 25
featured_media: 121
---

みなさん自身のサイトのセキュリティー対策は万全でしょうか？

セキュリティーってなんか難しいイメージありますよね。

でもWordPressは結構セキュリティーホールになりやすいのでしっかり対策をしないと、個人情報を抜かれてしまったり、最悪サイトやブログを壊されてしまったりします。

そうならないためにも最低限のセキュリティー対策をしていきましょう。

今回はそのセキュリティー対策の第1歩としてSiteGuard WP Pluginというプラグインを紹介します。

## SiteGuard WP Pluginとは

SiteGuard WP Pluginは、不正なアクセスからwebサイトを守るためのプラグインです。

WordPressの管理画面は登録時 `https://xxx.com/wp-admin` または `https://xxx.com/wp-login.php` というurlでアクセス出来るようになっています。

このままの状態にしておくと管理画面への不正アクセスが容易にできてしまいます。

そのような脆弱性をこのプラグインを入れることでカバーができます。

## SiteGuard WP Pluginの設定手順

では早速プラグインを導入してみましょう！！

### プラグインの有効化

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でWordPressを導入した場合 SiteGuard WP Plugin は既にインストールされていました。

プラグイン画面を確認すると無効状態でインストール済みなのが分かると思います。

ここから有効化を押すとプラグインが有効化されます。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg)

もしSiteGuard WP Pluginがインストールされていない場合検索バーからインストールし、有効化することが可能です。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/3.png)](https://shiimanblog.com/wp-content/uploads/2021/09/3.png)

### 管理画面ログインURLの変更

プラグインを有効化すると直ぐに管理画面へのログインURLが変更になります。

こちらをメモしておかないと管理画面に入れなくなってしまうので必ずブックマークしてください！！

有効化するとプラグイン上部に「ログインページURLが変更されました。」と表示されます。

ここから直ぐに「新しいログインページURL」を選択しブックマークしてください。

これで `https://xxx.com/wp-admin` または `https://xxx.com/wp-login.php` では管理画面が開かなくなりました。新たなURLでログインすることになります。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/5-1-1024x307.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/5-1.jpg)

ブックマークが完了したら、「設定変更はこちら」をクリックします。

するとログインページ変更画面に切り替わります。

SiteGuard WP Pluginを使ってログインページを変更すると、初期値で `https://xxx.com/login_<5桁の乱数>` というパターンで設定されてしまいます。それを知っている攻撃者はそのパターンに当てはまるURLをひたすら試すといずれログインページに辿り着いてしまうので、セキュリティー上よくありません。

ですのでここで**もう一度ログインURLをご自身の手で変更**してください。

変更したら**再度ブックマーク**を忘れないでください。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot-1024x421.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Screenshot.jpg)

### その他の設定

ログインURLの変更が終わったらその他の設定をしておきます。

私の場合は更新通知の設定をOFFにしています。なぜならWordPressの更新やプラグインの更新はそれなりに頻度が高く、そのたびにメールで通知されてしまうからです。なるべくメールは意味のあるものだけ受け取るようにしないとメール自体見なかったり、見逃したりするのでこちらはOFFにしています。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/7-1024x545.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/7.jpg)

後は基本的にデフォルトの設定のままで問題ありません。

設定が完了したら、実際にログインをし直して管理画面のURLが変更されていることを確認しましょう。

WordPressの右上に表示されている「こんにちは、ユーザ名さん」にカーソルを当てるとログアウトが表示されますのでそこからログアウトをします。

ログアウトする前に確実に変更後のURLがブックマークされているか、もしくはメモをしたかを確認してください。

### 設定の確認

ログアウト後に、設定したURLにアクセスすると下記のようなログイン画面が表示されると思います。

するとユーザ名とパスワードの他に画像による文字入力の設定が追加されています。

こちらはボット対策です。この3種類を入力することにより今後はログイン可能になります。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/8.png)](https://shiimanblog.com/wp-content/uploads/2021/08/8.png)

## まとめ

今回はWordPressのセキュリティー対策として SiteGuard WP Plugin を導入してみました。

とりあえずこちらを導入すると、管理画面への不正アクセスを許す可能性を大幅に減少させることができます。もちろんセキュリティー対策はこれで完璧というものではありませんので、次回移行もセキュリティー関連のプラグインを紹介していきたいと思います。

皆さんの大事なブログやwebサイトを悪意のある攻撃者から守れるのは皆さん自身です。

少しずつでいいのでセキュリティーに対する意識を高めながら、対策をしていきましょう。