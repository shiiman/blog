---
id: 671
title: 【WordPress】おすすめプラグイン part7 〜 Jetpack 〜
slug: jetpack
status: publish
date: 2021-09-19T19:30:00
modified: 2021-09-14T17:30:50
excerpt: WordPressの記事投稿時にSNSへ自動連携するプラグイン「Jetpack」の導入・設定方法を紹介します。
categories: [10, 3]
tags: [6, 23, 46]
featured_media: 672
---

こんばんは、しーまんです。

皆さん記事を投稿した際に自動でSNSに告知したいと思ったことはありませんか。

通常だとwebサイトの更新後、そのURLをSNSに貼り付けて告知という感じで行いますよね。

ですが、手作業でこの作業を行うと時間がかかるうえにミスをする可能性もあります。また、複数のSNSを使用している方であればその手間も倍増します。

私はブログを更新したらSNSにも自動で通知いくようにしたいなと思い、そんなプラグインがないのかなと探していたところ「Jetpack」というプラグインを見つけました。

本記事ではそんなJetpackを使用してSNSへ自動通知する方法と、Jetpackのその他機能について説明していきます！

## Jetpack

JetpackとはWordPress.comのプラグインで、SNS自動通知以外にもたくさんの機能が備わっています。下記が設定できる項目の分類になります。

- セキュリティ
- パフォーマンス
- ライティング
- 共有
- ディスカッション
- トラフィック

ただし、**Jetpackは高機能ゆえに少し重かったり、他のプラグインと競合するものも多い**です。

私の場合は「共有」設定にあるSNS自動通知機能と「トラフィック」設定にあるサイトマップの作成機能のみを使用しておりますので、今回はこの2つの設定について紹介します。

その他の機能でも便利なものがありますので、こちらについてはご自身の使用用途に合わせて設定してみてください。

Jetpackを使用するのは無料で利用できますが、WordPress.comのアカウントが必要になります。なので登録手順から紹介していきます。

### WordPress.comアカウント設定

まずはWordPressの設定からプラグイン -> 新規追加から「Jetpack」を検索します。

検索に表示されたプラグインから今すぐインストール、有効化していきましょう。

[![Jetpack - 登録1](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack-800x494.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack.png)

するとWordPressの設定画面に遷移しますので、まずはWordPress.comアカウントを作成します。

「Jetpackを設定」をクリックして遷移しましょう。

[![Jetpack - 登録2](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting1-800x734.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting1.jpg)

Googleアカウントでの登録ができますので、使用するGoogleアカウントが選択されていたら「承認」しましょう。もし別のアカウントが選択されていたら「アカウントを切り替え」します。

![Jetpack - 登録3](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting2-800x902.png)

アカウント選択した後はプランの選択画面に遷移します。

最初は「Free」プランを選択しましょう。

[![Jetpack - プラン選択](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack.com_.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack.com_.png)

プラン選択の後はアンケート的なものが続きますので、どんどん次に進みましょう。

特に設定しなくても問題ないので「後で」を選択して問題ありません。

[![Jetpack - 登録4](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting3-800x578.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting3.jpg)[![Jetpack - 登録5](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting4-800x554.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting4.png)[![Jetpack - 登録6](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting5-800x531.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting5.jpg)[![Jetpack - 登録7](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting6-800x543.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting6.png)[![Jetpack - 登録8](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting7-800x558.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting7.png)[![Jetpack - 登録9](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting8-800x662.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting8.jpg)

こちらでJetpackの初期設定は終了です。

### SNS連携

登録が終わりましたら、早速SNS連携をしていきましょう。

2021年9月現在連携可能なSNSは以下になります。

- Facebook
- Twitter
- LinkedIn
- Tumblr
- Google Photos
- MailChimp
- Instagram

今回はTwitter連携の手順を紹介します。

他のSNS連携も同じ手順でできますので、ご自身でお試しください。

まずは共有タブから「ソーシャルネットワークに自動共有」をONにします！

[![Jetpack - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting9-800x293.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting9.png)

設定をONにすると「ソーシャルメディアアカウントを接続する」リンクが表示されますのでこちらをクリックします。

[![Jetpack - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting10-800x328.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting10.png)

するとWordPress.comの管理画面に移動しますので、連携タブのTwitterの横にある「連携」ボタンを選択しましょう。

[![Jetpack - 設定3](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting11-800x420.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting11.jpg)

すると連携中のTwitterが表示されますので正しいアカウントであれば、「接続」をクリックします。

[![Jetpack - 設定4](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting12-800x492.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting12.png)

WordPress.comとTwitterの連携承認画面が表示されますので、最後にこちらの内容で問題なければ「連携アプリを承認」をクリックしましょう。

[![Jetpack - 設定5](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting13-1-800x765.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting13-1.jpg)

こちらでWordPressで記事を投稿するとTwitterにも自動で通知が送られるようになります。

記事投稿された後にTwitterで自動通知されることを確認しましょう。

### サイトマップ作成

次はサイトマップの生成の設定を行います。

サイトマップというのは、Googleなどの検索エンジンがご自身のwebサイトのインデックスを作成しています。そのインデックスを収集する役割のロボット(クローラー)が、インデックスを作成しやすくする働きをするものがご自身のサイトのサイトマップです。

サイトマップを作成して、検索エンジンに通知することで、ご自身のサイトのインデックスがしっかりと検索エンジン側に伝えることが可能です。これはSEO観点でも非常に重要ですので、サイトマップの生成と通知は設定しておくようにしましょう。

私はJetpackを使用する前は「Google XML Sitemaps」というプラグインを使用して、サイトマップを使用していました。「Google XML Sitemaps」のようにサイトマップを設定するプラグインは数多く存在しますので、どれか1つで必ず設定するようにしましょう。

それでは「Jetpack」でのサイトマップの設定をみていきます。

サイトマップの設定はとても簡単です。

Jetpackのトラフィックタブを選択するとサイトマップという項目があります。

そこに「XMLサイトマップ生成する」がありますので、こちらをONにするだけです。

こちらをONにするだけでサイトマップが更新されるたびに検索エンジニに自動で通知してくれます。

[![Jetpack - 設定6](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting14-800x282.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Jetpack_setting14.png)

## まとめ

今回はJetpackというプラグインを使用して、SNSへの自動通知とサイトマップの作成について解説しました。

SNSの通知は今の時代たくさんの方にご自身の記事を読んでもらうためには必須の流れかなと思います。また、サイトマップの作成もSEO観点から必ずやってほしい設定になります。

今回紹介した機能以外にもJetpackは画像の遅延ロード機能や、総当たり攻撃からの保護機能、アクセス解析機能など、便利な機能が豊富に存在します。

興味がありましたら、こちらも合わせて設定してみてください。