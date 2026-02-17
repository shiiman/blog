---
id: 1049
title: 【2021年】ボットを作成してSlack APIをGASから叩く方法 | SlackAppライブラリ未使用
slug: gas-for-slack-bot
status: publish
date: 2021-10-04T19:30:00
modified: 2021-10-01T19:22:16
excerpt: Google Apps Script（GAS）からSlack APIを叩いてBotを作成する方法を、ライブラリ未使用で一から解説します。
categories: [18, 64, 65]
tags: [66, 67, 68]
featured_media: 1051
---

こんばんは、しーまんです。

最近多くの会社ではSlackというチャットツールを使用してコミュニケーションをしてるところが増えました。コロナ化時代ではチャットツールの使用が必須化されているということですね。

私もSREという職業をしているとSlackを使うことがよくあり、ボットを作成したりAPIを叩いたりなどはよくしております。

そこで今回はSlackボットの作成方法とGASを使用して簡単にSlackAPIを叩くまでを紹介していこうと思います。

## やることとやらないこと

SlackのAPIに関してはとても数が多いので、今回の記事だけでは紹介しきれません。

またボットの機能や、GASについても全てを紹介しだしたらきりがありません。

ですので最初に今回紹介することとしないことを明確にしておきたいと思います。

### やること

- Slackボットの作成
- GASでSlackに対してユーザ取得とDMを送るAPIを実装

### やらないこと

▼ Slack

- incoming webhook

  -\> 手軽にslackと連携が出来るが、webhookの管理や管理者自体の管理が煩雑になるので現在非推奨です。よって今回も扱いません

- Interactive Message

  -\> 対話的なボットの作成に使用するInteractive Messageですが今回はGASからの片方向通信で作成します。よくGASをpublicで公開するような記事を見かけますが、こちらはちゃんと作らないとセキュリティーリスクになりますので、今回は扱いません。

▼ GAS

- SlackAppライブラリの使用

  -\> 使用できるAPIが少なすぎるため今回は扱いません。

- スプレッドシードとの連携

  -\> 本来GASはスプレッドシートと連携してこそ力を発揮しますが、今回はSlack APIを叩くために使用しますので、今回は扱いません。

ということで前提条件を設定できましたので、実際に作成してきましょう。

## Slackボット作成

まずはSlack側でbotの作成をしていきます。

