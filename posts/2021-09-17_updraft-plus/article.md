---
id: 637
title: 【WordPress】おすすめのプラグイン Part 5 〜 Updraft Plus 〜
slug: updraft-plus
status: publish
date: 2021-09-17T19:30:00
modified: 2021-09-14T00:46:53
excerpt: WordPressのバックアッププラグイン「UpdraftPlus」の導入・設定方法を紹介。自動バックアップで万が一に備えます。
categories: [10, 3]
tags: [6, 44]
featured_media: 638
---

こんばんは、しーまんです！

皆さんはWordPressのバックアップを取っていますでしょうか？

WordPressは比較的簡単に動かなくなったり、動作がおかしくなったりしますので、定期的なバックアップをおすすめします。

その際に以前「BackWPup」というバックアップを行うプラグインを紹介しました。

こちらの記事を読んでいない方は今回の記事と合わせて御覧ください。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/wordpress-1-320x180.png)\
\
【Wordpress】おすすめのプラグイン Part 1 〜 BackWPup 〜\
\
2021年8月にConoHa WingでWordpressを使用してブログを始めたので、Wordpressの初期設定として入れたプラグインを紹介します！\
初回はBackWPupというもので、wordpressのバックアップを行うプラグインです。まず最初にこちらをインストールし設定しましょう。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.09.01](https://shiimanblog.com/wordpress/backwpup/ "【Wordpress】おすすめのプラグイン Part 1 〜 BackWPup 〜")

今回は私が「BackWPup」から「Updraft Plus」というプラグインに変更したので、その理由を紹介していきたいと思います！

## Updraft Plus

私はWordPressを使ってブログを書くようになってから、いろいろなプラグインを試しています。

そんな時に「Updraft Plus」というプラグインを見つけました。

私がなぜ「BackWPup」から「Updraft Plus」に乗り換えたかというと、2点あります。

1. バックアップの保存場所としてGoogle Driveを無料で選択可能
2. バックアップしたファイルから簡単に復元ができる！

特に2番の簡単に復元できるという点は、「BackWPup」ではできない優れた点に感じました！

この2つの理由から私はバックアッププラグインの変更を決めました。

### インストール手順

それでは早速インストールしていきましょう。

※ インストールをする前に競合する「BackWPup」は一応アンインストールをしておきました。

WordPressの「プラグイン」-> 「新規追加」から「Updraft Plus」を検索します。

検索で出てきたら「今すぐインストール」->「有効化」をしていきましょう。

[![Updraft Plus - インストール](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus-800x538.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus.png)

### バックアップ手順

インストールが終わったら早速バックアップをしてみましょう。

WordPressの左メニューもしくは上部に「UpdraftPlus」の設定が追加されています。

そこから「バックアップ/復元」タブを開き「今すぐバックアップ」を押してみましょう。

すると直ぐにバックアップが始まりまる。

これで手動バックアップは完了です。テーマの変更やプラグインの追加前はこのように手動でバックアップを取ってから変更する癖をつけましょう。

[![Updraft Plus - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting1-800x369.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting1.jpg)

### Google Driveに保存する設定

続いてバックアップのスケジュール設定と、保存場所を設定していきましょう。

まずは「設定」タブを開き「ファイルバックアップのスケジュール」と「データベースバックアップのスケジュール」を設定していきましょう。こちらでバックアップ頻度と保存ファイル数を設定可能です。

無料版では細かい時間の設定はできません。

細かい時間指定が必要な場合は有料版をお試しください。

次に「保存先を選択」していきます。

UpdraftPlusの保存先はかなり多い選択しを用意してくれているので、ご自身でバックアップファイルを保存したい場所を選択しましょう。

今回はGoogle Driveを選択していきます。

[![Updraft Plus - 設定5](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting5-800x531.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting5.jpg)

Google Driveを選択すると最初に注意ポップアップが表示されますので、リンクをクリックします。

[![Updraft Plus - 設定6](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting6-800x236.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting6.png)

するとGoogleのログインアカウント選択を求められますので、紐付けるアカウントを選択します。

[![Updraft Plus - 設定7](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting7-1-800x1254.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting7-1.png)

アカウントの選択を行うとUpdraftPlusがGoogleアカウントに対してアクセス許可を求めるポップアップが出ますので、こちらを確認して「許可」をしましょう。

[![Updraft Plus - 設定8](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting8-800x1371.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting8.png)

「許可」が確認されるとオレンジ色の画面が表示されます。

こちらが表示されたらGoogle Driveとの連携は完了です。「Complete setup」を押しましょう。

[![Updraft Plus - 設定9](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting9-1-800x321.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting9-1.png)

### Google Driveへの保存確認

Google Drive連携の設定ができたら、実際にバックアップができるか確認しましょう。

「バックアップ/復元」タブにもどり「今すぐバックアップ」を再度実行します。

バックアップが完了するとご自身のGoogle Driveに「UpdraftPlus」ディレクトリができているはずです。

[![Updraft Plus - 設定10](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting10-800x406.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting10.png)

中身を確認してみましょう。こちらがバックアップの中身になります。

[![Updraft Plus - 設定11](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting11-800x437.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting11.png)

ここまででバックアップの設定は完了です。

### 復元手順

バックアップの設定が完了したら、復元方法も確認しておきましょう。

Pudraft Plusは復元も画面上から出来るのでとても簡単です。

まずは「バックアップ/復元」タブを開きましょう。すると画面の下部に「既存のバックアップ」という箇所が見つかるはずです。こちらにバックアップされた履歴が表示されます。

その横に「復元」ボタンが存在します。対象のバックアップを選択して「復元」ボタンを押しましょう。

[![Updraft Plus - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting2-800x415.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting2.jpg)

次に復元するコンポーネントを選択します。

分からない場合は全て選択すると、バックアップ時点の状態に全て戻ります。

[![Updraft Plus - 設定3](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting3-800x324.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting3.png)

コンポーネントを選択して「次」を選択すると確認画面が表示されますので、問題なければ「復元」を押しましょう。するとバックアップファイルからWordPressの復元が行われます。

[![Updraft Plus - 設定4](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting4-800x257.png)](https://shiimanblog.com/wp-content/uploads/2021/09/updraftplus_setting4.png)

以上で復元は終わりです。

本当に簡単ですね。

## まとめ

今回は私がWordPressのバックアッププラグインを「BackWPup」から「Updraft Plus」に乗り換えた理由を説明した後、「Updraft Plus」のインストール・設定方法を紹介しました。

また、バックアップファイルからの復元方法も合わせて紹介しました。

設定も、復元もとても簡単に行えたと思います。

バックアップはこのように手軽に行うことが可能ですので、WordPressの設定変更時はぜひバックアップをしてから行うようにしましょう！！