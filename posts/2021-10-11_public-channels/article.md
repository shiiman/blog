---
id: 1154
title: 【slack】GASを使ってslackのpublicチャンネル一覧をスプレッドシートに出力する方法
slug: public-channels
status: publish
excerpt: こんばんは、しーまんです。 チャットツールの定番Slackでチャンネル数が増えすぎてどんなチャンネルが有るのか分からなくなったことはありませんか。また自分の所属しているチャンネルが増えすぎて、どのチャンネルに投稿した方が \[…\]
categories:
    - 18
    - 64
    - 65
tags:
    - 66
    - 67
    - 68
    - 75
featured_media: 1156
date: 2021-10-11T19:30:00
modified: 2022-06-30T18:00:20
---

こんばんは、しーまんです。

チャットツールの定番Slackでチャンネル数が増えすぎてどんなチャンネルが有るのか分からなくなったことはありませんか。また自分の所属しているチャンネルが増えすぎて、どのチャンネルに投稿した方がよいのか迷ってしまったことはないでしょうか。

Slackは手軽にチャンネルを増やすことが出来る一方、手軽すぎるが故に気付いたら不要なチャンネルが大量にできてしまうというデメリットがあります。

ですので、そういった不要なチャンネルを定期的に整理することで、スッキリと業務を行うことができます。

今回はそんな増えすぎたslackのチャンネルを整理したい場合や、どんなpublicチャンネルがあるのか一覧を確認したい場合などに使える、スプレッドシートにpublicチャンネルを一覧化させる方法を紹介します！

## Slack Bot作成

まずはSlack Botを作成してGASと連携させていきます。

ボットの作成方法については以前の記事で解説しておりますので、そちらをご確認ください。

