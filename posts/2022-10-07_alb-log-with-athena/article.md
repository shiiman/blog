---
id: 1748
title: 【AWS】Application LoadbalancerのログをAthenaでクエリする方法
slug: alb-log-with-athena
status: publish
date: 2022-10-07T19:30:00
modified: 2022-10-07T17:57:47
excerpt: S3に保存されたALBのアクセスログをAthenaでクエリする方法を、テーブル作成からクエリ実行まで解説します。
categories: [19, 18, 21]
tags: [123, 124, 125]
featured_media: 1749
---

こんばんは、しーまんです。

AWSではいろいろなリソースのログを [S3](https://aws.amazon.com/jp/s3) に保存することができます。

しかし、そのままでは保存したログを見ることができません。

今回はS3に保存された [Application Loadbalancer](https://aws.amazon.com/jp/elasticloadbalancing/application-load-balancer/) のログを [Athena](https://aws.amazon.com/jp/athena) を使用してクエリすることで、ログを検索・調査できるようにします。

AWSでのログ運用で困っている方の参考になればいいなと思います。

## 事前準備

まず前提としてApplication Aladbalancer(ALB)のアクセスログをS3に配置済みの状況とします。

もしまだALBのログ出力を有効にしていない場合は有効にしておきましょう。

参考までに [Terraform](https://www.terraform.io/) でのログの有効化設定の例を載せておきます。

```
resource "aws_lb" "alb" {
  name  = "test-alb"
  ....省略

  access_logs {
    bucket   = "alb-log-bucket"
    prefix     = "test-alb"
    enabled = true
  }
}
```

上記の設定だとALBのアクセスログは以下のパスに保存されます。

```
s3://alb-log-bucket/test-alb/AWSLogs/{account_id}/elasticloadbalancing/{region}/yyyy/MM/dd/xxxx.log.gz
```

{account\_id}はご自身のAWSアカウントのID、{region}にはALBの配置されているリージョンが入ります。

## Athena

さて、出力されているログの情報が分かりましたので、ここからAthenaでクエリできるように設定していきます。

AthenaとはサーバーレスでS3内のデータに対して標準SQLを使用して簡単に分析できるサービスです。

こちらを使用するためにはまずデータベースを作成し、次にテーブルを作成します。

作成したテーブルに対してクエリを発行し分析を行うという形ですね。

ですのでまずはデータベースを作成していきます。

Terraformで作成する場合は下記のように `name` と `bucket` を指定することで作成可能です。

```
resource "aws_athena_database" "athena_database" {
  name   = "alb_athena_db"
  bucket = "athena_query_log_bucket"
}
```

bucketはALBログの配置されているbucketとは別物です。

athena自身のクエリログを出力するためのbucketになります。

## テーブル作成

データベースの作成が終わったら、いよいよクエリを行うためのテーブルを作成していきます。

Athenaは [Glueデータカタログ](https://aws.amazon.com/jp/glue/) と統合されているためそちらのリソースで定義してもよいのですが、 [Create SQLで定義](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/application-load-balancer-logs.html) しておいた方が分かりやすいため、SQLで用意します。

またクエリを発行する際にスキャンデータ量を抑えるためにパーティションを作成することが推奨されています。以前は手動で作成したり、Lambdaを利用して定期的に作成する必要がありました。しかし2020年のアップデートで [パーティション射影](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/partition-projection.html) という機能がリリースされ、自動でパーティションを追加できる仕組みが追加されました。

今回用意するSQLもパーティション射影の設定を追加しています。

サンプルとして下記の **create\_alb\_log\_table.sql.tpl** を用意します。

```
CREATE EXTERNAL TABLE IF NOT EXISTS ${athena_database_name}.${athena_table_name} (
    `type` string,
    `time` string,
    `elb` string,
    `client_ip` string,
    `client_port` int,
    `target_ip` string,
    `target_port` int,
    `request_processing_time` double,
    `target_processing_time` double,
    `response_processing_time` double,
    `elb_status_code` int,
    `target_status_code` string,
    `received_bytes` bigint,
    `sent_bytes` bigint,
    `request_verb` string,
    `request_url` string,
    `request_proto` string,
    `user_agent` string,
    `ssl_cipher` string,
    `ssl_protocol` string,
    `target_group_arn` string,
    `trace_id` string,
    `domain_name` string,
    `chosen_cert_arn` string,
    `matched_rule_priority` string,
    `request_creation_time` string,
    `actions_executed` string,
    `redirect_url` string,
    `lambda_error_reason` string,
    `target_port_list` string,
    `target_status_code_list` string,
    `classification` string,
    `classification_reason` string
)
PARTITIONED BY
(
    `day` STRING
)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.RegexSerDe'
WITH SERDEPROPERTIES (
'serialization.format' = '1',
'input.regex' =
'([^ ]*) ([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*)[:-]([0-9]*) ([-.0-9]*) ([-.0-9]*) ([-.0-9]*) (|[-0-9]*) (-|[-0-9]*) ([-0-9]*) ([-0-9]*) \"([^ ]*) (.*) (- |[^ ]*)\" \"([^\"]*)\" ([A-Z0-9-_]+) ([A-Za-z0-9.-]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^\"]*)\" ([-.0-9]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^ ]*)\" \"([^\s]+?)\" \"([^\s]+)\" \"([^ ]*)\" \"([^ ]*)\"')
LOCATION 's3://${log_bucket_path}/AWSLogs/${aws_account_id}/elasticloadbalancing/${region}/'
TBLPROPERTIES
(
    "projection.enabled" = "true",
    "projection.day.type" = "date",
    "projection.day.range" = "2022/01/01,NOW",
    "projection.day.format" = "yyyy/MM/dd",
    "projection.day.interval" = "1",
    "projection.day.interval.unit" = "DAYS",
    "storage.location.template" = "s3://${log_bucket_path}/AWSLogs/${aws_account_id}/elasticloadbalancing/${region}/\$\{day\}"
)
```

こちらのSQLテンプレートをTerraformで反映していきます。

```
data "aws_caller_identity" "current_account_id" {}

data "aws_s3_bucket" "alb_log_bucket" {
  bucket = "alb-log-bucket"
}

data "template_file" "create_alb_log_table_sql" {
  template = file("./create_alb_log_table.sql.tpl")
  vars = {
    athena_database_name = resource.aws_athena_database.athena_database.name
    athena_table_name         = "alb_logs"
    log_bucket_path               = "${data.aws_s3_bucket.alb_log_bucket.id}/test-alb"
    aws_account_id                = data.aws_caller_identity.current_account_id.account_id
    region                                 = "ap-northeast-1"
  }
}

resource "aws_athena_named_query" "athena_named_query" {
  name        = "create_alb_log_table"
  database = resource.aws_athena_database.athena_database.id
  query       = data.template_file.create_alb_log_table_sql.rendered
}
```

上記のTerraformファイルを実行するとAthenaの「保存したクエリ」に `create_alb_log_table` ができていますので、こちらを実行してテーブル作成します。

## テーブルに対してクエリしてみる

最後にきちんと設定できているかどうか確認するために実際にクエリを発行してみましょう。

```
SELECT * FROM "alb_athena_db"."alb_logs" limit 10;
```

表示されれば正しくテーブルの設定ができています。

あとは解析したいクエリをこちらのテーブルに対して実行するだけになります。

クエリサンプルに関しては [公式ドキュメント](https://docs.aws.amazon.com/ja_jp/athena/latest/ug/application-load-balancer-logs.html) に少し例がありますので、参考にしてください。

## まとめ

今回は、Application LoadbalancerのログをAthenaでクエリする方法について解説しました。

今までALBのログをS3には配置していたけど、出力していただけで利用していなかったという方がいるのではないでしょうか？こちらのログを分析に使用することで多くの利点があります。

是非今回の設定を試して、いろいろな解析に利用してみてください。

今回の記事がAWSを使用する方にとって少しでも参考になりましたら幸いです。