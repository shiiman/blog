---
id: 437
title: 【docker】ローカルでマイクロサービス開発する方法
slug: microservice
status: publish
excerpt: 最近dockerやkubernetesなどといった技術が出てきたことにより、コンテナを使用してマイクロサービスでシステムを開発する手法が流行っています。 しかし実際にローカルでマイクロサービスの開発ってどうやっているので \[…\]
categories:
    - 18
    - 22
tags:
    - 35
    - 36
    - 37
featured_media: 438
date: 2021-09-08T19:30:00
modified: 2021-09-08T01:38:08
---

最近dockerやkubernetesなどといった技術が出てきたことにより、コンテナを使用してマイクロサービスでシステムを開発する手法が流行っています。

しかし実際にローカルでマイクロサービスの開発ってどうやっているのでしょうか。

今回はそんな疑問をお持ちの方に1つのマイクロサービスの開発例として、複数のコンテナに対してURIパスルーティングによるアクセス手法を紹介します。

マイクロサービスのアクセスパターンは「コンテナ間通信パターン」と「URIパスルーティングパターン」がありますが、今回紹介するのは「URIパスルーティングパターン」です。

![マイクロサービスパターン](https://shiimanblog.com/wp-content/uploads/2021/09/microservice_pattern-800x529.png)

複数のコンテナを扱いローカル開発を行う場合 docker-compose がよく利用されます。

そこで今回もdocker-composeを使ってローカル環境を作成していきます。

## ディレクトリ構成

最終的なディレクトリ構成を先に示します。

これからこちらのようにファイルを作成していきます。

```
.
├── README.md
├── conf
│   └── localhost
├── docker-compose.yml
├── service-a
│   ├── dockerfile
│   └── main.go
├── service-b
│   ├── dockerfile
│   └── main.go
└── service-c
    ├── dockerfile
    └── main.go
```

## 3つのアプリケーションを用意する

まずは3つのウェブアプリケーションを作っていきましょう。

今回アプリケーション自体は簡単な表示を返すだけのwebアプリケーションをGo言語で作成していきます。下記のようなmain.goファイルを用意します。

```
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service A")
	})
	http.ListenAndServe(":8080", nil)
}

```

実行させて http://localhost:8080/ にアクセスしてみましょう。

```
go run main.go
```

下記のように表示されたらOKです。

[![](https://shiimanblog.com/wp-content/uploads/2021/09/service_a_accecc.png)](https://shiimanblog.com/wp-content/uploads/2021/09/service_a_accecc.png)

サービスAのmain.goができたら続いてB, Cと作っていきましょう。

## dockerfile作成

アプリケーションができたら次にdockerfileを作成してきます。

以下のコードを各サービスディレクトリに配置しましょう。

```
FROM alpine:3.13 as tzdata
RUN apk --no-cache add tzdata

FROM golang:1.16.1

COPY --from=tzdata /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ='Asia/Tokyo'
ENV LANG='ja_JP.UTF-8'

```

## docker-compose.yml

次にいよいよ複数のコンテナを作成するdocker-compose.ymlファイルを作成していきます。

それそれのサービスをbuildしていきます。

複数のコンテナを操作する場合共通部分は抜き出しておきましょう。

作業ディレクトリの設定と8080ポートでのアクセスを許可しておきます。

```
x-common: &common
  working_dir: /go/src/microservice
  tty: true
  expose:
      - 8080
  command: bash -c "go run main.go"
```

サービスAの設定はビルドとボリュームのアタッチです。

ローカルのディレクトリを先程共通部分で設定した作業ディレクトリにアタッチしましょう。

```
service-a:
    <<: *common
    build:
      context: "./service-a"
    volumes:
      - ./service-a:/go/src/microservice
```

同じようにB, Cのコンテナも設定します。

## nginx-proxy

ここまでの設定だと8080ポートにアクセスした際に振り分け先が3つあるのでうまく動作しません。

そこでローカルにproxyを立てて、アクセスの振り分けをいきます。

今回はnginxを使用して簡単にプロキシできる [jwilder/nginx-proxy](https://hub.docker.com/r/jwilder/nginx-proxy/) を使用します。

### ルーティング先

まずは各コンテナへのルーティング先を決めておきます。

- サービスA

URI: localhost/api/service-a

コンテナ: service-a

- サービスB

URI: localhost/api/service-b

コンテナ: service-b

- サービスC

URI: localhost/api/service-c

コンテナ: service-c

### confファイル

次に上記で決めたルーティング先にアクセスを割り振るために設定ファイルを作成してきます。

localhostでアクセスをするのでconfファイル名は localhost として作成します。

そして各URIパスとコンテナを下記のように設定します。

```
location /api/service-a/ {
    proxy_pass http://service-a:8080;
}
location /api/service-b/ {
    proxy_pass http://service-b:8080;
}
location /api/service-c/ {
    proxy_pass http://service-c:8080;
}
```

proxy\_passに設定しているURLは http://\[コンテナ名\]:8080 となります。

コンテナ名はdocker-compose.ymlで設定しているサービス名です。

サービス名を変更する場合はこちらも合わせて変更してください。

### docker-compose.ymlへの追記

confファイルの作成が終わりましたら、docker-compose.ymlにプロキシコンテナのサービスを追記していきます。

先程作成したconfファイルを /etx/nginx/vhost.d にアタッチします。

```
nginx-proxy:
    image: jwilder/nginx-proxy:alpine
    ports:
      - "80:80"
    volumes:
      - /var/run/docker.sock:/tmp/docker.sock:ro
      - ./conf:/etc/nginx/vhost.d
    depends_on:
      - service-a
      - service-b
      - service-c
```

また今回上記で設定したルーティング以外のアクセスを受け入れるコンテナを作っておきます。

サービスA,B,Cどれかに設定してもいいのですが、今回は分かりやすいように別のコンテナを立てていきます。

OS情報とHTTPリクエスト内容を出力するGo言語で書かれた [containous/whoami](https://hub.docker.com/r/containous/whoami) を使います。

`nginx-proxy` のリバースプロキシ先として VIRTUAL\_HOST にホスト名を定義します。

今回でいうとlocalhostですね。

```
termination:
    image: containous/whoami
    environment:
      - VIRTUAL_HOST=localhost
```

## 動作確認

### 最終 docker-compose.yml

最終的に出来上がったdocker-compose.ymlは下記のようになります。

```
version: "3"

x-common: &common
  working_dir: /go/src/microservice
  tty: true
  expose:
      - 8080
  command: bash -c "go run main.go"

services:
  nginx-proxy:
    image: jwilder/nginx-proxy:alpine
    ports:
      - "80:80"
    volumes:
      - /var/run/docker.sock:/tmp/docker.sock:ro
      - ./conf:/etc/nginx/vhost.d
    depends_on:
      - service-a
      - service-b
      - service-c

  termination:
    image: containous/whoami
    environment:
      - VIRTUAL_HOST=localhost

  service-a:
    <<: *common
    build:
      context: "./service-a"
    volumes:
      - ./service-a:/go/src/microservice

  service-b:
    <<: *common
    build:
      context: "./service-b"
    volumes:
      - ./service-b:/go/src/microservice

  service-c:
    <<: *common
    build:
      context: "./service-c"
    volumes:
      - ./service-c:/go/src/microservice

```

### 起動

コンテナを起動させていきましょう。

下記コマンドで起動させることができます。

```
docker-compose up -d
```

### ブラウザアクセス

コンテナが起動したら各コンテナにブラウザからアクセスできるか確認してみましょう。

それぞれのアクセス先は下記になります。

▼ サービスA

http://localhost/api/service-a/

▼ サービスB

http://localhost/api/service-b/

▼ サービスC

http://localhost/api/service-c/

## まとめ

今回はdocker-composeを使用してローカルにマイクロサービス開発する環境を構築してみました。

実際にはwebサーバだけでなく共通DBとかキャッシュも使用すると思いますが、そこはそれほど難しくありません。接続先情報はdocker-composeの共通設定に追記してあげましょう。

はい、これでIngressのようなURIルーティングをローカルで再現させることができました。

基本的にはマイクロサービスなのでそれぞれサービス毎に開発を行えばいいのですが、クライアント側との接続をする場合などは、全てのアクセスが叩けないと効率が悪い場合があります。

そんな時に今回作成したURIルーティングを使用して複数コンテナを同時実行させると開発が効率的です。ぜひこの方法で開発してみてください。

## github

今回作成したソースはgithubにアップしておきましたので、よければ参考にしてください。

[https://github.com/shiiman/local-microservice](https://github.com/shiiman/local-microservice)