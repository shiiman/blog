---
title: App Store と Google Play の申請からリリースまで（審査・却下対応・公開操作）
slug: indie-app-submission-to-release
date: '2026-08-24T10:00:00+09:00'
categories:
  - app
tags:
  - store-review
  - app-store-connect
  - google-play-console
  - in-app-purchase
draft: true
id: 2251
excerpt: 個人開発アプリをApp StoreとGoogle Playに申請してから公開するまでの記録です。Appleは初回審査でGuideline 2.1（Information Needed）に落ち、実機収録と審査ノートの書き直しで再提出して通りました。却下理由の全文、審査にかかった時間、対応の中身、初回リリース特有の提出の制約、承認メールに書かれたリリースの前提、届いたメールと使ったURLをまとめています。Playは審査中なので結果と公開操作は追記します。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発の iOS / Android アプリ **Bucket List Commit** を App Store と Google Play へ同時提出しました。この記事は**申請してから公開するまで**の記録です。提出までにやったことは別記事にしたので、ここは審査に出したあとの話だけを扱います。

Apple は初回審査で 1 回落ちました。指摘は **Guideline 2.1 – Information Needed** の 1 件だけで、コードの不具合はゼロでした。それでも対応の過程で実在のバグを 1 件見つけ、さらに初回リリース特有の提出の制約でつまずきました。最終的に再提出で通って、いまは公開操作を待つ状態です。

Google Play はこの記事を書いている時点でまだ審査中です。結果と、両ストアの公開操作については**この記事に追記**します。

> 提出までの手順は別記事です。[App Store Connect 編](/engineering/app/app-store-connect-review-submission-notes/) / [Google Play Console 編](/engineering/app/google-play-console-review-submission-notes/) / [事務手続き編](/engineering/app/indie-app-store-registration-paperwork/) / [周辺インフラ編](/engineering/app/indie-app-domain-site-email-oauth-setup/) / [全工程まとめ](/engineering/app/personal-app-cross-store-release-full-journey/)

### この記事の表記について

固有の識別子は伏せています（アプリ ID は `1234567890`、Bundle ID は `com.example.app`）。アプリ名 `Bucket List Commit` は実名です。

### 前提

- 新規アプリの初回リリース（バージョン 1.0）
- 課金あり（サブスク 2 種 + 買い切り 6 種）・UGC あり・配信地域は日本 + 米国
- 両ストアに同じビルド番号（`3581`）で提出
- リリースは両ストアとも**手動**設定

---

## 審査の返答までにかかった時間

一番知りたかったのはここでした。**提出してから結果のメールが来るまで**の時間です。

| 回 | 提出したビルド | 結果 | 返答まで |
|---|---|---|---|
| 1 回目 | 3581 | 却下（Guideline 2.1 – Information Needed） | **14 時間 20 分** |
| 2 回目 | 3583 | 承認 | **31 時間 36 分** |

**再審査のほうが長い**のは意外でした。「情報を足しただけだから早いだろう」という想定は当たりませんでした。1 回目が半日で返ってきたのも想定より早く、**返答時間は読めない**というのが結論です。

なお 2 回目は情報の追加だけでなく、収録中に見つけたバグの修正（build 3583）と課金商品 8 件を含む 10 項目の提出でした。**審査対象の量が増えている**ので、同じ条件での比較ではありません。

却下から再提出までにやったことは、この順番でした。所要時間は書きません（私が放置していた時間が混ざるだけで、参考にならないからです）。

1. 審査ノートを書き直す
2. 実機で画面収録する
3. 収録で見つけたバグを直してビルドを上げる
4. 提出物を組み直して再提出する

それぞれ後の節に書きました。

---

## 1 回目: Guideline 2.1 – Information Needed で却下

メールは 2 通同時に届きます。件名は `Your App Review Feedback` と `There's an issue with your〈アプリ名〉(iOS) submission.` で、内容は同じものです。

指摘は 1 件だけでした。

| | |
|---|---|
| 該当 | **Guideline 2.1 – Information Needed – New App Submission** |
| コードの不具合指摘 | **0 件** |
| ビルド | 3581 は `VALID` / 失効なし |
| 再提出 | **不要**（メールに明記） |

