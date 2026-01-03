---
id: 1521
title: 【AWS】Firelensを使用してS3にログをリアルタイムで出力する方法 | 【Fluentbit】【Kinesis Data Firehose】
slug: firelens-s3
status: publish
excerpt: こんばんは、しーまんです。 今回は久しぶりにAWSネタです。一番業務で使用しているのに逆に記事にすることが少なかったことに気づきました。これからはAWSの記事頻度を上げられればと思います。 皆さんはアプリケーションのログ \[…\]
categories:
    - 19
    - 18
tags:
    - 98
    - 99
    - 100
featured_media: 1524
date: 2022-01-24T19:30:00
modified: 2022-01-21T17:31:34
---

こんばんは、しーまんです。

今回は久しぶりにAWSネタです。一番業務で使用しているのに逆に記事にすることが少なかったことに気づきました。これからはAWSの記事頻度を上げられればと思います。

皆さんはアプリケーションのログをどのように扱っているでしょうか？AWSを使用する場合とりあえず [S3](https://aws.amazon.com/jp/s3/) に配置してそこから [Athena](https://aws.amazon.com/jp/athena) なり [Redshift](https://aws.amazon.com/jp/redshift/) なりで分析しているのではないでしょうか！

そこで今回はECSで動作しているアプリケーションのログを**リアルタイムでS3に出力する方法**を紹介致します。ログの扱いをどのようにしようかと頭を悩ませている方には参考になると思いますので、ぜひご覧になってください。

## 前提

今回は [ECS](https://aws.amazon.com/jp/ecs/) のコンテナ上で動作しているアプリケーションのログをS3に出力することを前提としています。ただし、他の環境で動作しているアプリケーションのログに関しても出力する場所が変わるだけで、その後の動作は参考になると思います。

また今回配置するS3のバケット名は ｢ **firelens-test-bucket**」とします。こちらはご自身の環境に合わせて随時読み変えてください。

## 全体設計

今回構築するAWSの全体的な設定は下図のようになっています。

[![システム構成図](https://shiimanblog.com/wp-content/uploads/2022/01/firelens_diagram-1-800x287.png)](https://shiimanblog.com/wp-content/uploads/2022/01/firelens_diagram-1.png)

ログを扱う時は以下の点に注意して設計するようにしましょう。

- 出力されるログの形式は決まっているか(決まっていない場合はログ定義から始めましょう)
- 更新頻度はどの程度か(リアルタイム性の有無)
- 出力されるログはどのくらいの量か
- ログの欠損が許容されるかどうか(許容される場合は許容値の目安を決めておきましょう)

## 今回使用するリソース

今回使用するリソースは全体設計の図でも記載がありますが、「 **Firelens**」「 **Fluentbit**」「 **Firehose**」を主に使用します。こちらのリソースは名前が似ていますが、それぞれの役割をちゃんと把握して使用できるようになりましょう。

またこれは私が勝手に覚えるために自分で作った名称ですが、AWSでログを扱う時はFから始まるリソースを3つ使用するので、AWSのS3のように「 **F3**」と読んでいます。

皆さんも**AWSのログ = 「F3」**と覚えてしまいましょう。

### Firelens

それでは、それぞれのリソースの説明と設定方法を解説していきます。

まずは「 **Firelens**」からみていきましょう。

FirelensとはECSで使用できるログルーターのことです。タスク定義の中に含めてサイドカーとして配置しつつ、他のコンテナからはログドライバーとして使用します。

例としてnginxのログをfluentbitに流すタスク定義を紹介します。

※ <>で囲んでいる項目はご自身の環境に合わせて置き換えてください。

```
{
    "family": "<FAMILY>",
    "taskRoleArn": "<TASK_ROLE_ARN>",
    "executionRoleArn": "<EXECUTION_ROLE_ARN>",
    "networkMode": "awsvpc",
    "requiresCompatibilities": [
        "FARGATE"
    ],
    "cpu": "<CPU>",
    "memory": "<MEMORY>",
    "containerDefinitions": [
        {
            "name": "web",
            "image": "<NGINX_IMAGE>",
            "logConfiguration": {
                "logDriver": "awsfirelens"
            },
            "portMappings": [
                {
                    "containerPort": 80,
                    "hostPort": 80,
                    "protocol": "tcp"
                }
            ],
            "essential": true
        },
        {
            "name": "firelens",
            "image": "<FLUENTBIT_IMAGE>",
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "/ecs/fluentbit",
                    "awslogs-region": "ap-northeast-1",
                    "awslogs-stream-prefix": "<ENVIRONMENT>"
                }
            },
            "firelensConfiguration": {
                "type": "fluentbit",
                "options": {
                    "enable-ecs-log-metadata": "true",
                    "config-file-type": "file",
                    "config-file-value": "/fluent-bit/etc/fluent-bit-custom.conf"
                }
            },
            "essential": true
        }
    ]
}
```

タスク定義は解説しだすとキリがないので大切な部分だけピックアップして解説します。

#### ログを送る側のコンテナ設定

まずはログを送る側のコンテナ設定です。

下記で定義しているように **logDriver** に awsfirelens を指定します。

```
"logConfiguration": {
    "logDriver": "awsfirelens"
},
```

この指定をすることで標準出力されたログをサイドカーとして起動している **firelensコンテナ** に送ることが出来ます。

#### ログを受け取る側のコンテナ設定

次にログを受け取る側のコンテナ設定です。

下記で定義している **firelensConfiguration** の設定を行います。

```
"firelensConfiguration": {
    "type": "fluentbit",
    "options": {
        "enable-ecs-log-metadata": "true",
        "config-file-type": "file",
        "config-file-value": "/fluent-bit/etc/fluent-bit-custom.conf"
    }
},
```

受け取る側のコンテナでは **fluentbit** を使用するのがスタンダードなので、そちらの設定を保存したイメージを作成して指定してあげましょう。

最初は [ドキュメント](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/userguide/firelens-taskdef.html#firelens-taskdef-customconfig) の記載の通り **config-file-type** を **s3** として設定ファイルをS3から取得するように設定したのですが、うまくいきませんでした。

再度 [ドキュメント](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/userguide/firelens-taskdef.html#firelens-taskdef-customconfig) をよーく読み返したところ下記の1文がありました。

AWS Fargate でホストされるタスクは、file 設定ファイルタイプのみをサポートします。

私はこちらの注意書きを見逃していてエラーにハマっていました。

早く [Fargate](https://aws.amazon.com/jp/fargate/) でも対応してくれたら嬉しいですね。

ということで [Fargate](https://aws.amazon.com/jp/fargate/) を使用している皆さんはご注意ください!!

Fargateを使用する場合は**S3からFluentbitのconfig fileを読み込むことが出来ません**。

仕方がないので、Fargateを使用する場合は configファイルを配置するためだけに下記のような **dockerfile** を作成して、それを使ったイメージを指定するようにします。

```
FROM amazon/aws-for-fluent-bit:2.21.5
COPY ./fluent-bit-custom.conf /fluent-bit/etc/fluent-bit-custom.conf
```

こちらの設定で、nginxのログをfirelensコンテナに流し、fluentbitの処理でログをフィルタリングしたり、送信先を指定できるようになりました。

### Fluentbit

次はFluentbitの設定になります。

上記で配置した **fluent-bit-custom.conf** ファイルに様々な処理を定義していきます。

今回やることはnginxのログを firehoseに送る設定ですね。

下記に例を上げておきます。

```
[SERVICE]
    Flush        1
    Grace        30
    Log_Level    info
    Parsers_File /fluent-bit/parsers/parsers.conf

[FILTER]
    Name         parser
    Match        web-firelens-*
    Key_Name     log
    Parser       json
    Reserve_Data True

[OUTPUT]
    Name            firehose
    Match           web-firelens-*
    region          ap-northeast-1
    delivery_stream <KINESIS_DATA_FIREHOSE_NAME>
```

Firelensを使用するとログを特定する場合 「 **\[コンテナ名\]-firelens-**」 で送られてきますので、今回の場合 「 **web-firelens-\***」 を指定しています。

ログを出力するコンテナが複数ある場合はこちらでフィルタリングしてあげるとよいでしょう。(今回は1コンテナなので 「 **\***」 を指定しても問題ありません)

そして **\[OUTPUT\]** 項目で出力先を **firehose** に指定しています。その他のアウトプット先がどのようなものがあるのかを知りたい場合 [ドキュメント](https://docs.fluentbit.io/manual/pipeline/outputs) を参考にしてみてください。

### Kinesis Data Firehose

最後にKinesis Data Firehoseの設定になります。

こちらは先程「 **fluent-bit-custom.conf**」で設定した **delivery\_stream** の値と同じ名前で作成します。

AWSコンソールから「 **Create delivery stream**」を押します。

[![firehose - 設定1](https://shiimanblog.com/wp-content/uploads/2022/01/firehose_setting1-800x178.png)](https://shiimanblog.com/wp-content/uploads/2022/01/firehose_setting1.png)

するとdelivery stream作成画面に遷移しますので、下記を設定します。

「 **Source**」に Direct PUT

「 **Destination**」に Amazon S3

「 **Delivery stream name**」に 先程delivery\_streamに設定した名前

[![firehose - 設定2](https://shiimanblog.com/wp-content/uploads/2022/01/firehose_setting2-800x649.png)](https://shiimanblog.com/wp-content/uploads/2022/01/firehose_setting2.png)

そして最後に出力先バケットを指定します。

最初に定義したとおり、今回は「 **firelens-test-bucket**」を指定します。

[![firehose - 設定3](https://shiimanblog.com/wp-content/uploads/2022/01/firehose_setting3-800x193.png)](https://shiimanblog.com/wp-content/uploads/2022/01/firehose_setting3.png)

S3に出力する場合、ディレクトリを指定したり、圧縮設定をしたりといったことが可能ですので、必要な場合は利用してみましょう！

ここまでで全ての設定が完了になります。

あとは実際にnginxにアクセスをしてログがS3に配置されることを確認しましょう。

Firehoseのデフォルト設定では「**Buffer sizeが5MiB**」「 **Buffer intervalが300secound**」となっているためそれらの既定値が超えたタイミングでS3にファイル出力されます。

アクセスして直ぐにS3に書き込まれるわけではないので注意しましょう。

負荷やリアルタイム性を考慮して上記の設定を調整してください。

## まとめ

今回はAWSでログをリアルタイムにS3へ配置する方法として「 **F3**」を使用する方法を解説しました。F3とは「 **Firelens**」「 **Fluentbit**」「 **Firehose**」のことで、近年コンテナでアプリケーションを動かすことが増えた環境ではスタンダードな方法といえます。

用語は紛らわしいですが、一つ一つ役割を理解すれば難しいことはありません。

振り分けルールを細かく設定したい場合 **Fluentbit** の設定を理解する必要がありますが、Fluentbitにはしっかりと [ドキュメント](https://docs.fluentbit.io/manual/) が用意されていますので、そこまで苦労しないでしょう。

AWSで分析基盤を作成している方や、ログについて手軽に設定したい方はぜひ参考にしてみてください。今回の記事が少しでも、お役に立てれば幸いです。