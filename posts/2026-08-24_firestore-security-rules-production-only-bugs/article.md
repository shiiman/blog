---
title: Security Rules が本番だけ壊れる（エミュレータでは再現しない話）
slug: firestore-security-rules-production-only-bugs
date: '2026-08-24T18:00:00+09:00'
categories:
  - app
tags:
  - firebase
  - firestore
  - security-rules
  - indie-dev
draft: false
id: 2257
excerpt: 個人開発アプリのFirestoreとCloud StorageのSecurity Rulesで踏んだ罠をまとめました。request.resourceはcreate/updateにしか存在せず、read/deleteで触ると式全体がエラー化してfail-closedでdenyされる。しかもエミュレータは解決してしまうので再現しません。372件のエミュレータテストがgreenのまま本番のバックアップ復元が一度も成功していなかった話と、RulesはPRのマージでは本番に反映されない話です。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** で、Firestore と Cloud Storage の Security Rules に踏んだ罠をまとめます。

このアプリは Cloud Functions をほとんど使わず、**アクセス制御のほぼ全部を Security Rules に寄せています**（[技術構成](/engineering/app/indie-app-technical-architecture/)に書きました）。運用を個人で背負わないための選択でしたが、そのぶん Rules の間違いが直接ユーザーの機能を殺します。

そして実際に殺しました。**エミュレータのテストが 372 件通ったまま、本番のバックアップ復元が一度も成功していませんでした。**

> シリーズ「アプリの中身と実装」全 12 本の 6 本目です。最初から読むなら [① アプリの概要](/engineering/app/indie-bucket-list-app-overview/) から。ストアに出すまでの手続きは[別シリーズ](/engineering/app/personal-app-cross-store-release-full-journey/)にまとめました。

---

## `request.resource` は create と update にしか存在しない

一番大きな罠でした。

Cloud Storage の Rules に、こう書いてありました。

```
allow write: if isOwner(uid) && request.resource.size < 5 * 1024 * 1024;
```

サイズ上限を掛けるつもりの、ごく普通の行に見えます。でもここには 2 つ問題が重なっていました。

### `allow write` は delete を含む

`write` は `create` / `update` / `delete` の総称です。つまりこの行は**削除にも適用されます**。

### read と delete で `request.resource` に触ると式全体がエラーになる

`request.resource` は **create / update の評価中にしか存在しません**。read や delete の評価で触ると `Property resource is undefined on object.` になり、**式全体がエラー化して fail-closed で deny されます**。

結果として、この 1 行が**削除を全部殺していました**。

そしてもう 1 つ、誤解しやすい点があります。

```
// これも「触る」ので同じくエラーになる
allow write: if isOwner(uid) && (request.resource == null || request.resource.size < N);
```

**`request.resource == null` というガードも「触る」ので回避になりません。** 「DELETE では `request.resource` が null になる」という誤った理解が出発点で、そのガードが 1 年以上 Rules に残っていました。

## エミュレータは再現しない

ここが本当に厄介なところです。

**ローカルの Rules エミュレータは、read や delete でも `request.resource` を解決してしまいます。** つまり**同じ Rules がローカルでは通り、本番だけ deny されます**。

| 環境 | read / delete で `request.resource` を触ると |
|---|---|
| ローカルのエミュレータ | 解決される（通る） |
| **本番** | **undefined でエラー化 → deny** |

私のところでは、**372 件のエミュレータテストが green のまま本番が壊れていました**。しかも Storage の delete のテストは 1 件も書いていませんでした。

### 実害

`backups/{uid}/snapshot.enc` の read が本番で 403 になっていました。

これが意味するところがこうです。

- アップロードは `create` なので**通る**
- 復元の read は**通らない**

つまり **Android のバックアップ復元は本番で一度も成功していませんでした**。「バックアップはあるのに戻らない」という状態です。ユーザーから見れば、安全網が存在しないのと同じでした。

### どう直したか

`request.resource` を触る条件は、**`create` と `update` だけに書きます**。

```
// 分ける
allow create, update: if isOwner(uid) && request.resource.size < 5 * 1024 * 1024;
allow read, delete:   if isOwner(uid);
```

そして**エミュレータで再現しないものは、エミュレータのテストでは守れません**。ここは構造的な検査（Rules のテキストを走査して「read / delete の allow に `request.resource` が現れないこと」を検査する）を CI に置きました。

## Rules は PR のマージでは本番に反映されない

これも実害が出ました。

**CI は Rules のテストを走らせるだけで、デプロイはしません。** Rules を変更した PR をマージしても、本番は古いままです。私のところでは、デプロイ後に 3 回 Rules を変更していて、**本番が陳腐化していました**。

