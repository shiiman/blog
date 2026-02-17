---
id: 464
title: 【docker】ローカルでWordPressを立ち上げる方法
slug: wordpress
status: publish
date: 2021-09-09T19:30:00
modified: 2021-09-11T03:32:18
excerpt: Dockerを使ってWordPressのローカルテスト環境を構築する方法を紹介。プラグインやテーマの検証を安全に行えます。
categories:
    - 18
    - 22
tags:
    - 6
    - 35
featured_media: 465
---

WordPressの設定やテーマを変えたり、新しいプラグインのインストールをした際に予想外の挙動を起こした経験はありませんでしょうか。

私は何度もあります。突然画像が表示されなくなったり、仮の文章が出てしまったり、表示崩れを起こしてしまったり、最悪の場合だとwebサイトにアクセスできなくなってしまうということもあります。

なので今回はローカルにテスト環境を構築してテーマの変更だったり、プラグインのテストだったりを確認したい。確認してから本番に反映したいと考えた事がある方のお悩みを解消していきます。

docker-composeを使うことにで簡単にローカルにテスト環境が構築できます。

ローカルならいくら環境を壊しても大丈夫です。

では早速やっていきましょう。

## ディレクトリ構成

まずはじめに今回作成する、最終的なディレクトリ構成を示します。

といっても docker-compose.yml ファイルとあといくつかディレクトリがあるだけのとてもシンプルな構成になっています。

```
.
├── README.md
├── docker-compose.yml
├── html
└── mysql
    ├── conf.d
    │   └── my.cnf
    ├── data
    └── log
```

## docker-compose.yml

では早速docker-composeファイルを作成してきましょう。

WordPressを実行するにはWordPress本体のコンテナとMysqlデータベースが必要です。

### WordPressコンテナ

WordPressのコンテナには3つの環境変数をセットします。

- WORDPRESS\_DB\_HOST : \[データベースサービス名\]:\[ポート番号\]
- WORDPRESS\_DB\_USER : データベースのユーザ名
- WORDPRESS\_DB\_PASSWORD : データベースのパスワード

\[データベースサービス名\]とは次に docker-compose.yml に設定するデータベースのサービス名になります。

次に先に示したhtmlディレクトリをコンテナの /var/www/html にアタッチします。

これはWordPressの設定ファイルがこのディレクトリ配下に配置されるため、アタッチしておくとローカルのhtml配下にファイルが作成されます。

/var/www/html をアタッチするメリットは3点あります。

1. 毎回立ち上げた際にそのディレクトリの設定を読み取るので、これまで設定した内容が消えずに再度中断時点から再開が可能になる。
2. そのディレクトリをgit管理にすることで、WordPressの状態をバージョン管理することが可能(DBデータを除く)
3. WordPressのソースコードを直接編集して修正・変更を加えることができる(上級者のみ)

上記の設定を反映させたのが以下のコードになります。

```
wordpress:
     depends_on:
       - db
     image: wordpress:latest
     volumes:
        - ./html:/var/www/html
     ports:
       - "8000:80"
     restart: always
     environment:
       WORDPRESS_DB_HOST: db:3306
       WORDPRESS_DB_USER: wordpress
       WORDPRESS_DB_PASSWORD: wordpress
```

### Mysqlコンテナ

次にMysqlのデータベースのコンテナを作成していきます。

Mysqlのコンテナは以下4つの環境変数を設定します。

- MYSQL\_ROOT\_PASSWORD : rootパスワード
- MYSQL\_DATABASE : データベース名
- MYSQL\_USER : ユーザ名
- MYSQL\_PASSWORD : パスワード

MYSQL\_USERとMYSQL\_PASSWORDはWordPressコンテナで設定したものを同じものを設定しましょう。

後は先に上げたディレクトリのうち以下のような紐付けでディレクトリをアタッチします。

mysql/log -> /var/log/mysql （ログディレクトリ）

