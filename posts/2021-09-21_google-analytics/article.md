---
id: 716
title: 【2021年9月】初心者用 Google Analytics 登録と設定 | Cocoon
slug: google-analytics
status: publish
excerpt: こんばんは、しーまんです！ 今回はGoogle Analyticsというツールを紹介していきたいと思います。 Google AnalyticsとはGoogleが提供する無料で高機能なアクセス解析ツールです。webサイトを \[…\]
categories:
    - 13
tags:
    - 48
featured_media: 719
date: 2021-09-21T19:30:00
modified: 2021-09-15T01:50:28
---

こんばんは、しーまんです！

今回はGoogle Analyticsというツールを紹介していきたいと思います。

Google AnalyticsとはGoogleが提供する無料で高機能なアクセス解析ツールです。

webサイトを運営している方にとっては必須ともいえるツールになります。

同じようなweb運営の必須ツールにGoogle Search Consoleがあります。

こちらは別の記事で導入方法を解説していおりますので、合わせて御覧ください。

 [![](https://shiimanblog.com/wp-content/uploads/2021/09/eyecatch_google_search_console-320x180.png)\
\
【2021年9月】初心者用 Google Search Console 登録と設定 \| Cocoon\
\
WordPressでwebサイトを構築した際に、どのような検索ワードで自身のサイトにアクセスしたのかを調査したい場合Google Search Consoleを使用します。\
そのGoogle Search Consoleの登録と設定を解説していきたいと思います。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.09.20](https://shiimanblog.com/wordpress/analytics/google-search-console/ "【2021年9月】初心者用 Google Search Console 登録と設定 | Cocoon")

## Google Analytics

### Google Analyticsとは

Google Analyticsは、webサイトのアクセスを分析するためのツールです。ご自身のサイトにアクセスしているユーザの動向などが分かるため、SEOの改善などにも役立てることができます。

Google Analyticsでは多くの機能がありますが、例えば次のような分析が可能です。

- どこの地域からのアクセスが多いかを分析
- どの記事へのアクセスが多いかを分析
- ブラウザやPC/スマホのアクセス割合を分析

Google Search Consoleとの違いはGoogle Search Consoleがユーザがどのように検索して自身のwebサイトにアクセスしたのかを把握するツールなのに対して、Google Analyticsはアクセスしたユーザがどのように分類されるのかを分析するツールです。

両方とも重要な情報を無料で得ることが可能ですので、合わせて設定して分析に使用することをおすすめします。

### Google Analyticsの登録

まずGoogle Analyticsを使用するにはGoogleアカウントの登録が必要です。

こちらがまだの方は [アカウント作成ページ](https://accounts.google.com/signup/v2/webcreateaccount?flowName=GlifWebSignIn&flowEntry=SignUp) より作成お願いします。

Googleアカウントの登録が完了したら、 [Google Analyticsの登録ページ](https://analytics.google.com/analytics/web/?hl=ja) にアクセスしましょう。

[![Google Analytics - 登録1](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics0-800x497.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics0.png)

まずはアカウント名を設定します。

1つのアカウントで複数のwebサイトの登録が可能ですので、webサイトの名前とは別の名前を登録しましょう。

[![Google Analytics - 登録2](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics1-1-800x429.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics1-1.jpg)

次の項目でwebサイトの情報を入力していきます。

プロパティ名にご自身のwebサイト名を入力しましょう。

続いてレポートのタイムゾーンと通貨を設定します。日本の場合は「日本」「円」を登録します。

[![Google Analytics - 登録3](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics2-800x429.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics2.png)

「詳細オプションを表示」を押すと、ユニバーサルアナリティクスプロパティの作成画面が開きます。

Google Analyticsは現在「 **Googleアナリティクス4プロパティ(GA4プロパティ)**」と「 **ユニバーサルアナリティクスプロパティ(UAプロパティ)**」の2種類が存在しています。デフォルトではGA4プロパティのみが作成されるようになっており、GA4プロパティの方が新しくできた分析測定方式になります。

ただ、まだGA4プロパティの方はUAプロパティに比べ分析できる項目も少なく、完全に移行できているかといったらできていないのが現状です。

ですので分析を今までの方式で行いたい方や、細かい分析を行っていきたい方は、こちらでUAプロパティを作成し、 **両方のプロパティを使用することを強く推奨** します。

※ UAプロパティは後から追加で作成することも可能ですので、よく分からない方はそのまま進んでいただいて問題ありません。

最後に「業種」「ビジネス規模」「利用目的」を選択して作成を押します。

[![Google Analytics - 登録4](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics4-800x429.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics4.jpg)

作成後はデータストリームの作成を行います。

Choose a platform でウェブを選択し、その後ご自身のwebサイトのURLとストリーム名を登録しましょう。ストリーム名はご自身のwebサイトの名前で問題ありません。

[![Google Analytics - 登録5](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics5-1-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics5-1.jpg)[![Google Analytics - 登録6](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics6-1-800x425.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics6-1.jpg)

ここまででGoogle Analyticsの大枠の登録が完了です。

このあと、サイト側に測定するための設定をおこなっていきます。

### WordPress側にトラッキングコードの埋め込み

Google Analytics側の設定が終わったら、WordPressにトラッキングコードの埋め込みを行っていきます。こちらをご自身のwebサイトに埋め込むことでwebサイトのいろいろな情報を取得して分析できるということです。

まずはトラッキングコードの取得から行っていきましょう。

データストリームの登録が終わるとウェブストリームの詳細が表示されていると思います。

そこから「グローバル サイトタグ（gtag.js） ウェブサイト作成ツールや、CMS でホストされるサイトをご使用の場合、このタグを設定」をクリックします。

[![Google Analytics - 登録7](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics6-800x409.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics6.jpg)

クリックするとサイトごとに個別のトラッキングコードが表示されますのでこちらのタグをコピーしておきます。

[![Google Analytics - 登録8](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics8-800x335.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics8.png)

### Cocoon設定でトラッキングコード埋め込み

今回はWordPressのテーマCocoonの設定を利用してトラッキングコードを埋め込んでいきます。

Cocoonを使用することで簡単にトラッキングコードを埋め込むことが可能です。

その他のテーマを使用している方は「外観」-\> 「テーマエディタ」でheader.phpを編集し、トラッキングコードに埋め込んでください。

WordPressのメニューから「Cocoon 設定」を選択し、「アクセス解析・認証」タブを選択します。

その画面で「ヘッド用コード」と入力出来る箇所がありますので、先程コピーしたトラッキングコードを貼り付けます。貼り付けできたら忘れずに保存しましょう。

[![Google Analytics - 登録9](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics9-1-800x448.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/google_analytics9-1.jpg)

こちらでWordPress側の設定も完了です。

Google Analyticsにデータが送られるようになったと思います。

実際にご自身のサイトにアクセスして、その後Google Anlyticsのデータが更新されていることを確認しましょう。

## まとめ

今回はGoogleが提供する無料で高機能なwebサイト分析ツールであるGoogle Analyticsの登録手順を解説しました。少しコードを埋め込んだりするところが複雑ですが、webサイトの運営には必須のツールですので、ぜひゆっくり記事を見直して導入してみてください。

また、GA4プロパティとUAプロパティのところは初心者にとって難しいところでしょう。

これは新しい分析ツールの移行に伴い並行して利用可能になっているので分かりづらくなっています。

ただし今後GA4プロパティの方が更新されていくことはGoogleからも正式に発表されています。

GA4プロパティの分析機能が十分になれば将来的にGA4プロパティだけを使用すればよくなるでしょう。

今回の記事で少しでもGoogle Analytics導入の役に立つことができれば幸いです！