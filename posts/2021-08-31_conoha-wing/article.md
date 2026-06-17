---
title: 【2021年8月】初心者がConoHa WINGでWordPressを始めたので手順を1からまとめてみた
slug: conoha-wing
date: 2021-08-31T19:30:00.000Z
categories:
  - initialization
  - wordpress
tags:
  - wordpress
  - conoha-wing
draft: false
id: 21
modified: 2021-09-24T21:39:54.000Z
excerpt: WordPress初心者向けに、ConoHa WINGでのブログ開設手順をサーバー契約からWordPressインストールまで画像つきで丁寧に解説します。
eyecatch: ./assets/eyecatch.png
---

ブログを書こうと思ってからじゃあwordpressを使おうと思って色々調べて見ました。

最初は自前でサーバを立てようかなと思ったのですが、 [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) で管理するのが圧倒的に楽そうだったので [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) で始めることにしました。

 ![](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SNCY9) 

## ConoHa WINGとは

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) とは初心者がWordPressでブログやホームページを作る際によく使われるレンタルサーバです。

ここではその特徴を3つ紹介します。

### 特徴1: 簡単

まずは、初期設定が驚くほど簡単で初心者に優しいということです。

普通の方はwordpressの導入 / wordpressのテーマ / ドメイン / データベース / メール などなどサーバを開始するにあたり行う作業は結構あるのですが、その辺の作業を全て [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) の方で担ってくれるので、驚くほど簡単にブログやホームページが開設できます。

![](./assets/Screenshot-1-1.png)

### 特徴2: 安い

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) は初期費用がかからず、月額1000円を切る料金体系になっています。

まずはお試しで3ヶ月とかで試して、続けられそうなら12ヶ月以上の契約をするのがおすすめです。

![](./assets/Screenshot.png)

### 特徴3: 表示速度国内No.1

webサーバの処理速度も他社を圧倒すると謳っています。

これはおそらく [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) はコンテンツキャッシュとかも標準で備えているので、

そのあたりが関係しているのではないかと思いますが、それでも早いに越したことはないですね。

![](./assets/Screenshot-1.png)

## ConoHa WingでWordpressサイト作成手順

ではここから、実際に [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) に登録してwordpressを立ち上げ、

webページにアクセスするまでの手順をみていきましょう。

### 申し込み

下記ボタンから [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) の申し込みページに飛ぶことができます。

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM)

ページを開いたら「今すぐお申し込み」をクリックします。

![](./assets/18d76db387d73133ddcce97cd865dc7c-2.jpg)

これから新しく申し込みを行うので「初めてご利用の方」の欄を入力して次へ

![](./assets/52acc05ae390f13efec500155675a0d3.jpg)

### プラン選択

まずはプランの選択です。

- 料金タイプ: WINGパック
- 契約期間はお好みで選択してください。
- プラン: 最初はベーシックで十分です。
- 初期ドメイン: こちらは実際には使用しないので、何でも構いません。
- サーバ名: こちらはwordpressを導入するサーバの名前です。特に気にならなければデフォルトで問題ありません。
- WordPressかんたんセットアップ: 利用する

![](./assets/11.jpg)

次にドメインを設定します。

ドメインとはブラウザのバーに表示されるweb上の住所みたいなものです。

こちらは一度設定すると変更ができませんので、最初にしっかり考えてから入力しましょう。

- 独自ドメイン設定: あなたのサイトのurlになります。
- 作成サイト名: こちらは後ほど変更が可能です。
- WordPressユーザ名: WordPressにログインする際に使用します。
- WordPressパスワード: WordPressにログインする際に使用します。
- WordPressテーマ: こだわりがなければCocoonを選択しましょう。

![](./assets/2.jpg)

独自ドメインに希望のドメインを入れた際に取得できない場合があります。

原因は他の人が既にそのドメインを取得している場合が考えられます。

ドメインとは住所のようなものを説明しましたが、世界で1つである必要があります。

その場合はドメインを変更するか .com 部分を別のものに変えて設定してみてください。

### お客様情報入力

一通りプランの入力が終わったら、次にお客様情報の入力です。

ご自身の情報を入力していきましょう。

![](./assets/1.png)

個人情報の入力後、本人確認が入りますのでSMS認証か電話認証で確認しましょう。

![](./assets/2.png)

### お支払い設定

最後にお支払い情報の入力です。

内容が正しいか確認しお申し込みしましょう。

![](./assets/12657264-19dd991a7dd7a973854daec897539a6a-1.jpg)

## Webサイトへのアクセス

### Webサイトへアクセスしてみよう

申込みが完了すると [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) の管理画面にログインすることができます。

サイト管理 -> サーバ設定 -\> サイトURLにアクセスすると作成された初期webサイトが表示されます。

![](./assets/12659675-ba0e6df566ec04d496eb332322900971.png) こちらが初期のwebサイトの画面になります。

この時点ではssl設定が完了していませんので http:// でのアクセスになります。

### sslアクセスを確認しよう。

プラン選択の際に「WordPressかんたんセットアップ」を利用するを選択していた場合、

ssl化の設定まで自動で行ってくれます。

ただその進捗が見えないので、

いつになったらhttps://でアクセス可能になるのかが分かりませんでした。

私はここで自分でssl化対応をしなくてはと思い、色々いじっていたらサイトにアクセスできなくなりとても焦りました。

ssl対応が終わったかどうかは [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) の簡単SSL化のボタンがグレーアウトから水色のボタンに変わったら、サイトURLもhttp://からhttps://に変わっていました。

私の場合ssl化が終わるまで1~2時間くらいかかった印象です。もしかしたらもっと早く終わるかもしれませんし、もっと時間がかかるかもしれません。

ですのでwebサイトのアクセスまで確認できたら、ssl化が自動で終わるまでは逸る気持ちを抑えて他の作業でもして待っていましょう。

![](./assets/12659322-3280b18c21ebda346b5db994733633c3-1.png)

## まとめ

今回はConoHa WINGの登録からWordpressでwebサイトの表示までを行いました。

いかがだったでしょうか。

思っていたよりも、登録からwebサイトの表示まであっという間と感じたのではないでしょうか。

これくらいならできそうと少しでも思っていただけたら幸いです。

この後実際にブログを始めるにはWordPressの各設定をしなくてはなりません。

そちらに関してはまた別の記事でアップしていきたいと思います。
