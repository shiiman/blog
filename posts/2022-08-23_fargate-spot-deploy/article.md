---
id: 1736
title: 【AWS】Fargate Spotを設定した環境でBlue-Greenデプロイした時に躓いた話
slug: fargate-spot-deploy
status: publish
date: 2022-08-23T19:30:00
modified: 2022-08-23T00:30:38
excerpt: Fargate Spot環境でBlue-Greenデプロイを行う際に遭遇した問題と、CodeDeployの設定での解決方法を紹介します。
categories: [19, 18]
tags: [115, 121, 122]
featured_media: 1738
---

こんばんは、しーまんです。

先日 **Fargate Spot** を導入してインフラ費用の削減を行った話を投稿しました。

その際にBlue-Greenデプロイで少し躓いた点がありましたので、今回はその話をまとめてみました。

もしFargate Spotを導入検討されている方や、デプロイ周りで躓いている方は参考にしてみてください！

## Fargate Spot

まずはFargate Spotについて簡単におさらいしておきましょう。

Fargate SpotとはAWSが提供するマネージドコンテナサービスであるFargateを、余剰リソースを使用することで通常より安く使用できるサービスです。

コンテナ費用を削減したいという方が検討するべきサービスになります。

こちらの導入方法については前回の記事を参考にして導入してみてください。

[![](https://shiimanblog.com/wp-content/uploads/2022/05/eyecatch_fargatespot-320x180.jpg)\
\
【AWS】Fargate Spotを導入してコンテナ費用を7割削減した話\
\
AWSでコンテナを扱う選択肢として最近良く聞くFargateというサービスがありますが、そのFargateのスポットインスタンスのようなものが登場していたのでそちらを適応し、コンテナ費用を7割削減した話を解説致します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2022.05.13](https://shiimanblog.com/engineering/post-1651/ "【AWS】Fargate Spotを導入してコンテナ費用を7割削減した話")

## 導入時に躓いた点

今回はFargate Spot導入時に躓いた点を話したいと思います。

導入自体はとても簡単でそんなに苦労しなかったのですが、デプロイ周りで予想外の挙動が起こりました。今回アプリをデプロイする時に [Code Deploy](https://aws.amazon.com/jp/codedeploy/) による [Blue Greenデプロイ](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/deployment-type-bluegreen.html) を使用していました。

Fargate Spotを導入するにはECS Serviceに [キャパシティープロバイダー](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-capacity-providers.html) という設定を入れます。

しかしここが今回ハマった点で、デプロイを行うと折角設定したFargate Spotのキャパシティープロバイダー設定が **消えてしまい** ました。

でも大丈夫です。

Fargate Spotのリリース当初はCode DeployによるBlue Greenデプロイは未対応だったらしいですが、こちらはもう対応されています。

あまり情報が無かったですが、私の方で試してみたところ問題なく動作することが確認できます。

ですので、次項目の対応を行ってみてください。

## Blue-Greenデプロイの対応方法

Fargate Spotの設定であるキャパシティープロバイダー設定をCode Deployを使用したデプロイ時に維持できるようにしたいです。

そこで [Code Deployの設定ファイル](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/reference-appspec-file.html) であるappspec.ymlを編集します。

変更前のAppSpecファイルの例は以下です。

```
version: 0.0
Resources:
    - TargetService:
          Type: AWS::ECS::Service
          Properties:
              TaskDefinition: "<TASK_DEFINITION>"
              LoadBalancerInfo:
                  ContainerName: "web"
                  ContainerPort: 80
```

プロパティにタスク定義とロードバランサーの設定のみを記載しておりました。

続いて変更後の例は以下になります。

```
version: 0.0
Resources:
    - TargetService:
          Type: AWS::ECS::Service
          Properties:
              TaskDefinition: "<TASK_DEFINITION>"
              LoadBalancerInfo:
                  ContainerName: "web"
                  ContainerPort: 80
              CapacityProviderStrategy:
                  - CapacityProvider: "FARGATE"
                    Weight: 1
                    Base: 0
                  - CapacityProvider: "FARGATE_SPOT"
                    Weight: 2
                    Base: 0
```

プロパティに [キャパシティプロバイダーの設定](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/reference-appspec-file-structure-resources.html#reference-appspec-file-structure-resources-ecs) が追加されているのが分かると思います。

ここにキャパシティプロバイダーの設定を入れることでCode Deploy時にサービスのキャパシティプロバイダーを反映させることができます。

## まとめ

今回は、Fargate Spot 導入時にデプロイで躓いた点について解説しました。

Fargate Spotの設定であるキャパシティープロバイダーの設定がCode Deployのタイミングで消えてしまう問題です。こちらの問題はCode Deployの設定であるAppSpecにキャパシティープロバイダーの設定を追加することで解消されます。

なかなか検索しても見つからなかったので、私と同じようにハマった人もいたのではないでしょうか。

同じような事象に陥って困っている方の参考に少しでもなりましたら幸いです。