mysql/conf.d -> /etc/mysql/conf.d （my.cnf(設定ファイル)ディレクトリ）

mysql/data -> /var/lib/mysql （データディレクトリ）

こちらの設定を反映させたのが以下のコードになります。

```
db:
     image: mysql:5.7
     volumes:
       - ./mysql/log:/var/log/mysql
       - ./mysql/conf.d:/etc/mysql/conf.d
       - ./mysql/data:/var/lib/mysql
     restart: always
     environment:
       MYSQL_ROOT_PASSWORD: root
       MYSQL_DATABASE: wordpress
       MYSQL_USER: wordpress
       MYSQL_PASSWORD: wordpress
```

## 動作確認

ここまでの設定でもうWordPressは起動します。

### 最終形 docker-compose.yml

```
version: '3'

services:
   db:
     image: mysql:5.7
     volumes:
       - ./mysql/log:/var/log/mysql
       - ./mysql/conf.d:/etc/mysql/conf.d
       - ./mysql/data:/var/lib/mysql
     restart: always
     environment:
       MYSQL_ROOT_PASSWORD: root
       MYSQL_DATABASE: wordpress
       MYSQL_USER: wordpress
       MYSQL_PASSWORD: wordpress

   wordpress:
     depends_on:
       - db
     image: wordpress:latest
     volumes:
        - ./html:/var/www/html
     ports:
       - "8000:80"
     restart: always
     environment:
       WORDPRESS_DB_HOST: db:3306
       WORDPRESS_DB_USER: wordpress
       WORDPRESS_DB_PASSWORD: wordpress

```

### 起動確認

それでは実際に起動して確認していきましょう。

下記コマンドでコンテナを起動させます。

```
docker-compose up -d
```

起動が完了したら http://127.0.0.1:8000 にアクセスしてみましょう。

WordPressの初期登録画面が表示されれば成功です。

[![wordpress - 初期設定画面](https://shiimanblog.com/wp-content/uploads/2021/09/wordpress-initial.png)](https://shiimanblog.com/wp-content/uploads/2021/09/wordpress-initial.png)

## おまけ

### WordPressのディレクトリ

WordPressを起動した時点でhtmlディレクトリにたくさんのファイルが作成されたのが確認できると思います。ここではWordPressがどのようなディレクトリ構造になっているのか軽く紹介します。

まずトップに3つのディレクトリがあると思います。

この3ディレクトリが大切なので簡単に概要だけ載せておきます。

カスタムする方以外は特に必要ない情報ですので、「へぇ〜WordPressってこんな感じでできているんだ〜」くらいに思ってもらえれば十分です。

フォルダ名概要wp-admin管理画面に関するファイルを格納するディレクトリwp-contentユーザの設定したテーマやプラグインなどを管理するディレクトリwp-includesシステムに必要なプログラムなどを管理するディレクトリ

次に wp-context の中をみていきましょう。

なぜならこのディレクトリがユーザの設定などを主に管理しているディレクトリだからです。

フォルダ名概要languages言語を管理するプラグインpluginsプラグインを管理するディレクトリthemesテーマを管理するディレクトリupgradeWordPressのアップデートの際に使用するディレクトリuploads画像などのアップロードファイルを管理するディレクトリ

## まとめ

今回はdocker-composeを使用してローカルにWordPressの実行完了を作成してみました。

ローカル環境で一度テストしてから本番のWordPressの設定を変更することにより、プラグインの導入やテーマの変更などをより安全にできるようになりました。

また合わせてWordPressのディレクトリの中身を軽く紹介しました。

実際にディレクトリの中にはもっとたくさんのファイルたちが存在していて、それらによりWordPressが動作しています。実際に触ることは無くてもそんなふうになってるんだなと認識くらいはしておくとよいでしょう。

## github

今回作成したソースはgithubにアップしておりますので、よければ参考にしてください。

[https://github.com/shiiman/wordpress](https://github.com/shiiman/wordpress)