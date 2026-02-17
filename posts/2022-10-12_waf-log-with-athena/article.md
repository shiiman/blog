---
id: 1782
title: 【AWS】WAFのログをAthenaでクエリする方法
slug: waf-log-with-athena
status: publish
date: 2022-10-12T19:30:00
modified: 2022-10-11T12:06:53
excerpt: S3に保存されたAWS WAFのログをAthenaでクエリする方法を、テーブル作成からクエリ実行まで解説します。
categories: [19, 18, 21]
tags: [83, 124, 125]
featured_media: 1784
---

こんばんは、しーまんです。

前回の記事でAPI GatewayのアクセスログをAthenaでクエリする方法を紹介しました。

まだご覧になっていない方はぜひこちらも御覧ください。

 [![](https://shiimanblog.com/wp-content/uploads/2022/10/eyecatch_api_gateway_log_with_athena-320x180.jpg)\
\
【AWS】API GatewayのログをAthenaでクエリする方法\
\
AWSでAPI Gatewayのアクセスログをs3に保存し、Athenaでクエリする方法を解説致します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2022.10.11](https://shiimanblog.com/engineering/api-gateway-log-with-athena/ "【AWS】API GatewayのログをAthenaでクエリする方法")

さて今回はAthenaでログをクエリする紹介第4弾になります。前回はAPI Gatewayのログでしたが、今回は [WAF](https://aws.amazon.com/jp/waf/) のアクセスログを [Athena](https://aws.amazon.com/jp/athena) を使用してクエリしていこうと思います。

AWSでのログ運用で困っている方の参考になればいいなと思います。

## 事前準備

まず前提としてWAFのアクセスログをS3に配置済みの状況とします。

もしまだWAFのログ出力を有効にしていない場合は有効にしておきましょう。

参考までに [Terraform](https://www.terraform.io/) でのログの有効化設定の例を載せておきます。

```
# wafのログを保存するバケットの取得.
data "aws_s3_bucket" "waf_log_bucket" {
  bucket = "aws-waf-logs-bucket"
}

# wafの作成.
resource "aws_wafv2_web_acl" "wafv2_web_acl" {
    name  = "test-waf"
    scope = "CLOUDFRONT"

    ....省略
}

# wafのログ設定.
resource "aws_wafv2_web_acl_logging_configuration" "waf_logging_configuration" {
  log_destination_configs = data.aws_s3_bucket.waf_log_bucket.arn
  resource_arn                    = resource.aws_wafv2_web_acl.wafv2_web_acl.arn
}
```

WAFのログを保存するS3バケット名はプレフィックスとして「 **aws-waf-logs-**」をつける必要があります。

上記の設定だとWAFのアクセスログは以下のパスに保存されます。

```
s3://aws-waf-logs-bucket/AWSLogs/{account_id}/WAFLogs/{region}/{webacl}/yyyy/MM/dd/HH/mm/xxxx.log.gz
```

{account\_id}はご自身のAWSアカウントのID、{region}にはWAFの配置されているリージョン、{webacl}にはWAFの名前が入ります。

Cloud Frontに設定しているWAFの場合regionは「 **cloudfront**」になります。

## Athena

さて、出力されているログの情報が分かりましたので、ここからAthenaでクエリできるように設定していきます。

AthenaとはサーバーレスでS3内のデータに対して標準SQLを使用して簡単に分析できるサービスです。

こちらを使用するためにはまずデータベースを作成し、次にテーブルを作成します。

作成したテーブルに対してクエリを発行し分析を行うという形ですね。

ですのでまずはデータベースを作成していきます。

Terraformで作成する場合は下記のように `name` と `bucket` を指定することで作成可能です。

```
resource "aws_athena_database" "athena_database" {
  name   = "waf_athena_db"
  bucket = "athena_query_log_bucket"
}
```

bucketはWAFログの配置されているbucketとは別物です。

athena自身のクエリログを出力するためのbucketになります。

## テーブル作成

データベースの作成が終わったら、いよいよクエリを行うためのテーブルを作成していきます。

Athenaは [Glueデータカタログ](https://aws.amazon.com/jp/glue/) と統合されているためそちらのリソースで定義してもよいのですが、 [Create SQLで定義](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/waf-logs.html) しておいた方が分かりやすいため、SQLで用意します。

またクエリを発行する際にスキャンデータ量を抑えるためにパーティションを作成することが推奨されています。以前は手動で作成したり、Lambdaを利用して定期的に作成する必要がありました。しかし2020年のアップデートで [パーティション射影](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/partition-projection.html) という機能がリリースされ、自動でパーティションを追加できる仕組みが追加されました。

今回用意するSQLもパーティション射影の設定を追加しています。

サンプルとして下記の **create\_waf\_log\_table.sql.tpl** を用意します。

```
CREATE EXTERNAL TABLE IF NOT EXISTS ${athena_database_name}.${athena_table_name} (
    `timestamp` bigint,
    `formatversion` int,
    `webaclid` string,
    `terminatingruleid` string,
    `terminatingruletype` string,
    `action` string,
    `terminatingrulematchdetails` array<
                                        struct<
                                            conditiontype:string,
                                            location:string,
                                            matcheddata:array<string>
                                            >
                                        >,
    `httpsourcename` string,
    `httpsourceid` string,
    `rulegrouplist` array<
                        struct<
                            rulegroupid:string,
                            terminatingrule:struct<
                                                ruleid:string,
                                                action:string,
                                                rulematchdetails:string
                                                    >,
                            nonterminatingmatchingrules:array<string>,
                            excludedrules:string
                                >
                        >,
    `ratebasedrulelist` array<
                            struct<
                                ratebasedruleid:string,
                                limitkey:string,
                                maxrateallowed:int
                                >
                            >,
    `nonterminatingmatchingrules` array<
                                        struct<
                                            ruleid:string,
                                            action:string
                                            >
                                        >,
    `requestheadersinserted` string,
    `responsecodesent` string,
    `httprequest` struct<
                        clientip:string,
                        country:string,
                        headers:array<
                                    struct<
                                        name:string,
                                        value:string
                                        >
                                    >,
                        uri:string,
                        args:string,
                        httpversion:string,
                        httpmethod:string,
                        requestid:string
                        >,
    `labels` array<
                struct<
                    name:string
                        >
                    >,
    `captcharesponse` struct<
                            responsecode:string,
                            solvetimestamp:string,
                            failureReason:string
                            >
)
PARTITIONED BY (
    `day` string
)
ROW FORMAT SERDE
    'org.openx.data.jsonserde.JsonSerDe'
STORED AS INPUTFORMAT
    'org.apache.hadoop.mapred.TextInputFormat'
OUTPUTFORMAT
    'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
    's3://${log_bucket_path}/AWSLogs/${aws_account_id}/WAFLogs/${region}/${webacl}/'
TBLPROPERTIES (
 'projection.enabled' = 'true',
 'projection.day.type' = 'date',
 'projection.day.range' = '2021/01/01,NOW',
 'projection.day.format' = 'yyyy/MM/dd',
 'projection.day.interval' = '1',
 'projection.day.interval.unit' = 'DAYS',
 'storage.location.template' = 's3://${log_bucket_path}/AWSLogs/${aws_account_id}/WAFLogs/${region}/${webacl}/\$\{day\}'
)

```

こちらのSQLテンプレートをTerraformで反映していきます。

```
data "aws_caller_identity" "current_account_id" {}

data "aws_s3_bucket" "waf_log_bucket" {
  bucket = "aws-waf-logs-bucket"
}

data "template_file" "create_waf_log_table_sql" {
  template = file("./create_waf_log_table.sql.tpl")
  vars = {
    athena_database_name = resource.aws_athena_database.athena_database.name
    athena_table_name         = "waf_logs"
    log_bucket_path               = data.aws_s3_bucket.waf_log_bucket.id
    aws_account_id                = data.aws_caller_identity.current_account_id.account_id
    region                                  = "cloudfront"
    webacl                                 = resource.aws_wafv2_web_acl.wafv2_web_acl.name
  }
}

resource "aws_athena_named_query" "athena_named_query" {
  name        = "create_waf_log_table"
  database = resource.aws_athena_database.athena_database.id
  query       = data.template_file.create_waf_log_table_sql.rendered
}
```

上記のTerraformファイルを実行するとAthenaの「保存したクエリ」に `create_waf_log_table` ができていますので、こちらを実行してテーブル作成します。

## テーブルに対してクエリしてみる

最後にきちんと設定できているかどうか確認するために実際にクエリを発行してみましょう。

```
SELECT * FROM "waf_athena_db"."waf_logs" limit 10;
```

表示されれば正しくテーブルの設定ができています。

あとは解析したいクエリをこちらのテーブルに対して実行するだけになります。

クエリサンプルに関しては [公式ドキュメント](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/waf-logs.html) に少し例がありますので、参考にしてください。

## まとめ

今回は、WAFのログをAthenaでクエリする方法について解説しました。

今までWAFのログをS3には配置していたけど、出力していただけで利用していなかったという方がいるのではないでしょうか？こちらのログを分析に使用することで多くの利点があります。

是非今回の設定を試して、いろいろな解析に利用してみてください。

今回の記事がAWSを使用する方にとって少しでも参考になりましたら幸いです。