本文の要求はこれです（原文のまま）。

```
Reply in App Store Connect with all of the following information:

1. A screen recording captured on a physical device, running the latest
   operating system, demonstrating the app's functionality. The recording
   must begin with launching the app and show the typical user flow through
   its core features. If the app has any of the following, include them in
   the recording:
   - Account registration, login, and account deletion flows
   - Accessing paid content or features within the app, including any
     purchase or subscription flows
   - User-generated content, including content reporting and blocking
     mechanisms
   - Any prompts requesting access to sensitive data or device capabilities
     (for example, location, contacts, camera, or App Tracking Transparency)

2. A list of the device models and operating systems the app was tested on
   before submitting for review
3. A description of the app's functions and target audience, including the
   problem it solves and the value it provides
4. Instructions for setting up and accessing the app's main features,
   including any required login credentials or sample files
5. A list of the external services, tools, or platforms the app uses to
   deliver its core functionality
6. Describe any regional differences in the app's features or content, or
   confirm that the app functions consistently across all regions
7. If the app operates in a highly regulated industry or includes protected
   third-party material, provide any relevant documentation or credentials
```

7 項目のうち**本体は 1 番の画面収録**で、2〜7 は文章で答えられます。

そして重要な一文があります。**情報の追加だけなら再提出は不要で、App Store Connect の App Review ページから返信すればよい**とメールに書かれています。ビルドも失効していなかったので、ビルド番号を消費せずに済みました。

### 「情報が足りない」は個人開発だと踏みやすい

振り返ると、初回提出時の審査ノートは「デモアカウントの認証情報」と「報告 / ブロックの導線」しか書いていませんでした。上の 2〜7 のうち書けていたのは 4 だけです。

課金と UGC を持つ新規アプリでは、**実機の画面収録を最初から添えておく**のが正解だったと思います。ノートに 7 項目を書いておけば、この 1 往復ぶんの審査待ち（14 時間 + 32 時間）は消えていました。

---

## 却下理由は API では取得できない

対応の前に、状態を機械的に取れるか確かめました。**取れません。**

試して全部 404 だったエンドポイントです。

```
/v1/resolutionCenterThreads
/v1/apps/{appId}/resolutionCenterThreads
/v1/appStoreVersions/{versionId}/resolutionCenterThreads
/v1/appStoreReviewRejections?filter[app]={appId}
```

取れるのは状態だけです。

| エンドポイント | 取れる値 |
|---|---|
| `/v1/appStoreVersions/{id}` | `appStoreState = REJECTED` |
| `/v1/reviewSubmissions/{id}` | `state = UNRESOLVED_ISSUES` |
| `/v1/reviewSubmissionItems/{id}` | `state = REJECTED` |
| `/v1/appInfos/{id}` | `state = REJECTED` |

`appStoreVersions/{id}/appStoreVersionSubmission` は 404 です。新しい提出フローでは `reviewSubmissions` 側に移っています。

Google Play にも審査結果を返す API はありません。つまり**通知メールが唯一の情報源**です。審査に出したら、メールのフィルタと通知は先に整えておいたほうがいいです。私はこのタイミングで Gmail の API トークンが全プロファイル失効していて、メールを読むところから始まりました。

---

## 対応 1: 審査ノートの書き直し（上限は 4000 文字）

要求の 2〜7 を審査ノート（App Review Information の Notes）に入れました。API で更新できます。

```
PATCH /v1/appStoreReviewDetails/{id}
```

ここで上限にぶつかりました。**Review Notes は 4000 文字が上限**です。確実に超える長さ（10 万文字）を投げると、エラーメッセージが上限そのものを教えてくれます。

```
ENTITY_ERROR.ATTRIBUTE.INVALID.TOO_LONG
"Review Notes cannot be longer than 4000 characters."
```

**わざと超える値で叩いて上限を引き出す**のは、上書き事故が起きないので安全な調べ方です（不正な値は保存されません）。

旧 3954 文字を書き直して **3996 文字**に収めました。残り 4 文字なので、以後の追記の余地はありません。詳細版は返信のほうに置きました。

返信欄には文字数の上限はありませんが、**返信フォームにも 4000 文字の制限があります**。私の返信文は 8757 文字だったので、貼った瞬間に赤字で怒られました。

