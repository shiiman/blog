---
id: 1210
title: 【Snowflake】AWS S3からデータをロードし操作する方法 | 初期設定
slug: snowflake-s3
status: publish
date: 2021-10-26T19:30:00
modified: 2021-10-22T14:36:50
excerpt: AWS S3に配置されたCSVデータをSnowflakeにロードする方法を、外部ステージの設定からCOPYコマンドまで解説します。
categories:
    - 18
    - 77
tags:
    - 78
    - 79
featured_media: 1215
---

こんばんは、しーまんです。

Snowflakeでは外部のデータをロードすることでその恩恵を最大限活かすことができます。

今回はS3に配置されたCSVのデータをSnowflake上にロードする方法を紹介していきます。

外部連携周りの設定で苦労されている方の参考になりましたら幸いです。

## テーブル作成方法

Snowflakeでのテーブル作成や権限周りについては別記事に公開しておりますので、こちらを参照ください。

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
2021.10.22](https://shiimanblog.com/engineering/snowflake-permission "【Snowflake】ユーザの作成からロール権限設定までまとめてみた | 初期設定")

## データロード設定

SnowflakeにS3からデータをロードするためにはストレージ統合を行う必要があります。

他にも設定方法はありますが、ストレージ統合が一番セキュアですのでこちらの方法で設定することをおすすめします。

他の設定方法について知りたい場合は [ドキュメント](https://docs.snowflake.com/ja/user-guide/data-load-s3-config.html) を参考にしてください。

設定手順は以下の流れになります。

1. インテグレーション作成
2. ステージ作成
3. AWSでIAMロール作成

[ドキュメント](https://docs.snowflake.com/ja/user-guide/data-load-s3-config-storage-integration.html) ではまず先に仮のIAMロールを作成してからインテグレーションを作成し、その情報を元にIAMロールの設定を変更しています。この手順でもよいのですが、SnowflakeとAWSを行ったりきたりしないといけないので、上記の手順にしてステップ数を減らしています。

### インテグレーション

ではまずステップ１.インテグレーションの設定をしていきましょう。

作成時に予め決めておく必要があるのは後ほど作成する「IAMロール名」と「接続するバケット名」になります。

それぞれ「 **storage\_aws\_role\_arn**」と「 **storage\_allowed\_locations**」に設定しましょう。

※ 許可するバケットは複数設定可能です。

```
-- アカウントを操作できるロールに変更.
use role accountadmin;

--s3連携 development
create storage integration test_s3_int
  type = external_stage
  storage_provider = s3
  enabled = true
  storage_aws_role_arn = 'arn:aws:iam::<AWSアカウントID>:role/snowflake-role'
  storage_allowed_locations = ('s3://test-bucket/');

-- 作成したインテグレーションの情報を取得.
desc integration test_s3_int;
```

作成したインテグレーションの情報をこの後のAWSのIAMロール作成時に使用します。

「 **STORAGE\_AWS\_IAM\_USER\_ARN**」と「 **STORAGE\_AWS\_EXTERNAL\_ID**」の内容をメモしておきましょう。

[![DESC - INTEGRATION](https://shiimanblog.com/wp-content/uploads/2021/10/desc_integration-800x177.png)](https://shiimanblog.com/wp-content/uploads/2021/10/desc_integration.png)

### ステージ作成

インテグレーションの作成が終わったら次はステージの作成です。

このステージオブジェクトを通してSnowflakeのテーブルにデータをロードします。

```
create or replace stage test_stage
  storage_integration = test_s3_int
  url = 's3://test-bucket/csv/';
```

### IAM Role作成

最後にAWS側でIAMロールの作成を行います。

IAMロールの作成方法は [ドキュメント](https://docs.snowflake.com/ja/user-guide/data-load-s3-config-storage-integration.html#step-3-create-a-cloud-storage-integration-in-snowflake) どおりになります。

まずはポリシーを作成しその後ロールを作成していきます。

ロール作成時にインテグレーション作成時に設定した「IAMロール名」で作成するのと、「 **STORAGE\_AWS\_IAM\_USER\_ARN**」「 **STORAGE\_AWS\_EXTERNAL\_ID**」をそれぞれ設定するのを忘れないようにしましょう。

これでデータをロードする設定が完了しました。

## S3からデータロードをテストする

それでは先程設定したステージを使用して、テストテーブルにデータをロードしてみましょう。

データをロードするには「 **copy into <テーブル名> from @<ステージ名>/<ファイルパス> file\_format …**」でロードを行います。

今回は「s3://test-bucket/csv/test.csv」をtest\_tableにロードする例になります。

```
copy into test_schema.test_table
 from @test_stage/test.csv
 file_format = (
     type = csv
     skip_header = 1
     field_optionally_enclosed_by = '0x22'
 );
```

ファイルフォーマットがcsvで1行名はカラム名が入る想定なのでデータとしてはスキップしています。

またCSVのデータでtimestamp型のデータが「”(ダブルクウォート)」で囲まれているとテーブル定義と型が違うというエラーが出てしまうので、それを回避するために [field\_optionally\_enclosed\_by](https://docs.snowflake.com/ja/user-guide/data-load-considerations-load.html#csv-data-trimming-leading-spaces) を設定しています。

エラーが出なければS3からSnowflakeへのデータロードが完了するはずです。

## まとめ

今回はAWSのS3上に配置されたCSVファイルからSnowflakeにデータをロードする方法を解説いしました。

ストレージの連携にはいくつか手順が必要ですが、1度設定してステージを汎用的に利用すれば外部ストレージのデータを簡単にロードすることが可能です。

Snowflakeを使用する上で外部ストレージとの連携は必須ですので、ぜひ連携方法についてはおさえておきましょう！！

以上、Snowflakeの外部連携について参考になれば幸いです。