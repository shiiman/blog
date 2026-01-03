---
id: 1793
title: 【GCP】GKE Autopilot Spot Podを導入してPod料金を60〜91%削減した話
slug: spot-pod
status: publish
excerpt: こんばんは、しーまんです。 最近AWSのネタばかりでGCPのことに触れていなかったので、今回はGCP関連の記事になります。 GCPでインフラを構築しているサービスで、最近では GKE Autopilot を使用することが \[…\]
categories:
    - 20
    - 18
tags:
    - 128
    - 129
    - 130
featured_media: 1794
date: 2022-10-28T19:30:00
modified: 2022-10-28T18:22:38
---

こんばんは、しーまんです。

最近AWSのネタばかりでGCPのことに触れていなかったので、今回はGCP関連の記事になります。

GCPでインフラを構築しているサービスで、最近では [GKE Autopilot](https://cloud.google.com/blog/ja/products/containers-kubernetes/introducing-gke-autopilot) を使用することが多くなってきました。本番環境ではそこまで導入は多くないイメージですが、ホスト管理が不要になるので、開発環境ではよく利用しています。

今回はそんなGKE Autopilotを使用している環境でSpot Podを導入してインフラ費用を削減しましたので、その話題に触れていきたいと思います。

もしGKE Autopilotを既に使用しているけど [Spot Pod](https://cloud.google.com/blog/ja/products/containers-kubernetes/announcing-spot-pods-for-gke-autopilot) の導入をまだしていないという方は、是非この記事を参考にして導入をしてみてください。かなり効果が高い費用削減が行えると思います！

## Spot Podとは

Spot Podとは、2021年11月に発表されたGKE Autopilotの機能の一つです。

[GCE Spot VM](https://cloud.google.com/spot-vms?hl=ja) と同様、GCPで管理している余剰リソースをいつ落ちても問題ない代わりに通常より安く使用できるというものです。

クラウドの余剰リソースが少なくなってくると、Spod Podsのコンテナには強制終了の25秒前にSIGTERMシグナルが発せられ、そしてシャットダウンが実行されます。

また、現在ではプレビュー版扱いらしいので、本番環境への導入は待ったほうがよいと思います。あくまでPodがいつ落ちても問題ない環境やバッチなどのサービスに限定して導入することをおすすめします。

### 費用

ということで Spot Pod の費用をみていきましょう！

GKE Autopilotの価格は、以下のページにまとまっています。

[GKE Autopilot の料金](https://cloud.google.com/kubernetes-engine/pricing)

**東京リージョン** における、オンデマンド(通常のPod)とSpot Podの料金を比較したのが下記表です。

種別PodSpot PodPodに対するSpot Podの料金比率per vCPU per hour$0.0571$0.0171約30%per GB per hour$0.0063215$0.0018964約30%

つまり通常の7割り引きで使用できるということですね。

めっちゃ安くなります！

※ 割引率はリージョンにより異なります。

## 導入方法

料金が安くなるということで、開発環境であれば導入しない手はないでしょう！

ということで早速導入方法をみていきましょう！

### GKEバージョンの確認

まずSpot Podを導入するにはGKEバージョンを1.21.4 以降にする必要があります。

Autopilotを使用していればバージョンのアップデートも自動でされると思うので、ここは問題ないでしょう！

### マニュフェストファイルの編集

後はKubernetesのマニュフェストファイルに少し項目追加するだけで導入可能です。

マニュフェストのサンプルは下記になります。

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hello-app
  template:
    metadata:
      labels:
        app: hello-app
    spec:
      containers:
      - name: hello-app
        image: us-docker.pkg.dev/google-samples/containers/gke/hello-app:1.0
      nodeSelector:
        cloud.google.com/gke-spot: "true"
      terminationGracePeriodSeconds: 25
```

実際に追加が必要な箇所は以下のみです。

```
      nodeSelector:
        cloud.google.com/gke-spot: "true"
      terminationGracePeriodSeconds: 25
```

こちらを追記してデプロイするだけでSpot Podを利用することが可能になります。

## 導入した結果

では、実際に導入してみた効果を見てみましょう。

下記のグラフを御覧ください。縦軸が「 **GCP使用料金**」で横軸が「 **日付**」です。

Spot Podを導入したタイミングから「 **GCP使用料金**」が下がっていることが確認できると思います。

[![price](https://shiimanblog.com/wp-content/uploads/2022/10/image.png)](https://shiimanblog.com/wp-content/uploads/2022/10/image.png)

## まとめ

今回は、Spot Pod を使用してGKE Autopilot のPod料金を6〜9割削減する方法を解説致しました。

Spot Podはマニュフェストファイルに数行追加するだけで実装可能です。かなり手軽にコスト削減できますので、上記を参考に是非導入をしてみてください。今回の記事が GKE Autopilot を使用している方の参考になりましたら幸いです。