![App Store Connect の「App Review に返信」画面。返信欄が赤枠になり、下に「4000文字以内にする必要があります。」と残り文字数 -9268 が赤字で表示されている](./assets/review-reply-notes-limit.png)

`-9268` は超過分です。返信は分割して投稿しました。**長文を用意してから貼ると詰まる**ので、先に上限を確かめておくといいと思います。

---

## 対応 2: 実機で画面収録する

要求の 1 番です。ここで測っておいたことを並べます。

### 「最新 OS」の要件は実機の更新から

要求文に `running the latest operating system` とあります。収録当時の最新は **iOS 26.6.1**（2026-08-17 リリース）で、手元の実機は 26.6 でした。**更新しないと要件を外します**。

実機のモデルと OS は `xcrun devicectl` で取れます。

```sh
xcrun devicectl device info details | grep -i osVersionNumber
```

ペアリング済みなら **Developer Mode が OFF でも読めます**。

### Developer Mode は OFF のままでよい

TestFlight でのインストールと画面収録に Developer Mode は不要です。OFF のままで問題なく収録できました（`devicectl device info apps` だけが `CoreDeviceError 10005` で落ちます）。

ただし副作用が 1 つあります。**Developer Mode が OFF だと「設定 → デベロッパ → Sandbox Apple Account」の項目が無い**ので、TestFlight での購入は端末の App Store アカウントが Sandbox として使われます。登録済みの Sandbox テスターは使われません。実課金は発生しませんが、想定と違うアカウントで購入が走ります。

### 収録前に「Apple ID の使用を停止」する

Apple の要求は `Account registration ... flow` です。すでに一度サインインしている状態で撮ると、**初回の同意画面が出ず、ワンタップの「続ける」になってしまう**ので、登録フローの映像として弱くなります。

設定 → Apple アカウント → サインインとセキュリティ → Apple でサインイン → 該当アプリ → 「Apple ID の使用を停止」で、初回状態に戻せます。

### 台本を先に書く

要求された 4 種（登録 / ログイン / アカウント削除、課金、UGC の報告・ブロック、センシティブ権限のプロンプト）を 1 本に収める必要があるので、順番を先に決めました。

```
起動 → オンボーディング（規約同意ページを含む）
  → Apple でサインイン（初回同意画面から）
  → 主要機能をひと通り
  → ペイウォール → サブスク購入 → 権利が反映されるところ
  → 公開リストから他人の投稿 → 報告 → ブロック
  → 写真ライブラリの権限プロンプト
  → 設定 → アカウント削除（二重確認つき）
```

7 分 46 秒になりました。

### 添付は API から入れられる

動画は API で添付できます。

```
POST /v1/appStoreReviewAttachments   （appStoreReviewDetail に紐付け）
  → uploadOperations にマルチパートの PUT URL が返る
```

上限は実質ありません。`fileSize` に 5 GB を入れた予約が HTTP 201 で通りました。ただし**検証で作った予約は消してください**（`DELETE /v1/appStoreReviewAttachments/{id}`）。放置すると壊れた添付が残ります。

一方で、**Resolution Center への返信の投稿だけは API が非公開**です。ここは手作業になります。

### 収録すると本番にデータの残骸が出る

台本の最後にアカウント削除を入れたので、収録用アカウントのユーザーとバケットは消えます。それでも残るものがありました。

| 残るもの | 理由 |
|---|---|
| 報告のドキュメント | アカウント削除の対象外（別コレクションに書かれる） |
| 課金の権利・レシート | Sandbox 購入で生える |
| 退会済みユーザーの記録 | 退会直後の ID トークン保護用。時間が経つと消える |

提出前に掃除するスクリプトを 1 本用意して、全項目が通るまで消しました。**審査中に本番のデータを触るなら、何が残るかを先に洗い出しておく**必要があります。

---

## 収録したら実在のバグが 1 件出てきた

収録の途中で「購入を復元」が失敗トーストを出しました。**権利は無傷なのに失敗と表示される**バグでした。

