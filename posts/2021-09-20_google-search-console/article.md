---
id: 701
title: 【2021年9月】初心者用 Google Search Console 登録と設定 | Cocoon
slug: google-search-console
status: publish
excerpt: こんばんは、しーまんです！ Google Search Console(通称サチコ)はwebサイトの表示回数やユーザの行動を分析するためのツールです。こちらのツールを使用することで、自身のwebサイトの状態をグラフで確認 \[…\]
categories:
    - 13
tags:
    - 47
featured_media: 703
date: 2021-09-20T19:30:00
modified: 2021-09-21T18:32:40
---

こんばんは、しーまんです！

Google Search Console(通称サチコ)はwebサイトの表示回数やユーザの行動を分析するためのツールです。こちらのツールを使用することで、自身のwebサイトの状態をグラフで確認することができます。

webサイトを運営している方は必ずといっていいほど導入しているツールです。

無料で使用できますので、そのような便利なツールを活用できるように、本記事では登録・設定方法を解説します！

## Google Search Console

まずGoogle Search Consoleを使用するにはGoogleアカウントの登録が必要です。

こちらがまだの方は [アカウント作成ページ](https://accounts.google.com/signup/v2/webcreateaccount?flowName=GlifWebSignIn&flowEntry=SignUp) より作成お願いします。

### サイトの登録

アカウントの用意ができたら早速Search Consoleの登録をしていきましょう。

[登録ページ](https://search.google.com/search-console/about?hl=ja) にアクセスして「今すぐ開始」をクリックします。

[![Google Search Console - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console0-800x329.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console0.png)

まずはプロパティタイプの選択とご自身のサイト登録の画面に遷移します。

「ドメイン」と「URL プレフィックス」が選択できますが、この後の認証設定が変わります。

「ドメイン」の場合はご自身の契約しているサーバのDNS設定にTXTレコードを登録することになります。こちらの手順は登録しているサーバによって変わりますので、少し難しいです。

「ドメイン」での設定が難しそうな方は「URL プレフィックス」を選択しましょう。

今回はこちらで説明をしていきます。

[![Google Search Console - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console1-1-800x559.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console1-1.png)

「URL プレフィックス」を選択すると次は所有権の確認ページになります。

ここでは「その他の確認方法」->「HTMLタグ」を表示させましょう。

すると表示されたHTMLタグをご自身のページに貼り付けてくださいと出てきます。

Cocoonではこの設定が簡単に出来るように設定されています。

まずはHTMLタグの中から「content=””」で表示された値のみコピーしてください。

[![Google Search Console - 設定3](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console2-800x788.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console2.png)

コピーができたらWordPressの管理画面を開きます。

Cocoon設定から「アクセス解析・認証」タブを開くと真ん中あたりに「Google Search Console設定」という項目が現れます。

そこに「Google Search Console ID」という項目がありますので、先程コピーしたcontentの内容をコピーします。

[![Google Search Console - 設定4](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console3-1-800x192.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console3-1.png)

WordPress側の設定が完了したらSearch Consoleの画面に戻り「確認」を押していきましょう。

正しくIDを設定できていると「所有権を自動確認しました」というポップアップが表示されます。

こちらが表示されればサイトの登録は完了です。

[![Google Search Console - 設定5](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console4-800x385.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console4.png)

### サイトマップの登録

Google Search Consoleへのサイト登録が終わったら、サイトマップの登録も合わせておこなっておきましょう。登録しておくことでサイトの分析を取り込めるようになります。

サイトマップってなに？という方や、まだサイトマップを作成していない方はこちらの記事を参考にサイトマップの作成をしてみてください。

[https://shiimanblog.com/wordpress/jetpack/](https://shiimanblog.com/wordpress/jetpack/)

設定方法は簡単です。Google Search Console画面の左メニューから「サイトマップ」を開きます。

新しいサイトマップの追加でご自身のサイトのサイトマップURLを入力して送信しましょう。

JetpackやGoogle XML Sitemapsのようなプラグインを使用してサイトマップを作成している場合、デフォルトでは「sitemap.xml」にファイルが作られる設定のはずです。

サイトマップのファイルの場所はご自身の環境を確認して入力しましょう。

[![Google Search Console - 設定6](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console5-1-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_serch_console5-1.jpg)

以上で、Google Search Consoleの設定は完了です。

あとは日々データが収集されますので、ご自身のサイトのアクセス数やクリック数、流入経路などを確認して、サイトの最適化に役立てましょう。

## まとめ

今回はご自身のwebサイトの状態を確認するツールであるGoogle Search Consoleの設定方法について解説しました。

無料で使用できる強力なツールですので、ぜひ導入していない方は直ぐに設定をしてみてください。

こちらのツールは導入してからのデータしか保存されないためできるだけ早く導入することをおすすめします。

Google Search Consoleを使用して、webサイトの状態を常に把握し、どんどん改善してよりよいwebサイトにしていきましょう！