---
id: 4
title: Cocoon設定 Forbidden Accessになる理由は〇〇だった!!
slug: forbidden-error
status: publish
date: 2021-09-12T19:30:00
modified: 2021-09-24T21:49:11
excerpt: CocoonのWordPress設定保存時に発生する「Forbidden Access」エラーの原因と解決方法を紹介します。
categories: [11, 3]
tags: [6, 7, 31]
featured_media: 555
---

私はWordPressのテーマ設定でCocoonを使用しています。

ブログを始めたばかりの初心者である私はいろいろな記事をみてCocooの設定を試していました。そんな時に、設定をして「変更を保存」とすると

**閲覧できません(Forbidden access)**

[![Forbidden access](https://shiimanblog.com/wp-content/uploads/2021/09/Forbidden_access-800x318.png)](https://shiimanblog.com/wp-content/uploads/2021/09/Forbidden_access.png)

と表示されてしまうことがあります。

なんだこれは。。。

Forbidden accessというのは権限が無くて設定の変更やアクセスができませんという意味で、エンジニアをやっているとよく目にするのですが、ブログを始めたばかりの方がこのページに出くわしてしまったらとても焦ります。

あれ、壊れた！なにか悪いことしてしまったのかな？！

ブログ怖い！！

そう思ってしまっても仕方ありませんよね。

今回は実際にForbidden Accessに出くわしてしまった私が、解決に至ったまでの対応策を紹介します。

同じようにForbidden accessページが出てしまって困っている方がいれば参考にしてみてください。

## 結論

最初に結論から述べておきます。

私の場合は [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)というレンタルサーバーを利用してWordPressを作成しています。

そのConoHa WINGのセキュリティーシステムの1つである **WAF**が原因でした。

WAFとはweb application firewallの略です。

みなさんもファイアーウォールという言葉はなんとなく聞いたことがあるんではないでしょうか。

WAFとはそのファイアーウォールの1種で、怪しいアクセスや通信から自身のサーバーを守る役割をしています。しかしそのWAFが今回悪さをし、Cocoonの設定が悪質な通信だと判断してしまい、Fobidden accsssとしてしまっていました。

もしレンタルサーバーが [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)じゃないよって方もWAF(もしくはセキュリティー設定)を一度疑ってみてください。

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でなくてもWAFは存在します。例えばConHa WINGとよく比較されるエックスサーバというレンタルサーバでもWAFの設定があります。その設定を見直すことで解決する場合もあるかもしれません！

## WAFの確認

では実際に [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でWAFの設定を変更してFobidden accessの対応をしていきましょう。

まずはConoHa WINGの [コントロールパネル](https://manage.conoha.jp/) を開きましょう。

ページにログインできたら「サイト管理」->「サイトセキュリティ」->「WAF」と進んでいきましょう。すると攻撃ターゲットURLが表示されているのが確認できると思います。

こちらがForbidden accessとなったアクセスです。

cocoonの設定ではこのWAFに引っかかってしまうようなアクセスが存在するのです。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_waf-800x489.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_waf.jpg)

## 解決方法

解決方法を主に2つあります。

ConHa WINGのコントロールパネルから WAFの利用設定をOFFにすることができるので、OFFにすると解決はできますが、セキュリティー的に弱い設定になってしまいますのでOFFにするのはおすすめしません。

今回はそれ以外の解決方法を紹介します。

### WAFの除外設定を行う

先程のConHa WINGのコントロールパネルのWAFの画面をもう一度確認してみましょう。

「攻撃ターゲットURL」の左横に **「除外」** というボタンが表示されています。

こちらを押すと、対象になったURLのルールだけ除外することができます。

WAFの設定をOFFにせずに今回Forbidden accessになったものだけ除外にできるんですね。

こちらで除外することにより、今までForbidden accessになっていた設定も、問題なく完了させることができます！

### Conoha WINGプラグインから一時的にWAFをOFFにする

もう一つの方法が、WordPressのプラグインである「 **ConoHa WING コントロールパネルプラグイン**」を使用する方法です。

[ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)でWordPressを構築した方は自動で入っているプラグインになります。

もしインストールされていない場合は「新規追加」から検索してみてください。

ConoHa WING コントロールパネルプラグインを有効化するとWordPressの左メニューに「ConoHa WING」という項目が表示されます。そこから設定に進みましょう。

設定ページでは「セキュリティー設定」という項目があり、こちらのチェックボックスからWAFの設定をON, OFFさせることができます。

[![ConoHa Wing プラグイン](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_waf2-800x766.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/conoha_waf2.jpg)

ですのでわざわざConHa WINGのコントールパネルのページを開く必要はないということですね。

この設定を使用して、Cocoonの設定をする時のみWAFの設定をOFFにしましょう。

そしてCocoonの設定が完了したら、再度 [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)のWAFの設定をONにすることを忘れないようにします。

これで必要な時にWAFの設定をON, OFFすることでForbidden accessを回避できましたね。

## まとめ

今回はCocoonの設定時にFobidden accessと表示されてしまう原因について説明しました。

多くの場合はセキュリティー設定(WAF)が原因でFobidden accessと表示されてしまうのですね。

また [ConoHa WING](https://px.a8.net/svt/ejp?a8mat=3HKUB1+DBVECY+50+5SI7RM) ![](https://www19.a8.net/0.gif?a8mat=3HKUB1+DBVECY+50+5SI7RM)においてWAFの設定の確認方法と、実際にどのようにFobidden accessを回避すればよいかの具体的な方法を2つ紹介しました。

同じようなところで躓いてしまった方に少しでも役立つ情報になったら幸いです！