コードは「マージすれば CI が本番に出す」ものが多いので、Rules だけ手動という非対称に気づきにくいです。

### 本番の実物を取って diff する

推測ではなく、本番に載っている Rules そのものを取ってきて比較します。

```sh
T=$(gcloud auth print-access-token)
Q='x-goog-user-project: <project-id>'

# どの ruleset が本番か
curl -s -H "Authorization: Bearer $T" -H "$Q" \
  https://firebaserules.googleapis.com/v1/projects/<project-id>/releases

# その ruleset の中身
curl -s -H "Authorization: Bearer $T" -H "$Q" \
  https://firebaserules.googleapis.com/v1/projects/<project-id>/rulesets/<ID>
# → .source.files[0].content が本番の実物
```

ここでも 1 つ踏みました。**`x-goog-user-project` を付けないと 403 が返ります**（認証情報が quota project を持たないため）。そして私は**その 403 の JSON を「releases が 0 件」と誤読して、「Rules が未デプロイだ」と誤報しかけました**。

> **教訓**: レスポンスの件数を数える前に、エラーキーの有無を見る。
> 0 件とエラーは別物です。

diff を取るときの注意もあります。**コメント行も差分に出る**ので、実質的な差分だけを見て緊急度を判断します。私のケースでは、見かけの差分は大きくても実質は 2 行でした。

## 「全パス deny」を 18 ラウンドの監査が見逃した

これは Rules の話というより、**判断の型**の話です。

Storage の Rules が**全パス deny** という状態になっていました。原因は 1 行です。

```
match /b/{bucket}        // ← 誤り
match /b/{bucket}/o      // ← 正しい
```

この不具合を、リリース前の監査 **18 ラウンドすべてが見逃していました**。

見逃した理由が記録に残っていました。

> 未変更の main でも同じ症状で立証 = 変更起因でない

**ここで論理が飛んでいます。** main と同じ結果であることは「既存バグである」ことの証明にしかならず、「環境要因である」ことの証明にはなりません。回帰がないことの確認と、バグでないことの確認は別物です。

しかも「環境要因」という結論が記録に残り、**リリース必須要件の項目にまで「CI / JDK 17 / 実機必須」と書かれて既定事実化していました**。1 行直したら 272 passed / 7 failed が **279 passed** になりました。

### 見分け方があった

後から振り返ると、**失敗の内訳を見れば分かるサインが出ていました**。

**落ちているのが肯定のアサーションだけ**（「〜できる」「read できる」）で、**deny 系が全部通っている**なら、それは「対象が全面的に拒否されている」シグナルです。全 deny なら deny のテストは自明に通るので、**「一部だけ失敗している」ように見えます**。

| 状態 | 肯定テスト | deny テスト | 見え方 |
|---|---|---|---|
| 正常 | 通る | 通る | 全部 green |
| **全パス deny** | **落ちる** | **通る（自明に）** | **一部だけ失敗** |

> **教訓**: テストの失敗を「環境要因」と結論するには、環境を変えて通ることを実証する必要がある。
> 実証していないなら「原因未特定の既存失敗」と書く。

そして仮説を立てたら、**1 行の実験で切り分けます**。`/o` を足して 279 passed になった時点で確定しました。5 分の作業でした。

## ブロック機能が恒久 deny になっていた

もう 1 件、Rules の設計そのものの誤りです。

ユーザーをブロックする機能で、**Rules 側でブロック相手からの書き込みを恒久的に deny する**実装にしていました。

一見正しく見えますが、これには問題があります。**ブロックを解除しても、Rules の評価が過去のブロック記録を参照し続ける**設計になっていて、解除後も相手の操作が通らない状態が残りました。

対処は、ブロックの効果を「そのときの状態」で評価するように直すことでした。**Rules に「恒久的に禁止する」という状態を持たせない**、という整理にしました。

Rules は宣言的で読みやすい反面、**「いつの状態で評価されるか」が曖昧になりやすい**と感じました。ここは条件を関数に切り出して、名前で意図を書くようにしました。

## エミュレータで検証する手順と限界

それでもエミュレータのテストは必要です。書き方だけメモしておきます。

```sh
# JDK が必要（バージョンの相性で落ちることがある）
cd firebase && npm test
```

踏んだ環境の罠が 2 つありました。

- **`~/.npm/_cacache` に root 所有のエントリがあると `npm ci` が EACCES で落ちます。** キャッシュの場所を退避させて回避しました
- **JDK のバージョン**で Rules エミュレータが起動しないことがあります。`JAVA_HOME` を明示するようにしました

そして限界のほうが重要です。

| 検証できること | 検証できないこと |
|---|---|
| 条件式のロジック | **read / delete での `request.resource`** |
| 認証状態ごとの可否 | 本番にデプロイされているか |
| ドキュメント間の参照（`exists` / `get`） | 実際の Storage オブジェクトの有無 |