下記からslack appsのページを開きます。

 [![](https://s.wordpress.com/mshots/v1/https%3A%2F%2Fapi.slack.com%2Fapps?w=320&h=180)\
\
Slack API: Applications \| Slack\
\
![](https://www.google.com/s2/favicons?domain=https://api.slack.com/apps)\
\
api.slack.com](https://api.slack.com/apps "Slack API: Applications | Slack")

画面が遷移すると上部にある「 **Create New App**」をクリックします。

[![slackbot -設定1](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot1-800x235.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot1.png)

するとポップアップが表示されます。

アプリの作成方法を選択できますので、こちらは「 **From scratch**」を選択します。

[![slackbot -設定2](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot2-800x555.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot2.png)

次にアプリ名と連携するworkspaceを選択していきます。

アプリ名は自由に設定してください。

入力できたら、「 **Create App**」をクリックします。

[![slackbot -設定3](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot3-800x736.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot3.png)

アプリの作成ができると「 **Basic Information**」画面が表示されていると思います。

ここからボットの設定をしていきます。

まず左の項目かBasec Informationの右下にある「 **Permissions**」の画面に遷移し、ボットのSlackに対する権限を設定していきましょう。

[![slackbot -設定4](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot4-800x662.jpg)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot4.jpg)

遷移して少し下にスクロールすると「 **Scopes**」->「 **Bot Token Scopes**」をいう権限を設定する箇所が見つかります。

ここから「 **Add an Oauth Scope**」をクリックして権限を追加しましょう。

今回は「 **users.list**」と「 **conversations.open**」と「 **chat.postMessage**」3つのAPIを使用しますので、「 **users:read**」と「 **im:write**」と「 **chat:write**」の3つの権限を追加します。

APIに必要な権限は公式のドキュメントを確認しましょう。

「 [users.list](https://api.slack.com/methods/users.list)」「 [conversations.open](//api.slack.com/methods/conversations.open)」 「 [chat.postMessage](https://api.slack.com/methods/chat.postMessage)」

[![slackbot -設定5](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot5-1.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot5-1.png)

権限の設定が終わったら、ページの一番上に戻りworkspaceにアプリをインストールしていきます。

「 **Install to Workspace**」をクリックしましょう。

[![slackbot -設定6](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot6.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot6.png)

クリックすると連携確認画面に遷移しますので、問題なければ「 **許可する**」をクリックします。

[![slackbot -設定7](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot7.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot7.png)

連携が完了すると画面が戻り「 **OAuth Tokens for Your Workspace**」に「 **Bot User OAuth Token**」が発行されていると思います。こちらのトークンを後のGASで使用しますので、コピーしておきましょう。

[![slackbot -設定8](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot8.png)](https://shiimanblog.com/wp-content/uploads/2021/10/slackbot8.png)

こちらでボットの最低限の設定は終わりです。

## GASの作成

Slack側の設定が完了しましたので、次にGAS側の設定をしていきます。

Googleドライブ上で右クリックすると作成したいファイルの種類が表示されます。

ここから「 **Google Apps Script**」を選択して開きましょう。

[![GAS - 設定1](https://shiimanblog.com/wp-content/uploads/2021/10/gas1-800x688.png)](https://shiimanblog.com/wp-content/uploads/2021/10/gas1.png)

GASの編集画面が開いたら下記コードを登録していきます。

```
// SlackAPIで登録したボットのトークンを設定する.
let token = "xoxb-から始まる上記Slack側の設定で取得したトークン";
let method = "post";
let contentType = "application/x-www-form-urlencoded";

// 取得したslackのユーザから特定の名前のユーザにDMを送る
function postTestMessage() {
  let members = getSlackUsers();
  let message = "通知テスト!!!";

  for(let i=0;i<members.length;i++) {
    if(members[i].name == "xxxxx_xxxxx") {

      console.log(members[i].id);
      let channel_id = getChannelID(members[i].id);
      console.log(channel_id);
      postMessage(channel_id, message);
      break;
    }
  }
}

// slackユーザの取得.
// https://api.slack.com/methods/users.list
function getSlackUsers() {
  let options = {
    "method" : method,
    "contentType": contentType,
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

// メンバーIDを受け取りチャンネルIDを返す.
// https://api.slack.com/methods/conversations.open
function getChannelID(user_id) {
  let options = {
    "method" : method,
    "contentType": contentType,
    "payload" : {
      "token": token,
      "users": user_id
    }
  }

  // 必要scope = im:write.
  let url = 'https://slack.com/api/conversations.open';
  let response = UrlFetchApp.fetch(url, options);
  let obj = JSON.parse(response);

  return obj.channel.id;
}

// メッセージを送る.
// https://api.slack.com/methods/chat.postMessage
function postMessage(channel_id, message) {
  let options = {
    "method" : method,
    "contentType": contentType,
    "payload" : {
      "token": token,
      "channel": channel_id,
      "text": message
    }
  }

  // 必要scope = chat:write.
  let url = 'https://slack.com/api/chat.postMessage';
  UrlFetchApp.fetch(url, options);
}
```

簡単に説明すると一番上で token の設定をしています。これはSlackボットの設定時に取得したトークンになります。

あとは3つのAPIを呼び出すAPIをそれぞれ実装し、最後に postTestMessage という実際にテストメッセージを送るメソッドを追加しました。このメソッドを呼び出すことでSlackに通知を送ることができます。

成功すると下記のようなメッセージが送られてきます。

ちゃんと作成したボットからダイレクトメッセージが送られてきていますね。

[![GAS - 設定2](https://shiimanblog.com/wp-content/uploads/2021/10/gas2.png)](https://shiimanblog.com/wp-content/uploads/2021/10/gas2.png)

## まとめ

今回はSlackのボット作成方法とGASを使ってSlack APIを使用する方法を紹介しました。

今回のボットは一番シンプルなものを作成しましたが、アイデア次第で色々便利なボットを作成することが可能です。

例えばスプレッドシートからデータを読み込んでSlackに通知するとか、インタラクティブにボットに質問したら返答を返すとか。ボットの機能は無限大です。

業務の効率化を図るにはボットの作成はもってこいですので、ぜひ皆さんもいろいろなボットを作成してみてください。