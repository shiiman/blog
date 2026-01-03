---
id: 82
title: 【WordPress】おすすめのプラグイン Part 1 〜 BackWPup 〜
slug: backwpup
status: publish
excerpt: 前回の記事ででWordPressを導入するところまで行いました。ブログ開設までの流れを解説しておりますので、そちらも合わせて御覧ください。 さてブログはとりあえず公開することができましたが、そのままの設定だと実際に運営し \[…\]
categories:
    - 10
    - 3
tags:
    - 6
    - 23
    - 24
featured_media: 84
date: 2021-09-01T19:30:00
modified: 2021-09-24T21:40:59
---

前回の記事で [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でWordPressを導入するところまで行いました。

ブログ開設までの流れを解説しておりますので、そちらも合わせて御覧ください。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/cw-320x180.png)\
\
【2021年8月】初心者がConoHa WINGでWordpressを始めたので手順を1からまとめてみた\
\
2021年8月にConoHa WingでWordpressを使用してブログを始めたので、その時の画像を使って手順を1からまとめてみました。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.08.31](https://shiimanblog.com/wordpress/conoha-wing/ "【2021年8月】初心者がConoHa WINGでWordpressを始めたので手順を1からまとめてみた")

さてブログはとりあえず公開することができましたが、

そのままの設定だと実際に運営していくにはまだ色々足りません。

私がブログ開設にあたり少しづつ設定を増やしておりますので、

そちらももらさず紹介できればなと思っています。

今回はその中の1つ!! WordPressのプラグインについて紹介するシリーズになります。

## BackWPupとは

ブログを立ち上げてさぁドンドン記事を書いていくぞ〜

と思ったあなた。少しお待ち下さい！

記事を書いていく前にかならずやっておきたい設定が2つほどあります。

それが以下の2つ！

- バックアップ設定
- セキュリティー対策

なげこの2つが重要かというと

- サーバはいつ壊れるか分からない。
- サーバはいつ悪意のあるユーザに狙われるか分からない。

からです。。。

(SREという職業柄この辺のことがいの一番に心配してしまうのですよね。。。)

ということで今回はいつ壊れても復元出来るように

WordPressのバックアップの設定を出来る **BackWPup** というプラグインを紹介します。

実はConoHa WINGでは1日1回 14日分のバックアップを自動で行ってくれています。

なのでこれで十分だという方はここからの設定は不要かもしれません。

ただしConoHa WINGの自動バックアップでは以下のような問題点があります。

- いつの時点のバックアップか分からない
- 好きなタイミングでバックアップを行えない
- 14日以前の状態には戻せない

よって私は追加で **BackWPup** を設定しています。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/Screenshot-2-1024x417.png)](https://shiimanblog.com/wp-content/uploads/2021/08/Screenshot-2.png)

## BackWPupの設定手順

それでは早速プラグインの導入に参りましょう。

WordPressの画面から設定を行っていきます。

### プラグインの検索

画面左のメニューから「プラグイン – 新規追加」を選択。

表示画面の右の検索ウィンドウで「BackWPup」と入力します。

すると関連の強いプラグインが表示されます。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/1-1-1024x545.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/1-1.jpg)

### インストール

対象のブラグインが見つかったら「今すぐインストール」ボタンを押してインストールします。

プラグインはインストールを押しただけでは有効になりません。

インストールが終わると「今すぐインストール」ボタンが「有効化」ボタンに変わりますので、そちらをクリックすると有効になります。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/2-1.png)](https://shiimanblog.com/wp-content/uploads/2021/08/2-1.png)

### ジョブの設定

プラグインが正常にインストールされ、有効化状態になると左のメニューに「BackWPup」が追加されます。この状態だとまだ何も設定されておりませんので、実際にバックアップを行うジョブを設定していきます。

左メニューから「BackWPup -> ジョブ」を選択。

選択画面の上部に新規追加ボタンがありますので、こちらをクリックします。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/3-1024x545.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/3.jpg)

するとバックアップジョブの設定が表示されます。

基本はデフォルトのままでも良いですが、「一般」タブで設定する箇所は3箇所です。

1. このジョブの名前: 分かりやすい名前をつけておきましょう。
2. アーカイブ形式: Tar GZipを選択
3. バックアップファイルの保存方法: 今回はとりあえず フォルダーへバックアップを選択しましょう。

設定できましたら、忘れずに「設定を保存」ボタンを押しましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/4-998x1024.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/4.jpg)

次に「スケジュール」タブを開きます。

そしてジョブの開始方法で「WordPress の cron」を選択しましょう。

するといつバックアップを取るかをスケジューラーで設定できます。

今回は毎週日曜日 AM3時にバックアップを取るように設定していきます。

設定できたら「変更を保存」ボタンを押しましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/5-1024x637.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/5.jpg)

ここまで設定できたらバックアップ設定の完了です。

日曜日のAM3時時点のバックアップが毎週取られます。

### 手動バックアップ

では実際にバックアップを取ってみましょう。

先程作ったジョブを選択すると「今すぐ実行」というテキストリンクが表示されます。

こちらを押すと指定した時間以外にも、すぐにバックアップを取ることが可能です。

プラグインの導入やアップデート前は必ずバックアップを取る癖をつけましょう。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/6-1024x504.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/6.jpg)

実際にバックアップされたファイルはConoHa ファイルマネージャーから確認することができます。

```
[ドメイン]/wp-content/uploads/backwpup-[文字列]-backups
```

[![](https://shiimanblog.com/wp-content/uploads/2021/08/Screenshot-3-1024x658.png)](https://shiimanblog.com/wp-content/uploads/2021/08/Screenshot-3.png)

保存されるディレクトリを変更したい場合は設定から変更可能です。

[![](https://shiimanblog.com/wp-content/uploads/2021/08/Screenshot-1-2-1024x541.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/Screenshot-1-2.jpg)

## まとめ

今回はBackWPupプラグインを使用してWordPressのバックアップ方法を解説しました。

WordPressはプラグインのアップデートや、競合、予期せぬスクリプト変更などで結構直ぐに動かなくなります。その為そのような変更を行う場合はバックアップを取ってから行う癖をつけておきましょう。

バックアップを取ることで、今まで書いてきた記事が全て消えてしまったとか、サーバが壊れて修復できなくなったなどの事態を防ぐことができます。

ぜひWordPressを導入したら最初にバックアップ設定をしてみてください。