真因は分岐の順序でした。同期結果の失敗判定を、ローカルの権利保持チェックより**先に**見ていたため、「Pro とテーマを全部持っていて増分ゼロ」のケースが失敗に落ちていました。サーバー側のログは同時刻に権利あり（`isPro=true`）を返していました。

面白かったのは Android 側です。同じ機能の Kotlin 実装は権利チェックを先に見て「すでに有効」に倒しており、**iOS だけが非対称**でした。片方の OS だけ順序が違う、というのは自分で書いていても気づけません。

直したあと、**順序を元に戻すと 3 回とも失敗する**ことを負テストで確かめました。「直った」ではなく「壊すと落ちる」を測らないと、直したかどうか分かりません。

このバグ修正でビルドを焼き直したので、番号が **3581 → 3583** に進みました（ビルド番号はコミット数なので、PR のマージコミット分も増えます）。Play は 3581 で審査中だったので、**両ストアの番号は不一致**になりました。機能への影響は無いので、そのまま進めています。

> **教訓**: 審査に落ちて実機で収録すると、シミュレータでは出ないバグが出ます。「情報が足りないだけ」の却下でも、収録は実質的な最終テストになります。

---

## 審査ノートが自分の申告と矛盾していた

これは提出前に気づいてよかった失敗です。

書いた審査ノートの 7 番（第三者の保護素材）に「第三者の保護素材は無い」と書きました。ところが App Store Connect の申告値は違っていました。

```
GET /v1/apps/{appId}
  → contentRightsDeclaration = USES_THIRD_PARTY_CONTENT
```

Apple の 7 番はまさにそこを聞いています。**そのまま出すと虚偽申告**でした。実際にはブランドロゴを 7 件同梱していて（各サービスの公式配布物を無改変で使用）、加えて実行時に OG 画像を取得しています。正典はリポジトリ内の出所台帳のほうでした。

> **教訓**: 審査ノートを書く前に `/v1/apps/{id}` と `/v1/appInfos/{id}` の申告値を読んで突き合わせる。
> 仕様書に正典があるのに、読まずに書いたのが原因でした。

---

## 再提出でつまずいた: 初回は課金商品をバージョンと同じ提出物に入れる

ここが一番時間を取られました。

却下されたあと、課金商品（サブスク 2 件 + 買い切り 6 件 + サブスクグループ 1 件）を提出物に足そうとしたら、**「審査へ提出できません」で止まりました**。

![App Store Connect の「提出物の下書き」画面。黄色い警告に「新しいサブスクリプショングループは、そのグループに属する自動更新サブスクリプションとともに提出する必要があります」「項目を審査に提出するには、選択したプラットフォームのアプリバージョンを追加してください」と表示され、その下に提出準備完了の項目が並んでいる](./assets/review-submission-items.png)

画面の警告は要領を得ませんが、API を叩くと**唯一の明文**が返ってきます。

```
POST /v1/reviewSubmissionItems
→ 409 ENTITY_ERROR.RELATIONSHIP.REQUIRED
   /data/relationships/appStoreVersionForReview
   "App 1234567890 must have an approved appStoreVersions for platform IOS,
    or an appStoreVersions must be included in this review submission."
```

条件は択一です。

1. **承認済みの iOS バージョンが既にある**
2. **同じ提出物にバージョンを含める**

初回リリースは 1 が存在しないので、**2 しかありません**。「1.0 が審査中」では条件を満たしません。UI の「選択したプラットフォームのアプリバージョンを追加してください」はこれを指しています。

### バージョンは 1 つの提出物にしか属せない

ところがバージョンは却下済みの提出物に握られていて、追加しようとすると弾かれます。

```
STATE_ERROR.ITEM_PART_OF_ANOTHER_SUBMISSION
```

ここで**やってはいけない操作**を踏みました。バージョンページ右上の「App Review に再提出」を押すと、**確認ダイアログなしでバージョン単独が飛びます**。課金商品は置き去りになります。これで 1 件だけ提出されてしまい、やり直しになりました。

### 正しい手順

3 手で終わります。課金商品の追加だけ UI が必須です。

