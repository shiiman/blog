---
id: 1214
title: '【Snowflake】Macでsnowsqlを設定してみた | 「Failed to initialize log. No logging is enabled: [Errno 13] Permission denied」の対処方法'
slug: snwoflake-snowsql
status: publish
date: 2021-10-25T19:30:00
modified: 2021-10-22T12:40:03
excerpt: MacでSnowflakeのCLIクライアント「SnowSQL」をインストール・設定する方法と、Permission deniedエラーの対処法を紹介します。
categories:
    - 18
    - 77
tags:
    - 78
    - 80
featured_media: 1221
---

こんばんは、しーまんです。

Snowflakeにアクセスする際はWebからアクセスすることが多いと思います。その他のアクセス方法としては各プログラム言語のコネクタやドライバーが用意されているので、そちらでシステム化をすることが多いでしょう。

今回はもう一つよく使うCLIクライアントからのアクセス方法を紹介致します。

またその設定の際に、「Permission denied」エラーがでたので、そちらの解決方法も解説致します。

## snowsqlの設定方法

SnowflakeのCLIクライアントはsnowsqlといいます。

今回はMac環境でsnowsqlを設定していきます。

その他環境の設定に関しては [ドキュメント](https://docs.snowflake.com/ja/user-guide/snowsql.html) をご参考ください。

では早速インストールしてきましょう。

Macの場合はHomebrewで簡単にインストール可能です。

```
brew install --cask snowflake-snowsql
```

こちらでインストールは完了です。

デフォルトのターミナルシェルを使用する場合そのまま使用可能ですが、zshを使用している場合パスが通っておりませんので、「 **~/.zshrc**」にエイリアスを設定する必要があります。

私の場合もzshを使用しているので、「 **~/.zshrc**」に下記の設定を追加しました。

```
alias snowsql=/Applications/SnowSQL.app/Contents/MacOS/snowsql
```

これでパスも通せたので(正確にはエイリアスを貼っただけ)、インストール完了です。

## Permission denied エラーの対処方法

ここでインストール確認をしてみましょう。

```
snowsql -v
```

こちらはcliのバージョン情報を確認するコマンドです。

コマンドを実行すると下記エラーが表示されました。

```
Failed to initialize log. No logging is enabled: [Errno 13] Permission denied
```

何かしらをログに記録しようとして権限エラーが出ているようですね。

snowsqlはコマンド実行のログをデフォルトだと「 **../snowsql\_rt.log**」つまり、cliファイルの上の階層にログを出力しようとします。

Macの場合「/Applications/SnowSQL.app/Contents/MacOS/」ですね。

こちらへのファイル作成権限でエラーになっているようです。

対処方法は2つ考えられます。

1. 権限を変更する
2. ログの出力先を変更する

基本的には2のログファイル出力先を変更していきましょう。

ではどうやって出力先を変えていくかというと、snowsqlの設定ファイルというものが存在しておりますので、そちらの設定を変更することでログファイルの出力先を変更可能です。

設定ファイルは「 **~/.snowsql/config**」ですのでこちらのファイルを変更します。

変更点は以下です。

```
# main log file location. The file includes the log from SnowSQL main
# executable.
log_file = ../snowsql_rt.log

↓↓

log_file = ~/.snowsql/snowsql_rt.log
```

変更したら、再度「snowsql -v」コマンドを実行してみましょう。

今度はエラーなくcliのバージョンが表示されたと思います。

また「 **~/.snowsql/snowsql\_rt.log**」にログが出力されていることも確認出来ると思います。

snowflakeのエラーが表示された場合などはこのログを参照するようにしましょう！

## 実際にSQLを発行してみる

では実際にSnowflakeに対してSQLを発行してみましょう。

接続には以下の情報が必要です。

アカウント名: -a オプションで指定

リージョン: –region オプションで指定

ユーザ名: -u オプションで指定

パスワード: cliからの応答で指定

上記情報はSnowflakeのアカウント管理者に確認しましょう。

では実際に接続してみます。

```
snowsql -a <アカウント名> --region <リージョン> -u <ユーザ名>
```

コマンドを実行するとPasswordを聞かれますので、ユーザのログインパスワードを入力しましょう。

接続に成功するとSQL実行モードに入ります。

毎回接続情報を入力するのは面倒なので、接続先情報は先程「Permission denied」の時に編集した「 **~/.snowsql/config**」に設定することで接続情報の入力を省略することが可能です。

このSQLモードで与えられたユーザの権限の範囲でSQL発行が可能になります。

試しに下記コマンドを入力してみましょう。

ご自身で閲覧可能なデータベースの情報が確認出来ると思います。

```
show databases;
```

これでCLIを通してSnowflakeにアクセス出来るようになりましたね。

SQLモードから抜けるには「!exit」と入力します。

「!」を使用する点に注意しましょう。

## まとめ

今回はMac環境にSnowflakeのCLIクライアントであるsnowsqlの設定方法を解説しました。

インストールしてそのまま使用しようとすると「Permission denied」エラーが発生します。

その解決方法も合わせて確認できたと思います。

CLIを使えると、コマンドベースで簡単にアクセスできますので、コマンドに慣れている方はかなり助かるのではないでしょうか。

以上、少しでも参考になれば幸いです。