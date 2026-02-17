---
id: 1190
title: 【Snowflake】ユーザの作成からロール権限設定までまとめてみた | 初期設定
slug: snowflake-permission
status: publish
date: 2021-10-22T19:30:00
modified: 2021-10-22T01:04:14
excerpt: Snowflakeのユーザー作成からロール・権限設定まで、アクセス制御の基本をわかりやすく解説します。
categories: [18, 77]
tags: [78]
featured_media: 1191
---

こんばんは、しーまんです。

[Snowflake](https://www.snowflake.com/) を使用して権限周りでどうすればよいのか迷われた事はありませんか。

Snowflakeの権限周りはかなり細かく設定できる分、初めて触る方にとっては分かりづらい点があります。

そんな導入初期によく悩みになるユーザの権限周りの設定について今回解説できればなと思います。

また基本的なユーザ作成方法と権限周りの考え方を学び、ご自身の環境に合わせてカスタム出来るようになればとても良いと思っています。

そのための参考になれば幸いです。

## Snowflakeと管理オブジェクトについて

Snowflakeとアカウント内で操作・管理できるオブジェクトに関しては別記事で上げておりますので、こちらを参考にしてください。

 [![](https://shiimanblog.com/wp-content/uploads/2021/10/eyecatch_snowflake_object-320x180.jpg)\
\
【Snowflake】アカウント内のオブジェクトと組織の概念について \| 初期設定\
\
Snowflakeで使用するオブジェクトの意味を解説していきます。またその後オーガナイゼーションの概念について解説致します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2021.10.21](https://shiimanblog.com/engineering/snowflake-object "【Snowflake】アカウント内のオブジェクトと組織の概念について | 初期設定")

## 権限について

Snowflakeのアクセス制御は2つの仕組みが用意されています。

- **任意アクセス制御（DAC）：** 各オブジェクトに所有者がおり、所有者はそのオブジェクトへのアクセスを許可できます。
- **ロールベースのアクセス制御（RBAC）：** アクセス権限がロールに割り当てられ、ロールはユーザーに割り当てられます。

この2つの仕組みをベースに権限を管理することになります。

基本的にはロールを作って、ユーザに割り当てる。そのロールに割り当てられた権限の操作のみがユーザに許可されるという感じです。

またデフォルトではオブジェクトを作成したユーザのロールがそのオブジェクトの所有者となります。所有者と記載がありますが、ユーザではなくロールである点に注意してください。

詳細は [ドキュメント](https://docs.snowflake.com/ja/user-guide/security-access-control-overview.html) を参照ください。

### デフォルトロール

先の説明でロールをユーザに割り当てると説明しました。

Snowflakeではデフォルトで5つのロールが最初から存在しています。

- ACCOUNTADMIN
- SYSADMIN
- SECURITYADMIN
- USERADMIN
- PUBLIC

これらをユーザに割り当てて使用するのが最初のスタートです。

その後はご自身で権限を絞ってカスタムロールを作成し、必要な権限を渡す運用になるでしょう。

まずはデフォルトで用意されている5つのロールの説明をしていきます。

#### ACCOUNTADMIN

システムの最上位のロールであり、アカウント内の限られた/制御された数のユーザーにのみ付与する必要があります。「SYSADMIN」 および 「SECURITYADMIN」 を継承するロールです。

アカウント管理者の権限ですので、限られた数名に付与するロールになります。

#### SYSADMIN

アカウントでウェアハウスとデータベース（およびその他のオブジェクト）を作成する権限を持つロールです。

Snowflakeの権限管理で [推奨](https://docs.snowflake.com/ja/user-guide/security-access-control-considerations.html) されているのは、すべてのカスタムロールを 「SYSADMIN」 ロールに割り当てるロール階層を作成することです。このロールには、ウェアハウス、データベース、およびその他のオブジェクトに対する権限を他のロールに付与する機能もあります。

データを扱う管理者的なロールですね。

#### SECURITYADMIN

オブジェクトの付与をグローバルに管理し、ユーザとロールを作成、モニター、管理できるロールです。「USERADMIN」を継承するロールです。

#### USERADMIN

ユーザとロールの管理のみをおこなう専用のロールです。

- CREATE USERおよびCREATE ROLEのセキュリティ権限が付与されています。
- アカウントにユーザとロールを作成できます。このロールは、所有するユーザーとロールを管理することもできます。

#### PUBLIC

アカウント内のすべてのユーザーおよびすべてのロールに自動的に付与される疑似ロール。PUBLIC ロールは、他のロールと同様にセキュリティ保護可能なオブジェクトを所有できます。ただし、ロールが所有するオブジェクトは、定義上、アカウント内の他のすべてのユーザーとロールが使用できます。

通常、このロールは明示的なアクセス制御が不要で、すべてのユーザーがアクセスに関して平等であると見なされる場合に使用されます。

## Snowflakeの推奨どおりにオブジェクトを作成してみる

ではここからは実際にSnowflakeの推奨権限管理に則って各オブジェクトを作成してみましょう。

基本的には「ACCOUNTADMIN」を使用せずその他のロールでオブジェクトを作成していきます。

今回のゴールイメージは下記のようなイメージです。

- sysadmin -> データベース作成
- operator\_role -> 指定スキーマ内のテーブル作成+編集閲覧権限
- system\_role -> 指定スキーマ内のテーブル閲覧権限
- public -> readonly

どのオブジェクトにどんな権限をつけることが可能なのかはとても細かいので、 [ドキュメント](https://docs.snowflake.com/ja/sql-reference/sql/grant-privilege.html) を参照してください。

### ユーザ作成

まずはユーザを作っていきます。

今回作成するのは「一般ユーザ」と「BIユーザ」です。

ユーザを作成するのは「USERADMIN」ロールを使用します。

```
-- ユーザを操作できるロールに変更.
use role useradmin;

-- 一般ユーザ作成.
create user operator1
  password = 'operator1'
  login_name = 'operator1'
  email = 'operator1@gmail.com'
  display_name = 'operator1'
  must_change_password = true;

-- BIユーザ作成.
create user system1
  password = 'system1'
  login_name = 'system1'
  email = 'system1@gmail.com'
  display_name = 'system1'
  must_change_password = true;
```

### ロール作成

次にロールを作成してそれぞれのロールの継承関係を構築します。

そのあと、先程作成したユーザにロールを割り当てていきましょう。

ロールを作成するのは「SECURITYADMIN」ロールを使用します。

```
-- 権限を操作できるロールに変更.
use role securityadmin;

-- 一般ユーザロール作成.
create role operator_role;
-- sysadminはoperator_roleを継承.
grant role operator_role to role sysadmin;

-- システムロール作成.
create role system_role;
-- sysadminはsystem_roleを継承.
grant role system_role to role sysadmin;

-- 一般ユーザに一般ユーザロールを付与.
grant role operator_role to user operator1;

-- systemユーザにsystemロールを付与.
grant role system_role to user system1;
```

### ウェアハウス作成

次にSQLを実行するコンピュータであるウェアハウスを作成します。

ウェアハウスの作成には「SYSADMIN」ロールを使用します。

また今回ウェアハウスは誰でも使用可能なようにpublicロールに「usage」権限を付与します。

```
-- システムを操作できるロールに変更.
use role sysadmin;

-- warehouse作成.
create warehouse test_wh
    with warehouse_size = 'xsmall'
    warehouse_type = 'standard'
    auto_suspend = 10
    auto_resume = true;

-- test_whはpublic以上で使用可能な権限付与.
grant usage on warehouse test_wh to role public;
```

### データベース作成

続いてデータベースを作成しましょう。

データベースの作成には「SYSADMIN」ロールを使用します。

また今回DBは誰でも閲覧可能なようにpublicロールに「usage」権限を付与します。

```
-- システムを操作できるロールに変更.
use role sysadmin;

-- データベース作成.
create database test_db;
-- test_dbはpublic以上で使用可能.
grant usage on database test_db to role public;
```

#### スキーマ作成

続いて作成したデータベースにスキーマを作成しましょう。

ロールは引き続き「SYSADMIN」ロールを使用します。

```
-- システムを操作できるロールに変更.
use role sysadmin;

-- スキーマ作成.
create schema test_db.test_schema;
-- test_schemaはpublic以上で使用可能.
grant usage on schema test_db.test_schema to role public;
-- test_schemaはoperator_role以上で全権限.
grant all on schema test_db.test_schema to role operator_role;
```

#### テーブル作成

最後にテーブルを作成していきます。

テーブルを作成する前にスキーマ内の全てのテーブルに対して「将来の許可」を与えます。

毎回テーブル作成の度に権限を付与するのはとても大変なので、今後作られるテーブルに対して予めどんな権限を付与するのかを設定できます。

「将来の許可」は今回のようなスキーマに対する許可の他にデータベースレベルに対する設定も可能です。詳しくは [ドキュメント](https://docs.snowflake.com/ja/user-guide/security-access-control-configure.html#assigning-future-grants-on-objects) を参照ください。

```
-- 権限を操作できるロールに変更.
use role securityadmin;

-- operator_roleにtest_dbのtest_schema内のテーブルに対して編集権限を与える.
grant select, insert, update, delete on future tables in schema test_db.test_schema to role operator_role;
-- system_roleにtest_dbのtest_schema内のテーブルに対して閲覧権限を与える.
grant select on future tables in schema test_db.test_schema to role system_role;

-- システムを操作できるロールに変更.
use role sysadmin;

-- テーブル作成.
create table test_table (
    col1 int,
    col2 string
)
```

これで各オブジェクトへのアクセス許可が完了です。

「operator1」ユーザはテーブルの作成や編集・削除が可能で、「system1」ユーザは作成されたテーブルに対してselectクエリを発行することのみが可能になります。

## まとめ

今回はSnowflakeにおける権限の設定方法について解説しました。

まずはデフォルトロールの役割を解説した後、カスタムロールの作成から、ユーザにロールを割り当て、権限をつけていく方法を具体的に簡単な例を上げて説明しています。

今回の例では比較的簡単なアクセス制御でしたが、基本的には全てこれらの応用に過ぎません。今回の基礎的な権限割当が出来るようになれば、自身でどのように権限を管理したいかによって、正しい設定方法を設計出来るようになると思います。

権限周りは実際にサンプルを自分で試してみて、どの権限を与えればどのような挙動になるのか確かめることが理解をする上で大切になります。今回の例を自身で試して、挙動を体験してみましょう。

今回の記事が、Snowflakeでの権限設定の参考に少しでもなれば幸いです。