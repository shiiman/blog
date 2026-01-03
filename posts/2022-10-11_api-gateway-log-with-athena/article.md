---
id: 1772
title: 【AWS】API GatewayのログをAthenaでクエリする方法
slug: api-gateway-log-with-athena
status: publish
excerpt: こんばんは、しーまんです。 前回の記事でCloud FrontのアクセスログをAthenaでクエリする方法を紹介しました。まだご覧になっていない方はぜひこちらも御覧ください。 【AWS】Cloud FrontのログをAt \[…\]
categories:
    - 19
    - 18
    - 21
tags:
    - 124
    - 125
    - 127
featured_media: 1773
date: 2022-10-11T19:30:00
modified: 2022-10-11T11:27:40
---

こんばんは、しーまんです。

前回の記事でCloud FrontのアクセスログをAthenaでクエリする方法を紹介しました。

まだご覧になっていない方はぜひこちらも御覧ください。

 [![](https://shiimanblog.com/wp-content/uploads/2022/10/eyecatch_cloudfront_log_with_athena-320x180.jpg)\
\
【AWS】Cloud FrontのログをAthenaでクエリする方法\
\
AWSでCloud Frontのアクセスログをs3に保存し、Athenaでクエリする方法を解説致します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2022.10.08](https://shiimanblog.com/engineering/cloud-front-log-with-athena/ "【AWS】Cloud FrontのログをAthenaでクエリする方法")

さて今回はAthenaでログをクエリする紹介第3弾になります。前回はCloud Frontのログでしたが、今回は [API Gateway](https://aws.amazon.com/jp/api-gateway/) のアクセスログを [Athena](https://aws.amazon.com/jp/athena) を使用してクエリしていこうと思います。

AWSでのログ運用で困っている方の参考になればいいなと思います。

## 事前準備

まず前提としてAPI GatewayのアクセスログをS3に配置する必要があります。

API Gatewayはアクセスログの出力先として「 **[Cloud Watch Logs](https://docs.aws.amazon.com/ja_jp/AmazonCloudWatch/latest/logs/WhatIsCloudWatchLogs.html)**」か「 **[Kinesis Data Firehose](https://aws.amazon.com/jp/kinesis/data-firehose/)**」を選択できます。執筆時点ではS3に直接アウトプットすることはできません。

そこで今回はKinesis Data Firehoseを経由してS3にログを配置する設定を行っていきます。

ここで注意点です。Kinesis Data FirehoseでそのままS3に出力してしまうと、JSONログが繋がって出力されてしまいます。

Athenaでクエリするためには [JSONL](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/json-serde.html) で配置する必要があります。つまりJSON1つ毎に改行を入れる必要があるということですね。そこでJSON末尾に改行を入れるために [Lambda](https://aws.amazon.com/jp/lambda/) を使用します。

### Lambda

ということでまずは、JSONの末尾に改行を入れるLambdaを作成します。

Node.jsを使用すると下記のようなコードで行うことが可能です。

```
// https://dev.classmethod.jp/articles/convert-records-in-json-format-from-kinesis-data-firehose-to-json-lines-format-and-output-them-nodejs/

exports.handler = (event, context, callback) => {
    const output = event.records.map((record) => {

        //レコードのJSONの末尾に改行コード（\n）を連結する
        let entry = (Buffer.from(record.data, 'base64')).toString('utf8');
        let result = entry + '\n'
        const payload = (Buffer.from(result, 'utf8')).toString('base64');

        return {
            recordId: record.recordId,
            result: 'Ok',
            data: payload,
        };
    });
    console.log(`Processing completed.  Successful records ${output.length}.`);
    callback(null, { records: output });
};
```

TerraformでLambdaを作成

```
resource "aws_lambda_function" "lambda_function" { ....省略 }
```

### Kinesis Data Firehose

続いて、作成したLambdaを呼び出すようにKinesis Data Firehoseを作成します。

サンプルは以下です。

```
resource "aws_kinesis_firehose_delivery_stream" "delivery_stream" {
   ....省略

   processing_configuration {
      enabled = "true"

      processors {
        type = "Lambda"

        parameters {
          parameter_name  = "LambdaArn"
          parameter_value = "resource.aws_lambda_function.lambda_function.arn:$LATEST"
        }
      }
    }
}
```

### API Gateway

Kinesis Data Firehoseの設定が終わったら、API Gatewayの設定を変更してログをS3に配置できるようにしましょう。参考までに [Terraform](https://www.terraform.io/) でのログの有効化設定の例を載せておきます。

```
locals {
  api_gateway_log = {
    requestId                  = "$context.requestId"
    extendedRequestId = "$context.extendedRequestId"
    ip                               = "$context.identity.sourceIp"
    caller                         = "$context.identity.caller"
    user                           = "$context.identity.user"
    requestTime            = "$context.requestTime"
    httpMethod             = "$context.httpMethod"
    resourcePath           = "$context.resourcePath"
    status                       = "$context.status"
    protocol                    = "$context.protocol"
    responseLength      = "$context.responseLength"
  }
}

resource "aws_api_gateway_stage" "api_gateway_stage" {
  ....省略

  access_log_settings {
    destination_arn = resource.aws_kinesis_firehose_delivery_stream.delivery_stream.arn
    format                = jsonencode(local.api_gateway_log)
  }
}
```

ここまでの設定でアクセスログは以下のパスに保存されます。

```
s3://apigateway-log-bucket/test-apigateway/data/yyyy/MM/dd/xxxx.gz
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
  name   = "apigateway_athena_db"
  bucket = "athena_query_log_bucket"
}
```

bucketはCFログの配置されているbucketとは別物です。

athena自身のクエリログを出力するためのbucketになります。

## テーブル作成

データベースの作成が終わったら、いよいよクエリを行うためのテーブルを作成していきます。

Athenaは [Glueデータカタログ](https://aws.amazon.com/jp/glue/) と統合されているためそちらのリソースで定義してもよいのですが、 [Create SQLで定義](https://docs.aws.amazon.com/ja_jp/apigateway/latest/developerguide/set-up-logging.html) しておいた方が分かりやすいため、SQLで用意します。

またクエリを発行する際にスキャンデータ量を抑えるためにパーティションを作成することが推奨されています。以前は手動で作成したり、Lambdaを利用して定期的に作成する必要がありました。しかし2020年のアップデートで [パーティション射影](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/partition-projection.html) という機能がリリースされ、自動でパーティションを追加できる仕組みが追加されました。

今回用意するSQLもパーティション射影の設定を追加しています。

サンプルとして下記の **create\_apigateway\_log\_table.sql.tpl** を用意します。

```
CREATE EXTERNAL TABLE IF NOT EXISTS ${athena_database_name}.${athena_table_name} (
    `requestId` STRING,
    `extendedRequestId` STRING,
    `ip` STRING,
    `caller` STRING,
    `user` STRING,
    `requestTime` STRING,
    `httpMethod` STRING,
    `resourcePath` STRING,
    `status` STRING,
    `protocol` STRING,
    `responseLength` STRING
)
PARTITIONED BY
(
    `day` STRING
)
ROW FORMAT SERDE 'org.apache.hive.hcatalog.data.JsonSerDe'
LOCATION 's3://${log_bucket_path}/'
TBLPROPERTIES (
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

data "aws_s3_bucket" "apigateway_log_bucket" {
  bucket = "apigateway-log-bucket"
}

data "template_file" "create_apigateway_log_table_sql" {
  template = file("./create_apigateway_log_table.sql.tpl")
  vars = {
    athena_database_name = resource.aws_athena_database.athena_database.name
    athena_table_name         = "apigateway_logs"
    log_bucket_path               = "${data.aws_s3_bucket.apigateway_log_bucket.id}/test-apigateway/data"
  }
}

resource "aws_athena_named_query" "athena_named_query" {
  name        = "create_apigateway_log_table"
  database = resource.aws_athena_database.athena_database.id
  query       = data.template_file.create_apigateway_log_table_sql.rendered
}
```

上記のTerraformファイルを実行するとAthenaの「保存したクエリ」に `create_apigateway_log_table` ができていますので、こちらを実行してテーブル作成します。

## テーブルに対してクエリしてみる

最後にきちんと設定できているかどうか確認するために実際にクエリを発行してみましょう。

```
SELECT * FROM "apigateway_athena_db"."apigateway_logs" limit 10;
```

表示されれば正しくテーブルの設定ができています。

あとは解析したいクエリをこちらのテーブルに対して実行するだけになります。

## まとめ

今回は、API GatewayのログをAthenaでクエリする方法について解説しました。

API Gatewayはアクセスログを直接S3にアウトプットできません。その対応のためにKinesis Data Firehoseの設定と、ログに改行コードを埋め込むためにLambdaの設定を行っています。

今までAPI Gatewayのログを活用できていなかったという方がいるのではないでしょうか？こちらのログを分析に使用することで多くの利点があります。

是非今回の設定を試して、いろいろな解析に利用してみてください。

今回の記事がAWSを使用する方にとって少しでも参考になりましたら幸いです。