```
1. PATCH /v1/reviewSubmissions/{却下済みの提出物}  {"canceled": true}
     → CANCELING → 数秒で COMPLETE
     → バージョンが DEVELOPER_REJECTED になって解放される
     （ビルド / 審査ノート / 添付は無傷。実測済み）

2. POST /v1/reviewSubmissionItems
     → 課金商品の下書きに appStoreVersion を追加
     （DEVELOPER_REJECTED でも 201 で通る）

3. PATCH /v1/reviewSubmissions/{下書き}  {"submitted": true}
     → WAITING_FOR_REVIEW
```

**課金商品・サブスク・サブスクグループの項目は API から追加できません**（`reviewSubmissionItems` にその関係が無い）。各詳細ページの「審査用に追加」ボタンだけです。バージョンは API で出し入れできます。

最終形は **10 項目 1 本**でした（バージョン 1 + 買い切り 6 + サブスク 2 + サブスクグループ 1）。

### 空の下書きは API では消せない

`reviewSubmissions` は DELETE に対応しておらず、`canceled: true` も「提出済みのもの」しか受けません。空の下書きを API で作ると消せなくなります（バージョンを提出するとき App Store Connect が自動で片付けてくれます）。

### 提出物の中身は base64 から読める

公開 API は提出項目の関係を返さないので、項目 ID をデコードするのが唯一の中身確認手段でした。

```sh
echo '<itemId>' | base64 -d
# → {submissionId}|{type}|{resourceId}
```

`type` の対応は `6` = アプリバージョン / `17` = アプリ内購入 / `18` = サブスクリプション / `19` = サブスクグループです。

---

## 2 回目: 約 32 時間で承認

再提出から約 32 時間で通りました。API で確認した状態です。

| 対象 | 状態 |
|---|---|
| バージョン 1.0（build 3583） | `PENDING_DEVELOPER_RELEASE` / `releaseType = MANUAL` |
| 提出物 | `COMPLETE` |
| サブスク 2 件 | `APPROVED` |
| 買い切り 6 件 | `APPROVED` |

**課金商品を同じ提出物に入れた判断は正解でした。** 入れていなければ、アプリだけ承認されて本番では何も買えない状態になっていました（審査ノートで購入手順を案内しているのに、です）。

承認は**メールが 2 通**来ます。先に事務的な通知、その 23 分後にお祝いのほうが届きました。

```
1 通目  Review of your〈アプリ名〉(iOS) submission is complete.
        → "Review of your submission has been completed.
           It is now eligible for distribution."  + Submission ID

2 通目  Welcome to the App Store
        → "Congratulations! ... has been approved for distribution."
```

後者に**リリースの前提が 2 つ**書かれています。

> Please note that it can take up to 24 hours for apps to become available on the App Store after release.
> **If your contracts are not yet in effect, your app cannot be distributed.**

公開操作をしてから実際にストアに出るまで最大 24 時間かかること、そして**契約が有効でないと配信されない**ことです。後者は次の警告に直結します。

アプリ一覧の表示はこうなります。「デベロッパによるリリース待ち」で止まっていて、**承認だけでは公開されません**。

![App Store Connect のアプリ一覧。Apple Developer Program 使用許諾契約の更新・EU のトレーダーステータス・年齢制限のソーシャルメディア新質問の 3 つのバナーが並び、その下にアプリが「iOS 1.0 デベロッパによるリリース待ち」と表示されている](./assets/asc-pending-developer-release.png)

### 承認と同時に出ていた警告 3 本

画面上部のバナーは 3 本でした。判定はこうです。

| 警告 | 判定 |
|---|---|
| Apple Developer Program 使用許諾契約が更新されました | **これだけ対応が必要**。承認メールの「契約が有効でないと配信されない」がまさにこれを指します。Account Holder 本人が同意する（無料・即時）。API に契約状態のエンドポイントは無いので、同意したかどうかは `Agreement signed: Apple Developer Program License Agreement` のメールで確認しました |
| EU のトレーダーステータス | **無関係**。配信 ON は日本と米国の 2 か国だけ |
| 年齢制限のソーシャルメディア新質問 | **既に回答済み**。未設定はキッズカテゴリ用の項目と任意項目だけ |

契約の同意は承認の 1 時間半後に済ませました。つまり**承認された時点では契約が有効でなかった**わけで、これを放置したままリリースを押すと配信されなかったはずです。承認メールを読まずに画面のバナーだけ見ていたら「あとで見るか」で流していたと思います。