[![](https://shiimanblog.com/wp-content/uploads/2021/10/eyecatch_slackbot-320x180.png)\
\
【2021年】ボットを作成してSlack APIをGASから叩く方法 \| SlackAppライブラリ未使用\
\
Slackでボットの作成とGASでそのボットとの連携方法を紹介します。\
また、GAS側ではSlackAppというライブラリもあるのですが、使用できるAPIが少なすぎるため今回はライブラリなしで実装します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.10.04](https://shiimanblog.com/engineering/gas-for-slack-bot/ "【2021年】ボットを作成してSlack APIをGASから叩く方法 | SlackAppライブラリ未使用")

## スプレッドシートとGASの連携

ボットの作成ができたら次にスプレッドシートとGASを連携させていきます。

まず連携させたいスプレッドシートを開きメニューから「ツール」->「スクリプト エディタ」を選択します。するとGASのエディタが開きますので、こちらでプログラムを入力していきます。

[![GAS - 設定1](https://shiimanblog.com/wp-content/uploads/2021/10/gas1-1.png)](https://shiimanblog.com/wp-content/uploads/2021/10/gas1-1.png)

## 作成するコード

### Slack API

まずは今回必要なSlackのAPIについて洗い出しておきましょう。

今回使用するSlackのAPIは以下の2つです。

1. [チャンネルのリストを取得するAPI](https://api.slack.com/methods/conversations.list)
2. [ユーザのリストを取得するAPI](https://api.slack.com/methods/users.list)

※ それそれのAPIに必要な権限はドキュメントの「Required scopes」に記載されていますので、Botに権限付与することを忘れないでください。

### 関数

次に今回作成する関数をリストアップします。

今回作成する関数は以下5つのメソッドです。

1. 実行メソッド(処理の起点)
2. ユーザの情報を取得するメソッド
3. チャンネルのリストを取得するメソッド
4. チャンネルのリストを取得する際にページャを行うメソッド
5. スプレッドシートの内容を削除するメソッド

今回は実行メソッド内でスプレッドシート内に値を書き込みます。

こちらの処理を分けたい方はもう一つメソッドを追加してもいいかもしれませんね。

2\. 3. のメソッドはSlackのAPIで値を取得しているだけですので、その他のメソッドの作成について紹介していきます。

### スプレッドシートの操作

まずはスプレッドシードの操作方法です。

スプレッドシートを操作する場合、まずはオブジェクトを作成します。

作成方法はいろいろありますが、一番簡単なのは「アクティブなシートのオブジェクトを作成」する方法です。

下記のコードでスプレッドシートのシートオブジェクトを作成することが可能です。

```
let ss = SpreadsheetApp.getActiveSpreadsheet()
let sheet = ss.getActiveSheet()
```

次に操作方法ですが、こちらも簡単です。

スプレッドシートのセルに値を入れる方法は以下です。

```
sheet.getRange(行番号,列番号).setValue(入力値);
```

今回はこちらの簡単な方法のみを使用してスプレッドシートに値を入力しています。

### 注意点

処理を作成する上で今回は3点注意する必要があります。

1. ユーザAPIの取得方法
2. ページャーの操作
3. スプレッドシート更新処理

それぞれ解説していきます。

#### ユーザAPIの取得方法

まずスプレッドシートに出力したい内容を整理します。

今回出力したい内容は以下の4カラムです。

1. チャンネル名
2. トピック
3. 説明
4. 作成者

この内1〜3. の情報はチャンネルリスト取得時に合わせて情報を取得することが可能です。

しかし4\. に関しては **取得可能なのはSlackのユーザID** になります。

こちらをそのまま出力しても、ただのIDですので誰のことを指しているのか分かりません。

そこで **ユーザ情報を取得するAPIを別途使用してつきあわせる** ことで、ユーザの名前を取得します。

joinといった方が分かりやすいかもしれませんね。

しかし、チャンネル毎にユーザの情報を取得してしまうと直ぐにAPI制限に引っかかってしまいます。

そこで今回は全ユーザのリストを予め取得しておき、プログラム内で比較することによって、ユーザ情報取得APIの回数を1回にしています。

こちらはプログラムをやっているとよくあるつきあわせ処理ですので、覚えておくとよいでしょう。

※ プログラム初心者はAPIだけでなくデータベースなどでもループ処理でデータの取得をやってしまい、負荷を上げてしまうことがよくあります!!

#### ページャーの操作

次の注意点としてチャンネル取得APIを使用する際にページャー処理を追加する必要がある点です。

slackのチャンネル取得APIでは一度に1000件の情報しか取得できません。

しかしslackのチャンネルは規模が大きくなってくると1000件を超えてくることがあります。

その際にはページャーの処理を追加してあげる必要があります。

Slack APIのぺージャーの仕組みはレスポンスに「 **response\_metadata.next\_cursor**」の値が入っていれば、まだ次の値が存在するので、API実行時の値としてその「 **cursor**」の値をセットし、再度APIを呼び出すというものです。

もっと細かく知りたい場合は [公式ページ](https://api.slack.com/docs/pagination) をご確認ください。

#### スプレッドシート更新処理

3つ目の注意点はスプレッドシートに値を書き込む際の注意点です。

基本的にリストを更新する際は前に取得した値を一度消し込み、今回取得した値を書き込みます。

そうしないと、取得したリストの数によって、不要なレコードが残ってしまう可能性があるからです。

ですので、今回実行メソッドでは予めスプレッドシートの内容を消し込む処理を入れております！

## 最終的に出来上がったコード

最終的にできたコードは以下になります。

```
// スプレッドシート取得.
let ss = SpreadsheetApp.getActiveSpreadsheet();
let sheet = ss.getActiveSheet();

// SlackAPIで登録したボットのトークンを設定する.
let token = "xoxb-から始まる上記Slack側の設定で取得したトークン";

// 削除するレコード数.
let deleteLength = 2000;

// 実行メソッド.
function setChannelInfo() {
  deleteSpredsheetList();
  users = getSlackUsers();

  public_channels = getChannelList("public_channel");

  for(let i=0;i<public_channels.length;i++) {
      sheet.getRange(i+5,2).setValue(public_channels[i].name);
      sheet.getRange(i+5,3).setValue(public_channels[i].topic.value);
      sheet.getRange(i+5,4).setValue(public_channels[i].purpose.value);

      for(let j=0;j<users.length;j++) {
        if(users[j].id == public_channels[i].creator) {
          sheet.getRange(i+5,5).setValue(users[j].real_name + "(" + users[j].name + ")");
          break;
        }
      }
  }

  sheet.getRange(2,5).setValue(public_channels.length);
}

// slackユーザ取得.
function getSlackUsers() {
  let options = {
    "method": "get",
    "contentType": "application/x-www-form-urlencoded",
    "payload" : {
      "token": token
    }
  }

  // 必要scope = users:read.
  let url = 'https://slack.com/api/users.list';
  let response = UrlFetchApp.fetch(url, options);
  let obj = JSON.parse(response);
  return obj.members;
}

// チャンネルリストを消し込む.
function deleteSpredsheetList() {
  for(let i=0;i<deleteLength;i++) {
    sheet.getRange(i+5,2).setValue("")
    sheet.getRange(i+5,3).setValue("")
    sheet.getRange(i+5,4).setValue("")
    sheet.getRange(i+5,5).setValue("")
  }

  sheet.getRange(2,5).setValue(0);
}

// チャンネルリスト ページャー対応.
function getChannelList(types) {
  // ページャーの最大数.
  let maxPager = 10;
  let count = 1;
  let channelList = [];
  let cursor = "";
  while(count<=maxPager) {
    channelObj = getChannelObject(types, cursor);
    channelList = channelList.concat(channelObj.channels);

    if(channelObj.response_metadata.next_cursor == "") {
      break;
    }

    cursor = channelObj.response_metadata.next_cursor;
    count++;
  }

  return channelList;
}

// チャンネルリストを取得する.
// https://api.slack.com/methods/conversations.list
function getChannelObject(types, cursor) {
  let options = {
    "method" : "get",
    "contentType": "application/x-www-form-urlencoded",
    "payload" : {
      "token": token,
      "limit": 1000,
      "exclude_archived": true,
      "types": types,
      "cursor": cursor
    }
  }

  let url = 'https://slack.com/api/conversations.list';
  let response = UrlFetchApp.fetch(url, options);
  let obj = JSON.parse(response);

  return obj;
}

```

こちらを実行することで、以下のようにスプレッドシートにSlackのチャンネルの情報を一覧化させることができます。

[![GAS - 設定2](https://shiimanblog.com/wp-content/uploads/2021/10/gas2-1-800x407.png)](https://shiimanblog.com/wp-content/uploads/2021/10/gas2-1.png)

## 定期実行

スクリプトの作成が終わったら、最後に定期実行の設定を行いましょう。

今回は毎日1回深夜にスクリプトを実行することによってリストを常に最新化させます。

定期実行の設定はスクリプトを記入した画面から設定することが可能です。

まず左のメニューから「トリガー」を選択し遷移します。

[![GAS - 設定3](https://shiimanblog.com/wp-content/uploads/2021/10/gas3.png)](https://shiimanblog.com/wp-content/uploads/2021/10/gas3.png)

すると遷移した右下の方に「\+ トリガーを追加」ボタンが表示されていますので、こちらをクリックします。

するとポップアップ画面が出てきますので、こちらのポップアップでトリガーを設定します。

実行する関数を選択: **setChannelInfo**

イベントのソースを選択: **時間主導型** 時間ベースのトリーがのタイプを選択: **日付ベースのタイマー**

時刻を選択: **お好きな時間**

[![GAS - 設定4](https://shiimanblog.com/wp-content/uploads/2021/10/gas4.png)](https://shiimanblog.com/wp-content/uploads/2021/10/gas4.png)

こちらを設定して保存すれば対象のスクリプトが定時実行されます。

深夜にでもスクリプトを仕込んでおけば、寝ている間にリストが更新されるってことですね。

## まとめ

今回はSlackボットをGASから操作して、Slackのpublicチャンネル一覧を取得し、スプレッドシートに吐き出す方法を紹介しました。

GASでスプレッドシートを操作することはよく行う方法ですので、今回の記事を参考に是非いろいろ試してみてください。

ボットやGASなどはアイデア次第で業務の効率化が進みますので、「 **なんか毎日同じ作業を繰り返し行っているな**」とか、「 **slackとスプレッドシートを連携させたいな**」と思ったらボットやGASでどうにかできないか考えてみてください。

今回の記事が業務の効率化や改善に少しでも役立ちましたら幸いです！