**「エミュレータが green」は「本番で動く」の証明にはなりません。** ここを取り違えたのが、この記事の全部の始まりでした。

## iOS と Android で同じファイルを共有する

最後に、設計として効いたことを 1 つ。

**Security Rules は iOS と Android で同じファイルを共有しています。** コレクション設計も同じです。

これで「片方の OS だけ緩い」という事故が構造的に起きません。両 OS を書く個人開発だと、**同じものを 2 回書かない場所を探すのが効きます**。Rules はその最良の候補でした。

## まとめ

Security Rules で学んだことを 4 つ。

1. **`request.resource` は create / update だけ。** read / delete の allow に混ぜると、その操作が全部 deny になります
2. **エミュレータは再現しない。** 「エミュレータが green」は「本番で動く」の証明になりません
3. **Rules はマージでは本番に反映されない。** 本番の実物を API で取って diff します
4. **テストの失敗を「環境要因」で片付けない。** 落ちているのが肯定アサーションだけなら、全面 deny のサインです

Rules に寄せる設計そのものは、いまも正解だったと思っています。運用コストがゼロで、両 OS で共有できます。ただし、**本番でしか壊れない層を持つことになる**ので、そこを測る手段は別に用意する必要がありました。その話は[「動いている」を証明する方法](/engineering/app/proving-it-works-device-and-production-verification/)に書きました。

## 関連記事

**このシリーズの入口**

<a class="link-card" href="/engineering/app/indie-bucket-list-app-overview/">
<span class="link-card__thumb"><img src="/eyecatch/indie-bucket-list-app-overview.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">① バケットリストアプリを個人開発した記録（動機・画面・機能・4か月の内訳）</span>
<span class="link-card__excerpt">既存アプリへの 6 つの不満から何を作ると決めたか、画面と機能の全体像、約 4 か月・3,585 コミット・24 万行の内訳。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

<a class="link-card" href="/engineering/app/indie-app-technical-architecture/">
<span class="link-card__thumb"><img src="/eyecatch/indie-app-technical-architecture.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">② 個人開発アプリの技術構成（3層データ・認証・課金・双方向同期）</span>
<span class="link-card__excerpt">ローカルを正とする 3 層データ、OS ごとに非対称にした認証、サーバー権威の課金判定、GitHub への双方向同期キューの設計。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

### このシリーズの他の記事

#### 実装で踏んだ罠

- [③ SwiftUI と SwiftData で踏んだ罠（落ちる・効かない・反映されない）](/engineering/app/ios-swiftui-swiftdata-pitfalls/)
- [④ iOS の課金・ウィジェット・App Check で踏んだ罠（Debug では一度も通らない経路）](/engineering/app/ios-storekit-widget-pitfalls/)
- [⑤ Android 版で踏んだ罠（Compose・Glance・Play Billing・エミュレータ検証）](/engineering/app/android-compose-room-billing-pitfalls/)

#### 仕様と設計の判断

- [⑦ 個人開発アプリの課金とプラン設計（上限・買い切り・解約戻り・キルスイッチ）](/engineering/app/indie-app-pricing-and-entitlement-design/)
- [⑧ iOS と Android を往復できるようにした（アカウントリンクと単一端末リース）](/engineering/app/cross-platform-account-migration-ios-android/)
- [⑨ 報告とブロックを後から入れた話（ストアのUGCポリシー対応）](/engineering/app/adding-ugc-report-block-for-store-policy/)
- [⑩ iPad 対応を段階的に入れた（2ペイン化の判断とやらなかったこと）](/engineering/app/ipad-two-pane-staged-support/)

#### 運用と検証

- [⑪ 「動いている」を証明する方法（実機・本番・計測の作法）](/engineering/app/proving-it-works-device-and-production-verification/)
- [⑫ iOS と Android を日英対応にした運用（String Catalog と values-en）](/engineering/app/ios-android-string-catalog-i18n/)

### ストアに出すまでのシリーズ

- [① 個人開発でストアに登録する前の事務手続き（屋号・開業届・D-U-N-S）](/engineering/app/indie-app-store-registration-paperwork/)
- [② 個人開発アプリの周辺インフラを用意して踏んだ罠（サイト・メール・OAuth・バックエンド）](/engineering/app/indie-app-domain-site-email-oauth-setup/)
- [③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/app-store-connect-review-submission-notes/)
- [④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/google-play-console-review-submission-notes/)
- [⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程](/engineering/app/personal-app-cross-store-release-full-journey/)
- [⑥ App Store と Google Play の申請からリリースまで（審査・却下対応・公開操作）](/engineering/app/indie-app-submission-to-release/)