配信地域の確認で 1 回誤読しました。`/v1/apps/{id}/availableTerritories` は 404（関係が存在しない）で、正しくは次の順に辿ります。

```
/v1/apps/{id}/appAvailabilityV2?include=territoryAvailabilities   ← total を見る
/v2/appAvailabilities/{id}/territoryAvailabilities?limit=200       ← 実データ
```

`included` 側にページングのリンクが無いので、v2 の関連パスに切り替えないと 50 件で打ち切られます。1 回目はそれで「配信 ON は米国だけ」と誤読しました。

---

## 審査中にやってはいけなかったこと

審査に出したあと、本番を触りたくなる場面が何度もありました。触ってはいけないものを書いておきます。

- **審査ノートが参照している本番データを消さない。** 私のノートには「公開リストから他人の投稿を開いて報告・ブロックを確認できる」と書いてあります。見せ用の公開投稿を消すと、その導線に到達できず **Guideline 1.2 の定型リジェクト**に戻ります
- **審査アカウントへの権利付与を外さない。** サーバー権威で権利を判定している設計だと、外した瞬間に有料機能が審査できなくなります
- **手元の端末で審査アカウントにサインインしない。** 単一端末のリースを持つ設計だと、リースが手元に移ってレビュアーがロック画面に当たります

3 つとも「掃除したい」「確認したい」という善意で踏みそうな操作でした。**審査中は本番を凍結する**と決めておくのが安全だと思います。

---

## 届いたメール

審査に関して届いたものだけを並べます。提出までのメールは[全工程まとめ](/engineering/app/personal-app-cross-store-release-full-journey/)と各編にまとめてあります。

| 件名 | 送信元 | 届くタイミング |
|---|---|---|
| Thank You for Submitting Your App | `no_reply@email.apple.com` | 審査に提出した直後（1 分後） |
| **Your App Review Feedback** | `no_reply@email.apple.com` | 却下された時。**却下理由の本文が入っているのはこちら** |
| There's an issue with your〈アプリ名〉(iOS) submission. | `no_reply@email.apple.com` | 上と同時（4 秒差）。内容は同じ |
| 〈アプリ名〉1.0.0 (〈ビルド番号〉) for iOS is now available to test. | `testflight_no_reply@email.apple.com` | ビルドを上げるたび |
| **Review of your〈アプリ名〉(iOS) submission is complete.** | `no_reply@email.apple.com` | **承認された時**。`eligible for distribution` と Submission ID だけの短い通知 |
| **Welcome to the App Store** | `no_reply@email.apple.com` | 承認の 23 分後。**リリース後に最大 24 時間かかること、契約が有効でないと配信されないことが書かれている** |
| Agreement signed: Apple Developer Program License Agreement | `developer@email.apple.com` | 使用許諾契約に同意した時 |

読み方の注意です。

- **却下も承認も 2 通ずつ来ます。** 件名が違うので別件かと思いましたが、却下の 2 通は中身が同じでした。承認の 2 通は違う内容で、**リリースの前提が書いてあるのは 2 通目**です
- **却下理由の本文はメールにしか無い**（前述のとおり API には無い）ので、**このメールを消さない**でください
- 承認は API でも判定できます（`appStoreVersions.appStoreState = PENDING_DEVELOPER_RELEASE`）。ただし**契約の状態は API に無い**ので、そこはメールでしか確認できません

---

## 参考リンク

**すべて公開時点で開けることを確認しました**（コンソール系はサインインへリダイレクトされます）。

| 用途 | URL |
|---|---|
| App Store Review Guidelines | <https://developer.apple.com/app-store/review/guidelines/> |
| App Review の全体像（提出・却下・再提出） | <https://developer.apple.com/distribute/app-review/> |
| TestFlight | <https://developer.apple.com/testflight/> |
| App Store Connect | <https://appstoreconnect.apple.com/> |
| 「ビジネス」（契約 / 税務 / 銀行） | <https://appstoreconnect.apple.com/business> |
| API キーの発行 | <https://appstoreconnect.apple.com/access/integrations/api> |
| Apple Developer のアカウント（契約への同意はここ） | <https://developer.apple.com/account/> |
| Play Console | <https://play.google.com/console/about/> |
| 審査のためにアプリを準備する（Play） | <https://support.google.com/googleplay/android-developer/answer/9859455?hl=ja> |
| Play のデベロッパー プログラム ポリシー | <https://play.google/developer-content-policy/> |

