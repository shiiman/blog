---
id: 1494
title: 【Snowflake】snowpipeを使用して継続的にAWS S3からデータをロードする方法
slug: snowflake-snowpipe
status: publish
excerpt: こんばんは、しーまんです。 今回はSnowflakeのメイン機能であるsnowpipeについて紹介していきます。こちらの機能を使いこなすことによってSnowflakeでの分析がかなり容易になるかと思います。 Snowfl \[…\]
categories:
    - 18
    - 77
tags:
    - 78
    - 79
    - 97
featured_media: 1496
date: 2022-01-28T19:30:00
modified: 2022-01-19T19:41:38
---

こんばんは、しーまんです。

今回はSnowflakeのメイン機能であるsnowpipeについて紹介していきます。

こちらの機能を使いこなすことによってSnowflakeでの分析がかなり容易になるかと思います。

Snowflakeを使用して分析環境を構築しようかと検討している方には参考になる部分があると思いますのでぜひご覧いただければと思います。

## snowpipeとは

AWSにおけるS3用のsnowpipeとはS3バケットのSQS(Simple Queue Servive)通知を使用して **データのロードを自動化する仕組み** です。

![AWSでのsnowpipeの仕組み](https://shiimanblog.com/wp-content/uploads/2022/01/329671c453d6d608029602532e278fec.png)

設定方法については [ドキュメント](https://docs.snowflake.com/ja/user-guide/data-load-snowpipe-auto-s3.html) が充実しておりますので、一読しておくとよいかと思います。

## snowpipe設定手順

前述のドキュメントでは手順が多く、実際に設定する場合、最低限何をすればよいか分かりづらいので、簡単に手順をリスト化しておきます。

1. データベース/スキーマ/テーブルの作成
2. ストレージ統合(ステージ作成)
3. pipeline作成(snowpipe)
4. aws側で対象のS3に対してSQSの設定

こちらの4つの手順でデータのロードが継続的に自動化されます。今回はロードするファイルはjson形式という前提で話を進めていと思います。

では一つ一つ手順をみていきましょう！

### データベース/スキーマ/テーブルの作成

まずはデータをロードするテーブルを作成します。

テーブルを作成するまでの権限周りなどは別記事を公開しておりますので、そちらを参考にお願いします。

[![](https://shiimanblog.com/wp-content/uploads/2021/10/eyecatch_snowflake_user-320x180.jpg)\
\
【Snowflake】ユーザの作成からロール権限設定までまとめてみた \| 初期設定\
\
Snowflakeはデータを一元管理するためのクラウドデータプラットフォームです。そのアクセス管理権限について今回はまとめてみました。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.10.22](https://shiimanblog.com/engineering/snowflake-permission/ "【Snowflake】ユーザの作成からロール権限設定までまとめてみた | 初期設定")

例えば下記のようなテーブルを用意します。

```
CREATE OR REPLACE TABLE test_database.test_schema.test_table(
    user_id number,
    name varchar,
    create_at timestamp
);
```

### ストレージ統合(ステージ作成)

次に行うのがストレージ統合です。

ストレージ統合には以下の3点の設定が必要です。

1. インテグレーション作成
2. ステージ作成
3. AWSでIAMロール作成

それぞれの設定方法については以前、別記事として投稿しておりますのでそちらを参照ください。

[![](https://shiimanblog.com/wp-content/uploads/2021/10/eyecatch_snowflake_s3-320x180.jpg)\
\
【Snowflake】AWS S3からデータをロードし操作する方法 \| 初期設定\
\
AWSのS3上にあるデータをSnowflakeからアクセスして操作する方法を解説します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.10.26](https://shiimanblog.com/engineering/snowflake-s3/ "【Snowflake】AWS S3からデータをロードし操作する方法 | 初期設定")

ここでは **test\_stage** というステージを作成したとします。

### pipeline作成(snowpipe)

次にpipelineを設定していきます。

こちらの設定が実質snowpipeの設定になります。

```
CREATE OR REPLACE PIPE test_database.test_schema.test_pipe AUTO_INGEST=TRUE AS
  COPY INTO test_database.test_schema.test_table(user_id, name, create_at)
  FROM (select
    $1:user_id::number,
    $1:name::varchar,
    $1:create_at::timestamp
   FROM @test_database.test_schema.test_stage/ t)
  FILE_FORMAT = (TYPE = 'JSON');
```

pipelineが構築されるとS3からの通知を受け取る用のSQSが自動で設定されます。

AWS側の設定の際にSQSのARNが必要になりますので、確認しておきましょう！

```
use schema test_database.test_schema;
SHOW PIPES;
```

すると **notification\_channel** にSQSのARNが表示されますので、控えておきましょう！

ここまで設定すればSnowflake側の設定は完了です。

あとはAWS側の設定をしていきましょう！

### AWS側で対象のS3に対してSQSの設定

最後にAWS側でS3にSQS通知の設定をしていきます。

対象のS3バケットを選択し、プロパティタブを開きます。

[![snowpipe - 設定1](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting1.png)](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting1.png)

次に「 **イベント通知**」設定をしていきます。

「 **イベント通知を作成**」をクリックします。

[![snowpipe - 設定2](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting2-800x143.png)](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting2.png)

するとイベント通知作成画面が開きます。

必要に応じて「 **イベント名**」「 **プレフィックス**」「 **サフィックス**」を設定しましょう。

今回はjsonファイルのみをロードさせる設定にするためにサフィックスに **.json** と指定しています。

[![snowpipe - 設定3](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting3-800x505.png)](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting3.png)

次にイベントタイプの設定です。

「 **全てのオブジェクト作成イベント**」にチェックを入れましょう。

[![snowpipe - 設定4](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting4-800x373.png)](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting4.png)

最後に送信先の設定です。

送信先を「SQSキュー」に設定し、SQSキュー欄に先程控えておいたSQSのARNを入力しましょう！

[![snowpipe - 設定5](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting5-800x552.png)](https://shiimanblog.com/wp-content/uploads/2022/01/snowpipe_setting5.png)

以上でsnowpipeの設定は完了です。

思ったよりも簡単に出来ましたよね。

## 動作確認

設定が出来たら必ず動作確認を行いましょう。

下記のようなテスト用のjsonファイルを用意し、S3にアップロードします。

```
{"user_id":1,"name":"shiiman","create_at":"2022-01-28T12:00:00"}
```

Snowflake側でテーブルにデータがロードできていることが確認できるはずです。

## まとめ

今回はSnowflakeの主要機能であるsnowpipeの設定方法について解説してきました。

こちらの機能を使いこなすことがSnowflakeで分析環境を作る上では一番手軽かつ強固な設定になります。基本的にはSQLでの設定になり権限周りでのミスも起こりやすいですが、やっていることは単純ですので、1ステップずつ確認してみてください。

今回の記事がSnowflakeを使用する方の参考に少しでもなりましたら幸いです。