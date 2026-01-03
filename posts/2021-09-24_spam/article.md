---
id: 175
title: 【WordPress】お問い合わせが迷惑メールに入ってしまう件を解決 | Contact Form 7 | ConoHa WING
slug: spam
status: publish
excerpt: WordPressでお問い合わせフォームを作成する場合「Contact Form 7」を使用している方が多いと思います。Contact Form 7の設定方法は下記の記事を参考にしてください。 【WordPress】Co \[…\]
categories:
    - 4
    - 3
tags:
    - 6
    - 7
    - 51
featured_media: 805
date: 2021-09-24T19:30:00
modified: 2021-09-24T21:51:30
---

WordPressでお問い合わせフォームを作成する場合「Contact Form 7」を使用している方が多いと思います。Contact Form 7の設定方法は下記の記事を参考にしてください。

 [![](https://shiimanblog.com/wp-content/uploads/2021/09/eyecatch8-320x180.png)\
\
【WordPress】Contact Form 7でお問い合わせフォームを作成しよう / 初心者でも簡単！ビジネスや広告をやる方は必須\
\
WordPressのプラグインContact Form 7をしようして、お問合せフォームを作成します。\
初心者でも簡単にできますので、ビジネスや広告でサイトを使用している方はぜひご覧ください。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.09.06](https://shiimanblog.com/wordpress/contact/ "【WordPress】Contact Form 7でお問い合わせフォームを作成しよう / 初心者でも簡単！ビジネスや広告をやる方は必須")

今回は Contact Form 7の自動返信メールが迷惑メールになってしまうことがあり、その解決方法を紹介したいと思います。

## 迷惑メールになる原因

メールが迷惑メールになる原因は正確に説明するとかなり難しいです。

すごく簡単に説明すると対象のメールが「なりすまし判定」されてしまっていることが原因です。

ではなりすまし判定されてしまう要因は何でしょうか。

それは**フォームを設置しているサーバ** と **メールを送信するサーバ** が異なることが原因です。

その判定はいくつかありますが、一番分かりやすいのはメールを送信するサーバのドメインと返信を行うメールアドレスが異なる場合です。

具体的には以下のようなケースが考えられます。

サイトドメイン : https:// **shiimanblog.com**

送信元メールアドレス : info@ **gmail.com**

上記の太字青線部分が異なる場合がなりすましと判定されてしまうパターンになります。

上記のケース以外にもなりすまし判定されてしまうケースはありますが、「 **SPF**」とか「 **DKIM**」とか「 **DMARC**」とか言われるメールの知識が必要になりますので、今回は省きます。もっと詳しくこの辺のことを知りたい方は「 **SPF**」「 **DKIM**」「 **DMARC**」をキーワードとして検索してみてください。

## 解決方法

迷惑メール判定されてしまう原因は **フォームを設置しているサーバ** と **メールを送信するサーバ** が異なることが原因でした。ではどうすればその判定を回避できるのでしょうか。

一つの解決策としてSMTPサーバを通してメールを配信することです。SMTPとはSimple Mail Transfer Protocolのことで、それを使ったサーバはメールの送信元と宛先の間で、メールを送受信または中継することを主な目的としたアプリケーションを指します。

なぜ中継するサーバを入れるかというとSMTPサーバの設定を調整することでドメインが違っていてもなりすましではないですよと、承認処理を間に設定することが可能になるからです。

そこで今回はConoHa WingのSMTPと「WP Mail SMTP by WPForms」というプラグインを使用して解決する手順を紹介していきます。

## ConoHa WINGメール設定

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でWordPressを構築している場合はメールサーバも同時に用意してくれています。

こちらで新しく問い合わせ用のメールを作成して設定をしていきましょう。

他のレンタルサーバでもメールサーバを用意してくれていると思いますので、各レンタルサーバのマニュアルをご確認ください。

まずは新しいメールアドレスを取得します。

ConoHa WINGのコントロールパネルから「メール管理」-> 「メール設定」->「メールアドレス」タブを選択しましょう。すると右側に「 **\+ メールアドレス**」というボタンがありますので、こちらをクリックします。

[![ConoHa Wing - メール設定1](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail1-800x277.png)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail1.png)

新しいメールアドレス作成の画面になりますので、こちらでお好みのメールアドレスとパスワードを設定しましょう。

よくあるメールアドレスだと下記のような命名のものを見かけますね。ただこちらはお好みですのでご自身で自動返信にふさわしいメールアドレスを設定しましょう。

- info@
- support@

入力できたら保存をします。

[![ConoHa Wing - メール設定2](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail2-800x358.png)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail2.png)

すると新しいメールアドレスの取得が完了です。

[![ConoHa Wing - メール設定3](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail3-800x266.png)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail3.png)

次になりすましの判定をされないように設定をしてきます。

まず作成したメールアドレスをクリックしてメールアドレス詳細を開きます。

下の方に「DNS情報」と書かれた部分が存在します。

こちらをこの後の設定で使用しますので、メモしておきましょう。

[![ConoHa Wing - メール設定4](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail4-1-800x487.png)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail4-1.png)

メモをしたら左のメニューから「DNS」を選択し、ドメインリストからご自身のwebサイトのドメインを選択します。その画面の右側の鉛筆マークを選択すると編集出来るモードに変わります。

左下のプラスマークを押すと行を追加できますので、先程メモしたものとプラスで下記1行を追加していきます。追加する順番は関係ありませんので、どれから追加しても構いません。またTTLの部分は全て3600に設定しましょう。

```
_dmarc 3600 v=DMARC1; p=none; fo=1; rua=mailto:[作成したメールアドレス]
```

**全部で3行** 追加できたら保存ボタンをクリックしましょう。

これで [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)側のメール設定は完了です。

[![ConoHa Wing - メール設定5](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail5-2-800x444.png)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail5-2.png)

## SMTP設定

続いてWordPress側の設定をしてきます。まずは「WP Mail SMTP by WPForms」というプラグインをインストールしてきましょう。

### WP Mail SMTP by WPFormsインストール

WordPressの左メニューからプラグイン -> 新規追加を選択し「WP Mail SMTP by WPForms」を検索しましょう。そこで検索されたプラグインを今すぐインストール -> 有効化します。

[![WP Mail SMTG by WPForms](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms-800x426.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms.png)

### WP Mail SMTP by WPForms設定

プラグインの有効化ができたら設定を行っていきます。

[![WP Mail Smtp by_ WPForms - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms1-1-800x680.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms1-1.png)

SMTPメーラーを選択画面が表示されますので、「その他のSMTP」を選択します。

[![WP Mail Smtp by_ WPForms - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms2-800x518.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms2.png)

メーラーの設定を行っていきます。入力する内容は下記です。

- SMTP ホスト : ConoHa Wingのコントロールパネルから確認します。( **下図を参考**)
- 暗号化 : なし
- SMTP ポート : 587
- 認証 : あり
- SMTPユーザ名 : 登録したメールアドレス
- SMTPパスワード : メールアドレス登録時に設定したパスワード
- フォーム名 : ご自身のサイト名
- 送信元メールアドレス : 登録したメールアドレス

入力が終わりましたら保存して続行をクリックしましょう。

[![WP Mail Smtp by_ WPForms - 設定3](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms3-800x968.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms3.jpg)

**SMTPホスト** はConoHa WINGのコントロールパネルから確認ができます。

「メール管理」->「メール設定」-\> 「メールアドレス」タブから作成したメールアドレスをクリック。

メールサーバーの「SMTPサーバー」の情報を入力します！

[![WP Mail Smtp by_ WPForms - 設定4](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail6-800x343.png)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_mail6.png)

その後の設定はデフォルトで問題ありません。「保存して続行」をクリックします。

[![WP Mail Smtp by_ WPForms - 設定5](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms4-800x614.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms4.jpg)

メール通知の設定が表示されますがこちらは登録する必要はありませんので、Skipを選択しましょう。

[![WP Mail Smtp by_ WPForms - 設定6](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms5-800x769.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms5.png)

ライセンスの登録も求められますが、こちらもスキップで問題ありません。

[![WP Mail Smtp by_ WPForms - 設定7](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms6-800x807.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms6.png)

以上で、プラグインの設定は終了です。

[![WP Mail Smtp by_ WPForms - 設定8](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms7-800x633.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wp_mail_smtp_by_wpforms7.png)

### Contact Form 7の設定

プラグインの設定が完了したらContact Form 7の送信元メールアドレスを変更しましょう。

お問い合わせ -\> コンタクトフォーム の対象フォームを選択し送信元メールアドレスを今回登録したメールアドレスに変更します。自動返信用フォームを設定している場合は「メール」「メール（２）」の2箇所の設定があると思いますので、両方忘れずに設定しましょう。

[![Contact Form 7 - 送信元変更](https://shiimanblog.com/wp-content/uploads/2021/09/mail_setting-800x471.png)](https://shiimanblog.com/wp-content/uploads/2021/09/mail_setting.png)

送信元メールアドレスを変更し、保存をしたら今回の設定は全て終了です。

実際にフォームからお問い合わせを行いメールが届くことを必ず確認しましょう。

今回設定したメールアドレスから自動返信メールが届くのが確認できるはずです。

## まとめ

今回はお問い合わせフォームの自動返信メールが迷惑メールになってしまう問題の解決として、「ConoHa WINGでのメール設定」と「WP Mail SMTP by WPForms」プラグインの設定を解説しました。

なりすましとして迷惑メールになってしまうと大切なメールに気づかないという事態になってしまいます。これではお問い合わせフォームを設置した意味も薄れます。

きちんとメールが届くようにして、安心してお問い合わせフォームを活用できるように設定しておきましょう。今回の記事が、お問い合わせの自動返信メールが迷惑メールに入ってしまう事象が起こっている方の解決に少しでもお役立ちできれば幸いです！