App Review ページと Resolution Center は App Store Connect のアプリ配下（バージョンページの「App Review」）にあり、直接の URL はアプリ ID を含みます。**アプリ一覧から辿るほうが確実**です。

---

## Google Play（審査中）

この記事を書いている時点では、まだ結果が出ていません。API で確認した状態です。

| 対象 | 状態 |
|---|---|
| production トラック | リリース `1.0.0` / `versionCodes = ['3581']` / `status = completed` |
| 配信国 | 日本 / 米国の 2 か国 |
| 掲載言語 | 日本語 / 英語 |

`status = completed` は**公開済みという意味ではありません**。これはロールアウト割合の状態で、100% を指しています。設定漏れではなく Google の審査待ちです。ここは一度読み間違えて「承認と同時に配信開始」と書いてしまい、後から覆しました。**Play の状態を API の値から推論しない**ほうがいいです。

初回アプリの審査は数日かかることがあるとされています。Play 側の通知メールも、この記事を書いている時点では 1 通も来ていません（審査結果のメールを送信元で検索して 0 件でした）。

結果が出たら、次の内容をこの章に追記します。

- **審査に出してから結果が返ってくるまでの時間**（Apple の 14 時間 / 32 時間と並べて比較できるように）
- 通ったか、落ちたか（落ちた場合は指摘内容と対応）
- 届いたメールの件名と送信元

---

## リリース（公開操作）

まだやっていません。両ストアの承認が揃ってから押します。

現時点で分かっていることを先に書いておきます。

| ストア | 公開の設定 | 押す場所 |
|---|---|---|
| App Store | `releaseType = MANUAL`（自動公開しない） | App Store Connect のバージョンページ |
| Google Play | 管理対象の公開: オン | Play Console の「Google Play にアプリを公開する」 |

**どちらも「承認 = 公開」ではありません。** Apple は `PENDING_DEVELOPER_RELEASE` で止まり、Play は初回リリースだとチェックリストに公開ステップが別に並びます。2 つ押すだけで同時公開になります。

Apple の承認メールによれば、**リリースしてからストアに出るまで最大 24 時間**かかります。「同時公開」といっても、実際に両方のストアで見えるタイミングは揃わない可能性があります。

公開後にやる予定のことも決めてあります。

- 審査用アカウントへの権利付与を解除する
- 審査ノートが参照していた見せ用データを消す（ノートの該当記述も外す）
- 公式サイトのストアリンクを有効化する（**公開までストアページは 404 なので、死んだリンクを踏ませないため後回しにしています**）
- 課金レシート検証の本番疎通をもう一度確認する

ここも実際に押したあとで追記します。

---

## まとめ（Apple 分・暫定）

Apple 側で得たものを 5 つだけ。

1. **課金と UGC があるなら、初回提出から実機の画面収録を添える。** Guideline 2.1 の情報要求は、後から審査待ち 1 往復（私の場合は 14 時間 + 32 時間）を足します
2. **却下理由は API では取れない。** メールが唯一の情報源なので、通知の経路を先に確保しておく
3. **初回リリースは、課金商品をアプリバージョンと同じ提出物に入れる。** 明文は API の 409 だけです。UI の警告文は要領を得ません
4. **審査ノートを書く前に、自分の申告値を読む。** ノートと申告が食い違うと、それ自体が虚偽申告になります
5. **承認メールは 2 通目まで読む。** リリースの前提（契約が有効であること・最大 24 時間かかること）は 2 通目にしか書かれていません

Play の結果と公開操作が済んだら追記します。同じところで消耗する人が 1 人でも減れば幸いです。

## 関連記事

**ストアに出すまでのシリーズ**

