---
id: 1651
title: 【AWS】Fargate Spotを導入してコンテナ費用を7割削減した話
slug: post-1651
status: publish
excerpt: こんばんは、しーまんです。 今回は最近 Fargate Spot を導入してインフラ費用の削減を行いましたので、そのお話。FargateとはAWSが提供するコンテナ向けサーバーレスコンピューティングのことで、そのFarg \[…\]
categories:
    - 19
    - 18
tags:
    - 115
featured_media: 1644
date: 2022-05-13T19:30:00
modified: 2022-10-28T17:58:29
---

こんばんは、しーまんです。

今回は最近 **Fargate Spot** を導入してインフラ費用の削減を行いましたので、そのお話。

FargateとはAWSが提供するコンテナ向けサーバーレスコンピューティングのことで、そのFargateの費用最適化をしたということです。

もし [Fargate](https://aws.amazon.com/jp/fargate/) を既に使用しているけどSpotの導入をまだしていないという方は、是非この記事を参考にして導入をしてみてください。かなり効果が高い費用削減が行えると思います！

## Fargate Spotとは

まずは [Fargate](https://aws.amazon.com/jp/fargate/) とは何か、Fargate Spot とは何かを簡単に確認していきましょう。

Fargateとは冒頭でも述べたようにAWSが提供するコンテナ向けサーバーレスコンピューティングのことです。AWSで提供するコンテナサービスと言えば [ECS](https://aws.amazon.com/jp/ecs/) と [EKS](https://aws.amazon.com/jp/eks/) が思い浮かぶと思いますが、そのどちらでもFargateを使用することが可能です。

以前は ECSやEKSのバックグラウンドとしてEC2インスタンスを自身で管理する必要がありました。

Fargateはそのバックグラウンドのインスタンス管理をマネージド化したサービスです。

### Fargate Spotとは

そしてFargateの機能の一つとして AWS re:Invent2019 で発表されたのが、Fargate Spot です。

Fargate Spot とは [EC2のスポットインスタンス](https://aws.amazon.com/jp/ec2/spot/) と同様、AWSで管理している余剰リソースをいつ落ちても問題ない代わりに通常より安く使用できるというものです。

いつ落ちるか分からないということですが、Fargate Spot と通常のFargateは共存できます。

つまり 「Fargate Spot が落ちた際は通常のFargate が自動で立ち上がる」という設定も可能です！

とはいえ現状だとほぼ落ちることはない上に、万が一落ちる場合でもちゃんと通知を受け取ることが可能です。落ちても問題ない開発環境や、コストをできるだけ抑えたい環境で使用することがおすすめですね。

### 費用

ということで Fargate Spot の費用をみていきましょう！

Fargateの価格は、以下のページにまとまっています。

[AWS Fargate の料金](https://aws.amazon.com/jp/fargate/pricing/?nc=sn&loc=2)

**東京リージョン** における、オンデマンド(通常のFargate)とSpotの料金を比較したのが下記表です。

種別オンデマンドSpotオンデマンドに対するSpotの料金比率per vCPU per hour$0.050560$0.01516830%per GB per hour$0.005530$0.00165930%

つまり通常の7割り引きで使用できるということですね。

めっちゃ安くなります！

## 導入方法

料金が安くなり、落ちても問題ない設定が容易になので導入しない手はないですよね。

ということで早速導入方法をみていきましょう！

### Terraformで導入

ネットだとコマンドで導入している例とかが見受けられるのですが、私は [Terraform](https://www.terraform.io/) ユーザですので、Terraformで導入していきます。

#### キャパシティープロバイダー

Fargate Spotを導入するには [キャパシティープロバイダー](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-capacity-providers.html) という考え方必要です。

キャパシティープロバイダー(以降 CP)とは「 **FARGATE**」と「 **FARGATE\_SPOT**」を何対何の割合で起動するかを定義するものだと考えると理解が簡単です。

CPの設定には 「 **プロバイダー**」「 **ベース**」「 **ウェイト**」の3項目あり、例えば下記のように設定すると、最初立ち上がるタスクはbaseで設定された「FARGATE」が必ず立ち上がります。その後立ち上がるタスクについては 「FARGATE」：「FARGATE\_SPOT」 = 2：1 の割合で立ち上がります。

```
{
      capacity_provider = "FARGATE"
      base              = 1
      weight           = 2
 },
 {
      capacity_provider = "FARGATE_SPOT"
      base              = 0
      weight           =1
 }
```

このあたりの挙動は実際調整してみて確認するのが早いので、いろいろと弄って確かめてください。

### ECSにデフォルトキャパシティープロバイダーを設定

ということで、まずはCPをECSに設定していきます。

こちらで設定した値がデフォルトで使用される値になります。

```
/**
 * ECS Capacity Providers作成
 * https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_cluster_capacity_providers
 */
resource "aws_ecs_cluster_capacity_providers" "ecs_cluster_capacity_providers" {
    cluster_name       = "test_cluster"
    capacity_providers = ["FARGATE", "FARGATE_SPOT"]

    default_capacity_provider_strategy {
        capacity_provider = "FARGATE"
        base              = 0
        weight            = 1
     }

    default_capacity_provider_strategy {
        capacity_provider = "FARGATE_SPOT"
        base              = 0
        weight           = 2
     }
}
```

こんな感じで設定しておくと基本的には「Fargate Spot」でタスクが立ち上がり、もし落ちたとしても「Fargate」が立ち上がってくれます。

baseはどちらか1つしか有効に出来ません。

両方とも無効でも問題なく動作します。

### Service毎にキャパシティープロバイダーを変更したい場合

基本はクラスターのCPに乗っ取るので良いですが、作成するService毎にその値を調整したい場合はありますよね。そんな場合はService側の定義でCPの値を上書くことが可能です。

```
/**
 * ECS Sercice作成
 * https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service
 */
resource "aws_ecs_service" "ecs_service" {
    .... 省略

    launch_type を設定している場合は削除

    capacity_provider_strategy {
        capacity_provider = "FARGATE"
        base              = 1
        weight            = 2
     }

    capacity_provider_strategy {
        capacity_provider = "FARGATE_SPOT"
        base              = 0
        weight           = 1
     }
}
```

このように設定することで個別の設定が可能になります。

## まとめ

今回は、Fargate Spot を使用してコンテナ料金を7割削減する方法を解説致しました。

Fargate Spotはかなり手軽にコスト削減できるものですので、上記を参考に是非導入をしてみてください。今回の記事が Fargate を使用している方の参考になりましたら幸いです。