---
id: 1759
title: 【AWS】Cloud FrontのログをAthenaでクエリする方法
slug: cloud-front-log-with-athena
status: publish
date: 2022-10-08T19:30:00
modified: 2022-10-07T18:03:58
excerpt: S3に保存されたCloudFrontのアクセスログをAthenaでクエリする方法を、テーブル作成からクエリ実行まで解説します。
categories: [19, 18, 21]
tags: [124, 125, 126]
featured_media: 1760
---

こんばんは、しーまんです。

前回の記事でALBのアクセスログをAthenaでクエリする方法を紹介しました。

まだご覧になっていない方はぜひこちらも御覧ください。

 [![](https://shiimanblog.com/wp-content/uploads/2022/10/eyecatch_alb_log_with_athena-320x180.jpg)\
\
【AWS】Application LoadbalancerのログをAthenaでクエリする方法\
\
AWSでApplication Loadbalancerのアクセスログをs3に保存し、Athenaでクエリする方法を解説致します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2022.10.07](https://shiimanblog.com/engineering/alb-log-with-athena/ "【AWS】Application LoadbalancerのログをAthenaでクエリする方法")

さて今回はAthenaでログをクエリする紹介第2弾になります。前回はALBのログでしたが、今回は [Cloud Front](https://aws.amazon.com/jp/cloudfront/) のアクセスログを [Athena](https://aws.amazon.com/jp/athena) を使用してクエリしていこうと思います。

AWSでのログ運用で困っている方の参考になればいいなと思います。

## 事前準備

まず前提としてCloud Front(CF)のアクセスログをS3に配置済みの状況とします。

もしまだCFのログ出力を有効にしていない場合は有効にしておきましょう。

参考までに [Terraform](https://www.terraform.io/) でのログの有効化設定の例を載せておきます。

```
resource "aws_cloudfront_distribution" "cloudfront_distribution" {
   ....省略

  logging_config {
    include_cookies = false
    bucket                 = "cf-log-bucket"
    prefix                   = "test-cf"
  }
}
```

上記の設定だとCFのアクセスログは以下のパスに保存されます。

```
s3://cf-log-bucket/test-cf/xxxx.yyyy-MM-dd.xxxx.log.gz
```

## Athena

さて、出力されているログの情報が分かりましたので、ここからAthenaでクエリできるように設定していきます。

AthenaとはサーバーレスでS3内のデータに対して標準SQLを使用して簡単に分析できるサービスです。

こちらを使用するためにはまずデータベースを作成し、次にテーブルを作成します。

作成したテーブルに対してクエリを発行し分析を行うという形ですね。

ですのでまずはデータベースを作成していきます。

Terraformで作成する場合は下記のように `name` と `bucket` を指定することで作成可能です。

```
resource "aws_athena_database" "athena_database" {
  name   = "cf_athena_db"
  bucket = "athena_query_log_bucket"
}
```

bucketはCFログの配置されているbucketとは別物です。

athena自身のクエリログを出力するためのbucketになります。

## Lambdaでパーティション射影対応

データベースの作成が終わったら、テーブルの作成に行きたいところですが、今回は配置されたログに対して、加工処理を入れていきます。

ログの加工と言ってもログ自体の加工ではなく、ログのS3配置場所の移動になります。

これはなぜ必要かというと、Athenaでクエリを発行する場合、クエリを発行する際にスキャンデータ量を抑えるためにパーティションを作成することが推奨されています。

そのパーティションを自動で作成するための仕組みとして [パーティション射影](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/partition-projection.html) という機能があるのですが、その機能を利用するためにはデフォルトのS3へのパスでは対応できません。

ですので、パーティション射影によるパーティション作成が行えるようにログのパスを変更してあげる必要があります。

ではパスの移動を行うためのLambda関数を作っていきましょう。

[以下のソース](https://github.com/aws-samples/amazon-cloudfront-access-logs-queries) を使用してLambdaを作成します。(Node.js)

```
// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
const aws = require('aws-sdk');
const s3 = new aws.S3({ apiVersion: '2006-03-01' });

const datePattern = '[^\\d](\\d{4})-(\\d{2})-(\\d{2})-(\\d{2})[^\\d]';
const filenamePattern = '[^/]+$';

exports.handler = async (event, context, callback) => {
    const moves = event.Records.map(record => {
        const bucket = record.s3.bucket.name;
        const sourceKey = record.s3.object.key;
        const sourceRegex = new RegExp(datePattern, 'g');
        const match = sourceRegex.exec(sourceKey);

        if (match == null) {
            console.log(`Object key ${sourceKey} does not look like an access log file, so it will not be moved.`);
        } else {
            const [, year, month, day,] = match;
            const filenameRegex = new RegExp(filenamePattern, 'g');
            const filename = filenameRegex.exec(sourceKey)[0];
            const bucket_path = sourceKey.substr(0, sourceKey.length - filename.length - 1)
            const targetKey = `copy/${bucket_path}/${year}/${month}/${day}/${filename}`;
            console.log(`Copying ${sourceKey} to ${targetKey}.`);

            const copyParams = {
                CopySource: bucket + '/' + sourceKey,
                Bucket: bucket,
                Key: targetKey
            };
            const copy = s3.copyObject(copyParams).promise();
            const deleteParams = { Bucket: bucket, Key: sourceKey };

            return copy.then(function () {
                console.log(`Copied. Now deleting ${sourceKey}.`);
                const del = s3.deleteObject(deleteParams).promise();
                console.log(`Deleted ${sourceKey}.`);
                return del;
            }, function (reason) {
                var error = new Error(`Error while copying ${sourceKey}: ${reason}`);
                callback(error);
            });
        }
    });
    await Promise.all(moves);
};
```

そして [S3のイベント通知](https://docs.aws.amazon.com/ja_jp/AmazonS3/latest/userguide/NotificationHowTo.html) の設定により「 **すべてのオブジェクト作成イベント**」で作成したLambda関数を呼び出します。するとログが配置されると同時にLambdaが起動し、ファイルが下記パスに移動されます。

```
// 移動前
s3://cf-log-bucket/test-cf/xxxx.yyyy-MM-dd.xxxx.log.gz

// 移動後
s3://cf-log-bucket/copy/test-cf/yyyy/MM/dd/xxxx.yyyy-MM-dd.xxxx.log.gz
```

これでパーティション射影の準備が完了しました！

## テーブル作成

Lambdaの作成が終わったら、いよいよクエリを行うためのテーブルを作成していきます。

Athenaは [Glueデータカタログ](https://aws.amazon.com/jp/glue/) と統合されているためそちらのリソースで定義してもよいのですが、 [Create SQLで定義](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/cloudfront-logs.html) しておいた方が分かりやすいため、SQLで用意します。

今回用意するSQLもパーティション射影の設定を追加しています。

サンプルとして下記の **create\_cf\_log\_table.sql.tpl** を用意します。

```
CREATE EXTERNAL TABLE IF NOT EXISTS ${athena_database_name}.${athena_table_name} (
    `date` DATE,
    `time` STRING,
    `location` STRING,
    `bytes` BIGINT,
    `request_ip` STRING,
    `method` STRING,
    `host` STRING,
    `uri` STRING,
    `status` INT,
    `referrer` STRING,
    `user_agent` STRING,
    `query_string` STRING,
    `cookie` STRING,
    `result_type` STRING,
    `request_id` STRING,
    `host_header` STRING,
    `request_protocol` STRING,
    `request_bytes` BIGINT,
    `time_taken` FLOAT,
    `xforwarded_for` STRING,
    `ssl_protocol` STRING,
    `ssl_cipher` STRING,
    `response_result_type` STRING,
    `http_version` STRING,
    `fle_status` STRING,
    `fle_encrypted_fields` INT,
    `c_port` INT,
    `time_to_first_byte` FLOAT,
    `x_edge_detailed_result_type` STRING,
    `sc_content_type` STRING,
    `sc_content_len` BIGINT,
    `sc_range_start` BIGINT,
    `sc_range_end` BIGINT
)
PARTITIONED BY
(
    `day` STRING
)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY '\t'
LOCATION 's3://${log_bucket_path}/'
TBLPROPERTIES (
    "skip.header.line.count" = "2",
    "projection.enabled" = "true",
    "projection.day.type" = "date",
    "projection.day.range" = "2022/01/01,NOW",
    "projection.day.format" = "yyyy/MM/dd",
    "projection.day.interval" = "1",
    "projection.day.interval.unit" = "DAYS",
    "storage.location.template" = "s3://${log_bucket_path}/\$\{day\}"
)
```

こちらのSQLテンプレートをTerraformで反映していきます。

```
data "aws_caller_identity" "current_account_id" {}

data "aws_s3_bucket" "cf_log_bucket" {
  bucket = "cf-log-bucket"
}

data "template_file" "create_cf_log_table_sql" {
  template = file("./create_cf_log_table.sql.tpl")
  vars = {
    athena_database_name = resource.aws_athena_database.athena_database.name
    athena_table_name         = "cf_logs"
    log_bucket_path               = "${data.aws_s3_bucket.cf_log_bucket.id}/copy/test-cf"
    aws_account_id                = data.aws_caller_identity.current_account_id.account_id
    region                                 = "ap-northeast-1"
  }
}

resource "aws_athena_named_query" "athena_named_query" {
  name        = "create_cf_log_table"
  database = resource.aws_athena_database.athena_database.id
  query       = data.template_file.create_cf_log_table_sql.rendered
}
```

上記のTerraformファイルを実行するとAthenaの「保存したクエリ」に `create_cf_log_table` ができていますので、こちらを実行してテーブル作成します。

## テーブルに対してクエリしてみる

最後にきちんと設定できているかどうか確認するために実際にクエリを発行してみましょう。

```
SELECT * FROM "cf_athena_db"."cf_logs" limit 10;
```

表示されれば正しくテーブルの設定ができています。

あとは解析したいクエリをこちらのテーブルに対して実行するだけになります。

クエリサンプルに関しては [公式ドキュメント](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/cloudfront-logs.html) に少し例がありますので、参考にしてください。

## まとめ

今回は、Cloud FrontのログをAthenaでクエリする方法について解説しました。

Cloud Frontのログ出力はパーティション射影を活用する形でアウトプットしてくれません。(※投稿時点)ですので今回はLambdaを使用することでパーティション射影を利用できるようにアウトプットしています。

今までCFのログをS3には配置していたけど、出力していただけで利用していなかったという方がいるのではないでしょうか？こちらのログを分析に使用することで多くの利点があります。

是非今回の設定を試して、いろいろな解析に利用してみてください。

今回の記事がAWSを使用する方にとって少しでも参考になりましたら幸いです。