<a class="link-card" href="/engineering/app/indie-app-store-registration-paperwork/">
<span class="link-card__thumb"><img src="/eyecatch/indie-app-store-registration-paperwork.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">① 個人開発でストアに登録する前の事務手続き（屋号・開業届・D-U-N-S）</span>
<span class="link-card__excerpt">屋号の決定から屋号入り開業届の取り直し、D-U-N-S の 3 ルート比較まで。依存が一本道で並行できない工程の記録です。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.22</span>
</span>
</a>

<a class="link-card" href="/engineering/app/indie-app-domain-site-email-oauth-setup/">
<span class="link-card__thumb"><img src="/eyecatch/indie-app-domain-site-email-oauth-setup.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">② 個人開発アプリの周辺インフラを用意して踏んだ罠（サイト・メール・OAuth・バックエンド）</span>
<span class="link-card__excerpt">リダイレクトの末尾スラッシュで招待トークンが消える話、転送メールが返信できない話、OAuth 機密スコープ検証の差し戻し理由など。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.22</span>
</span>
</a>

<a class="link-card" href="/engineering/app/app-store-connect-review-submission-notes/">
<span class="link-card__thumb"><img src="/eyecatch/app-store-connect-review-submission-notes.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部</span>
<span class="link-card__excerpt">Individual 加入・W-8BEN・署名・TestFlight の 3 段・IAP 登録・年齢レーティング・審査ノート・提出まで。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.22</span>
</span>
</a>

<a class="link-card" href="/engineering/app/google-play-console-review-submission-notes/">
<span class="link-card__thumb"><img src="/eyecatch/google-play-console-review-submission-notes.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部</span>
<span class="link-card__excerpt">組織アカウントの本人確認・15% 手数料・署名鍵 3 種・アプリのセットアップ 11 項目・データセーフティ・課金テスト。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.22</span>
</span>
</a>

<a class="link-card" href="/engineering/app/personal-app-cross-store-release-full-journey/">
<span class="link-card__thumb"><img src="/eyecatch/personal-app-cross-store-release-full-journey.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程</span>
<span class="link-card__excerpt">3 週間のタイムライン、審査のために追加実装した機能、検証環境の作り分け、繰り返し踏んだ 9 つのメタな罠。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.22</span>
</span>
</a>

### アプリの中身と実装のシリーズ

- [① バケットリストアプリを個人開発した記録（動機・画面・機能・4か月の内訳）](/engineering/app/indie-bucket-list-app-overview/)
- [② 個人開発アプリの技術構成（3層データ・認証・課金・双方向同期）](/engineering/app/indie-app-technical-architecture/)
- [③ SwiftUI と SwiftData で踏んだ罠（落ちる・効かない・反映されない）](/engineering/app/ios-swiftui-swiftdata-pitfalls/)
- [④ iOS の課金・ウィジェット・App Check で踏んだ罠（Debug では一度も通らない経路）](/engineering/app/ios-storekit-widget-pitfalls/)
- [⑤ Android 版で踏んだ罠（Compose・Glance・Play Billing・エミュレータ検証）](/engineering/app/android-compose-room-billing-pitfalls/)
- [⑥ Security Rules が本番だけ壊れる（エミュレータでは再現しない話）](/engineering/app/firestore-security-rules-production-only-bugs/)
- [⑦ 個人開発アプリの課金とプラン設計（上限・買い切り・解約戻り・キルスイッチ）](/engineering/app/indie-app-pricing-and-entitlement-design/)
- [⑧ iOS と Android を往復できるようにした（アカウントリンクと単一端末リース）](/engineering/app/cross-platform-account-migration-ios-android/)
- [⑨ 報告とブロックを後から入れた話（ストアのUGCポリシー対応）](/engineering/app/adding-ugc-report-block-for-store-policy/)
- [⑩ iPad 対応を段階的に入れた（2ペイン化の判断とやらなかったこと）](/engineering/app/ipad-two-pane-staged-support/)
- [⑪ 「動いている」を証明する方法（実機・本番・計測の作法）](/engineering/app/proving-it-works-device-and-production-verification/)
- [⑫ iOS と Android を日英対応にした運用（String Catalog と values-en）](/engineering/app/ios-android-string-catalog-i18n/)
- [⑬ AI エージェントとひとりで両OSアプリを作った運用（仕様書・issue・構造ガード・記憶）](/engineering/app/solo-dev-with-ai-agent-workflow/)
