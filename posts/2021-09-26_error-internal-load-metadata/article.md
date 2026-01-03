---
id: 859
title: Dockerビルドエラー解決! | ERROR [internal] load metadata for &#8230;
slug: error-internal-load-metadata
status: publish
excerpt: こんばんは、しーまんです！！ 今回は私がdockerでbuildした際に「ERROR \[internal\] load metadata for ..」というエラーに遭遇した際の解決方法を紹介します。 ビルドした際に同じエ \[…\]
categories:
    - 18
    - 22
tags:
    - 35
featured_media: 860
date: 2021-09-26T19:30:00
modified: 2021-09-23T12:42:40
---

こんばんは、しーまんです！！

今回は私がdockerでbuildした際に「ERROR \[internal\] load metadata for ..」というエラーに遭遇した際の解決方法を紹介します。

ビルドした際に同じエラーが出て困っている方向けの記事になります。

遭遇したエラーは以下です。

```
[a12665@CA-20013877] % docker-compose up -d --build                           [main][11:20:59]
Building main
[+] Building 31.4s (4/4) FINISHED
 => [internal] load build definition from Dockerfile                                      0.0s
 => => transferring dockerfile: 452B                                                      0.0s
 => [internal] load .dockerignore                                                         0.0s
 => => transferring context: 68B                                                          0.0s
 => ERROR [internal] load metadata for docker.io/library/golang:1.17.1                   31.3s
 => ERROR [internal] load metadata for docker.io/library/alpine:3.13                     31.3s
------
 > [internal] load metadata for docker.io/library/golang:1.17.1:
------
------
 > [internal] load metadata for docker.io/library/alpine:3.13:
------
failed to solve with frontend dockerfile.v0: failed to create LLB definition: failed to do request: Head "https://registry-1.docker.io/v2/library/golang/manifests/1.17.1": net/http: TLS handshake timeout
ERROR: Service 'main' failed to build : Build failed
```

## エラー発生経緯

### エラー発生環境

エラーが発生した環境は下記です。

- MacOS Catalina 10.15.7
- Docker Desktop 4.0.1
- 自宅のwifi環境(楽天ひかり)

### エラー発生状況

エラーが発生した状況は下記です。

- goのバージョンを1.6.1から1.7.1にアップデートしようとした
- ローカルにpullしていたイメージを使用する場合問題なし

### エラー発生原因

エラーの文言からビルドに失敗しているというよりはイメージのpullに失敗しているように見えます。

なぜ突然エラーが出るようになったのか、詳しいところは分かりません。

Docker Desktop有料化とかの影響を受けているのでしょうか。。。

`docker login` をすると直るとかいう記事も見かけました。(こちらは検証しておりません)

いずれにせよDocker Desktopのバージョンアップの影響によるものが大きい気がしています。

(私の場合いつ頃バージョンアップしたのか覚えていませんでした)

## 解決方法

それでは私の場合の解決方法を紹介します。

まずDocker Desktopの設定を開きます。

右上のDockerアイコンをクリックし、「Preferences…」を選択します。

[![docker - 設定1](https://shiimanblog.com/wp-content/uploads/2021/09/docker_preferences.jpg)](https://shiimanblog.com/wp-content/uploads/2021/09/docker_preferences.jpg)

次に左メニューから「Docker Engine」を選択します。

右側にDockerの設定ファイルを変更出来るコード入力欄が表示されます。

そこに下記コードを設定します。

```
"insecure-registries": ["docker.io"],
```

私の場合はイメージがpullできないドメインが docker.io でしたので、そちらを指定しました。違うドメインで同じエラーが出た場合はそのドメインを設定してください。

[![docker - 設定2](https://shiimanblog.com/wp-content/uploads/2021/09/docker_setting-800x358.png)](https://shiimanblog.com/wp-content/uploads/2021/09/docker_setting.png)

上記の設定をするとDocker Desktopが再起動されます。

こちらの設定で私の場合はビルドが出来るようになりました。

## まとめ

今回はdocker build時に出たエラー「ERROR \[internal\] load metadata for …」についての解決方法を紹介しました。

Docker Desktop上でDockerのconfigファイルを修正するといったものでしたが、直接configファイルを修正して、再起動させても同じ結果になります。私の環境では下記がconfigファイルでしたね。

```
vi ~/.docker/daemon.json
```

またコマンドで回避する方法もありそうです。こちらは [公式リファレンス](https://docs.docker.com/engine/reference/commandline/dockerd/#insecure-registries) をご確認ください。

今回の記事で同じエラーが出て困っていた方の参考になりましたら幸いです。