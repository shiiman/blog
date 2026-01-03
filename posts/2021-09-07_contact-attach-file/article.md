---
id: 326
title: 【WordPress】Contact Form 7でお問い合わせにファイルを添付させる方法
slug: contact-attach-file
status: publish
excerpt: 前回の記事でWordPressのプラグインであるContact Form 7を使用してお問い合わせフォームを作成しました。こちらの設定までできていれば基本的なお問い合わせは問題ありません。 まだFormの設定ができていな \[…\]
categories:
    - 4
    - 10
    - 3
tags:
    - 6
    - 23
    - 28
    - 34
featured_media: 327
date: 2021-09-07T19:30:00
modified: 2021-09-08T01:37:08
---

前回の記事でWordPressのプラグインであるContact Form 7を使用してお問い合わせフォームを作成しました。こちらの設定までできていれば基本的なお問い合わせは問題ありません。

まだFormの設定ができていない方はこちらも合わせてご確認ください。

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
2021.09.06](https://shiimanblog.com/wordpress/contact "【WordPress】Contact Form 7でお問い合わせフォームを作成しよう / 初心者でも簡単！ビジネスや広告をやる方は必須")

しかし、バグ報告で画像を添付したり、データを送るためファイルを添付したいことがあると思います。今回はそんな時に、お問い合わせフォームに画像を添付させられる機能を追加でつけていきたいと思います。

## フォームにファイル添付追加

Contact Form 7でフォームを作成している場合、ファイルの追加はとても簡単です。

お問い合わせ -\> コンタクトフォームのフォームタブで「ファイル」を選択します。

[![Contact File](https://shiimanblog.com/wp-content/uploads/2021/09/contact-file-800x382.png)](https://shiimanblog.com/wp-content/uploads/2021/09/contact-file.png)

するとファイル添付の設定ポップアップが開きますので、必要な項目を記入します。

今回はファイルサイズを 2MBに、受け入れ可能なファイル形式を png,jpg,pdf に設定しました。

デフォルト値はそれぞれ下記です。

ファイルサイズの上限(バイト) : 1MB

受け入れ可能なファイル形式 : jpg、jpeg、png、gif、pdf、doc、docx、ppt、pptx、odt、avi、ogg、m4a、mov、mp3、mp4、mpg、wav, wmv

デフォルト値が設定されているので、特に何も入力しなくても動作します。

調整したい場合のみ値を入力しましょう。

入力できたら「タグを挿入」で設定を追加します。

[![ファイル添付セッテイング](https://shiimanblog.com/wp-content/uploads/2021/09/contact-file-setting-800x680.png)](https://shiimanblog.com/wp-content/uploads/2021/09/contact-file-setting.png)

するとフォームにタグが設定されます。このままだとボタンだけ表示されてしまいますので、他の設定と合わせて label タグで調整します。

こちらが編集できたら保存しておきましょう。

- ![](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-add-file-1-800x995.png)
- ![](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-add-file2-1-800x891.png)

次に受け取りメールの設定をしていきます。

受け取りメールの設定を忘れがちなので注意しましょう！！

webサイトのお問い合わせフォームの表示は上記の設定でファイルが添付できるようになりましたが、このままだと折角添付されたファイルを受け取ることができません。

メールタブで添付ファイルを受け取れるように設定しましょう。

添付ファイルのメールタグが上部に出ていますので、こちらを「ファイル添付」欄にコピーします。

こちらを設定するだけで、添付されているファイルをメールで受け取ることが可能になります。「メール(2)」の方のファイル添付にも同じメールタグを入力してあげましょう。そうすることで自動返信メールにも添付ファイルが送られるので、送信者がちゃんとファイル送信できたのかどうかが分かるようになり、より親切です。

[![コンタクトフォーム - メールタブ](https://shiimanblog.com/wp-content/uploads/2021/09/add-file-mail-800x459.png)](https://shiimanblog.com/wp-content/uploads/2021/09/add-file-mail.png)[![コンタクトフォーム - メールタグ(ファイル添付)](https://shiimanblog.com/wp-content/uploads/2021/09/add-file-mail-setting-800x300.png)](https://shiimanblog.com/wp-content/uploads/2021/09/add-file-mail-setting.png)

こちらで設定して完成です。

画面上部に表示されているショートコードをコピーして、固定ページに貼り付けてみましょう。

こんな感じでファイル添付ボタンが表示されました。

表示の確認後は実際にメールを送って、ファイルが添付されることを確認しましょう。

[![コンタクトフォーム - ファイル添付確認](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-file-confirm-800x783.png)](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-file-confirm.png)

これで添付ファイルの設定はおしまいです。

## その他のカスタタイズ

Contact Form 7 にはファイル添付のように簡単にフォームをカスタマイズ出来る要素が沢山用意されています。実際にコンタクトフォーム画面でどんなカスタムが可能なのか見てみましょう。

14もの項目があり、それぞれ簡単な設定だけでフォームをカスタマイズすることが可能です。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-custom-800x346.png)](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-custom.png)

こんな感じでラジオボタン(性別で使用されているもの)を配置したり数値入力欄(年齢で使用されているもの)を追加することも簡単にできます。お問い合わせフォームで必要なものを要件に合わせてカスタマイズしてみましょう。

[![コンタクトフォーム - 性別、年齢](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-age-800x334.png)](https://shiimanblog.com/wp-content/uploads/2021/09/contact-form-age.png)

## まとめ

今回はWordPressのプラグインであるContact Form 7 を使用してお問い合わせフォームにファイル添付機能を追加してみました。Contact Form 7 にはデフォルトでカスタムする項目が多く備わっているので特別な技術だったり追加のプラグインを必要としません。

簡単な設定を追加することでファイル添付の設定を追加することができました。

またファイル添付以外にも簡単にお問い合わせフォームをカスタムする項目がありますので、必要に応じてカスタマイズしてみましょう。

最後にお問い合わせフォームはカスタマイズしたら必ず送信確認をしてください。

設定を変更したらどこかに不具合が生じてメールが送られないということはよくあります。

実際にメールが送られてくるのを確認して、設定完了にするように心がけましょう。