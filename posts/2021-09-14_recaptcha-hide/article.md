---
id: 171
title: 【WordPress】Invisible reCaptchaの画像を非表示にする方法
slug: recaptcha-hide
status: publish
excerpt: WordPressのコメントやお知らせフォームからBotによる攻撃を防ぐためのプラグインとして「Invisible reCaptcha」というプラグインを導入している方は多いと思います。 こちらを導入して設定する方法は別 \[…\]
categories:
    - 3
tags:
    - 38
    - 42
featured_media: 583
date: 2021-09-14T19:30:00
modified: 2021-09-12T23:59:37
---

WordPressのコメントやお知らせフォームからBotによる攻撃を防ぐためのプラグインとして「Invisible reCaptcha」というプラグインを導入している方は多いと思います。

こちらを導入して設定する方法は別の記事で紹介しておりますので、まだ設定していない方はこちらを参考に設定してみてください。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/invisible_recaptcha-320x180.png)\
\
【Wordpress】おすすめのプラグイン Part 4 〜 Invisible reCaptcha 〜\
\
2021年8月にConoHa WingでWordpressを使用してブログを始めたので、Wordpressの初期設定として入れたプラグインを紹介します！\
今回はInvisible reCaptchaというもので、wordpressのセキュリティーアップを行うプラグインです。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.09.10](https://shiimanblog.com/wordpress/invisible_recaptcha/ "【Wordpress】おすすめのプラグイン Part 4 〜 Invisible reCaptcha 〜")

このプラグインを設定すると以下の様な画像がフォーム、もしくは画面下に表示されます。

こちらの画像が場合によってはデザインの邪魔であったり、UI的にいけていなかったりしてできるなら **非表示にしたい** ところです。

[![reCAPTCHA](https://shiimanblog.com/wp-content/uploads/2021/09/recaptcha.png)](https://shiimanblog.com/wp-content/uploads/2021/09/recaptcha.png)

しかし単純に非表示にすることはGoogleの規約上できません。

今回はGoogleの公式の方法を使用して画像の非表示をおこなっていきます！

もし同じように画像がちょっと邪魔だなと感じた方がいらっしゃれば、こちらの記事を参考に画像を非表示にしてみてください。

## Googleの規約確認

まずはGoogleの公式ページQ&Aページでやり方を確認しましょう。

[>\> Google reCAPTCHA公式ページQ&A](https://developers.google.com/recaptcha/docs/faq)

するとreCAPTCHAバッチを非表示にするには、以下のテキストを表示してください！と記載してあります。

以下のテキストはこちらの部分です。

```
This site is protected by reCAPTCHA and the Google
    <a href="https://policies.google.com/privacy">Privacy Policy</a> and
    <a href="https://policies.google.com/terms">Terms of Service</a> apply.
```

[![Google - reCAPTCHA](https://shiimanblog.com/wp-content/uploads/2021/09/google_recaptcha.png)](https://shiimanblog.com/wp-content/uploads/2021/09/google_recaptcha.png)

つまり、GoogleのreCAPTCHAを使ってサイトを守っているよ！という文言をフォームから見えるところに記載することによって、画像を非表示にしていいよ！ということです。

## 非表示設定

Googleの規約が確認できたということで、まずは文言を追加していきましょう！

### フッターに文言を追加する

お問い合わせフォームやコメントの下部に先程の定型文を埋め込んでも良いのですが、フォームを追加するたびにいちいち定型文を埋め込むのも大変なので、私はフッターメニューに常に定型文を含めることにしました。

フッターに含めるとこんな感じになります。

[![フッター - reCaptcha](https://shiimanblog.com/wp-content/uploads/2021/09/footer_recaptcha-800x60.png)](https://shiimanblog.com/wp-content/uploads/2021/09/footer_recaptcha.png)

こちらのように表示する設定をしていきます。

フッターへのHTMLの追加はCocoonというテーマを使用していれば簡単です。

「Cocoon 設定」->「フッター」タブのクレジット表記を変更していきます。

「独自表記」を選択し、下記のHTMLをコピー&ペーストして保存します。

```
This site is protected by reCAPTCHA and the Google <a href="https://policies.google.com/privacy">Privacy Policy</a> and <a href="https://policies.google.com/terms">Terms of Service</a> apply.<br>
Copyright © 2021 [サイト名] All Rights Reserved.
```

※ \[サイト名\]の部分を**ご自身のサイト名に変更**してご利用ください。

[![フッター - 設定](https://shiimanblog.com/wp-content/uploads/2021/09/footer_setting-800x396.png)](https://shiimanblog.com/wp-content/uploads/2021/09/footer_setting.png)

こちらの設定でフッターに定型文の表記ができました！

### reCaptcha画像を非表示にする

表記ができたらreCaptcha画像を非表示にしていきます。

WordPressの左メニューから「外観」-> 「カスタマイズ」を選択します。

すると下記のようなカスタマイズ画面に遷移しますので「追加 CSS」を選択します。

[![reCaptcha - 非表示設定1](https://shiimanblog.com/wp-content/uploads/2021/09/hide_css-800x811.png)](https://shiimanblog.com/wp-content/uploads/2021/09/hide_css.png)

GoogleのQ&Aページにも記載してあった下記のCSSコードをこちらに貼り付けます。

そして「公開」ボタンを押して反映しましょう。

```
.grecaptcha-badge { visibility: hidden; }
```

[![reCaptcha - 非表示設定2](https://shiimanblog.com/wp-content/uploads/2021/09/hide_css2-800x289.png)](https://shiimanblog.com/wp-content/uploads/2021/09/hide_css2.png)

これでコメントやお問い合わせフォームに表示されていたreCaptcha画像が表示されなくなったはずです。実際にご自身のページを確認して画像が非表示されていることを確認しましょう。

## まとめ

今回はコメントやお問い合わせフォームからBotの攻撃を防ぐプラグイン「Invisible reCaptcha」の設定で、表示されていたreCaptcha画像を非表示する方法を解説しました。

このreCaptcha画像はテーマによってはデザインが崩れたり、何かとかぶって表示されてしまったりと、困っていた方が多いのではないでしょうか。

私自身も最初テーマとの相性が悪く、表示が崩れておりました。

そこで今回のように画像を非表示にすることで、画面が崩れずスッキリさせることができたと思います。

ただし、Google公式で記載のある定型文を入れずにreCaptcha画像を非表示にするのは絶対にやめてください。しっかりルールを守りサイトを守るプラグインを使用していきましょう！！