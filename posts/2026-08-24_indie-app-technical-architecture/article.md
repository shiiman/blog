---
title: 個人開発アプリの技術構成（3層データ・認証・課金・双方向同期）
slug: indie-app-technical-architecture
date: '2026-08-24T12:00:00+09:00'
categories:
  - app
tags:
  - swift
  - kotlin
  - firebase
  - indie-dev
draft: false
id: 2253
excerpt: 個人開発でiOSとAndroidの両方を書いたバケットリストアプリの技術構成です。ローカルを正とする3層データレイヤー、iOSとAndroidで非対称にした認証、サーバー権威にした課金の権利判定、GitHubへの双方向同期キューの設計、Cloud Functionsをほぼ使わずSecurity RulesとApp Checkで組む方針をまとめました。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** の技術構成です。iOS と Android の両方をひとりで書いています。

何を作ったかは[前の記事](/engineering/app/indie-bucket-list-app-overview/)に、実装で踏んだ罠は iOS 編（[SwiftUI と SwiftData](/engineering/app/ios-swiftui-swiftdata-pitfalls/) / [課金とストア連携](/engineering/app/ios-storekit-widget-pitfalls/)）と [Android 編](/engineering/app/android-compose-room-billing-pitfalls/)に分けました。この記事は**どう組んだか**だけを書きます。

### この記事の表記について

Bundle ID は `com.example.app`、ドメインは `example.com` に置き換えています。アプリ名 `Bucket List Commit` は実名です。

---

## 全体の構成

| レイヤー | iOS | Android |
|---|---|---|
| 言語 / UI | Swift 6 / SwiftUI | Kotlin / Jetpack Compose |
| 非同期・状態 | Swift Concurrency / Observation | Coroutines + Flow / ViewModel + StateFlow |
| ローカル永続化 | SwiftData | Room |
| 端末間同期 | iCloud（CloudKit） | **なし**（バックアップと相互移行で代替） |
| データホーム | GitHub Issues + Projects v2 | 同じ方式を移植 |
| 共有・ソーシャル | Firestore | 同一コレクション・**同一 Security Rules** |
| 認証 | Apple Sign-In | **Google Sign-In** |
| クライアント保護 | App Check（App Attest） | App Check（Play Integrity） |
| 課金 | StoreKit 2 | Play Billing Library |
| ウィジェット | WidgetKit | Glance |
| 画像 | 自前のキャッシュ実装 | Coil |
| ネットワーク | URLSession を手で組む | Ktor |
| セキュアストア | Keychain | Keystore + EncryptedSharedPreferences |
| DI | なし（`@Environment` + サービス層） | **なし**（手書きの composition root） |

最低 OS は iOS 17 / Android 8（API 26）です。iOS 17 にしたのは SwiftData と Observation を使うためです。

DI ライブラリは両方とも入れていません。個人開発で層が浅いので、ライブラリを覚える手間のほうが大きいと判断しました。

## データは 3 層で、ローカルが正

一番大きな設計判断がここです。

```
[ローカル]  SwiftData / Room     ← これが Source of Truth
   ├→ [端末間]     iCloud（iOS のみ・自動）
   ├→ [データホーム] GitHub Issues + Projects v2（双方向）
   └→ [共有]       Firestore（公開・フレンド限定のものだけ）
```

**ローカルが正で、他はすべて派生**という原則にしました。理由は 2 つです。

1. オフラインで全機能が完結する。サインインしなくても使える
2. どこかのサービスが落ちても、ユーザーのデータは端末にある

この原則から出てくる帰結が 1 つあります。**Firestore は自分のデータの復元元にならない**（publish 専用で、ローカルへ書き戻さない）。だから安全網は GitHub とバックアップの 2 系統になります。Android には iCloud がないので、そのぶんバックアップを用意しました。

### 公開範囲で書き込み先が変わる

各バケットには 3 段階の公開範囲があります。ここが「どのストアに何が行くか」を決めています。

| 公開範囲 | ローカル | GitHub | Firestore | Cloud Storage |
|---|---|---|---|---|
| 非公開 | あり | あり | **なし** | **なし** |
| フレンド限定 | あり | あり | あり | あり |
| 公開 | あり | あり | あり | あり |

**非公開のものはクラウドの共有側に一切行きません。** これはプライバシーの方針でもあり、同時に**課金の構造でもあります**。Cloud Storage に画像が乗るトリガーが「非公開以外にした」の 1 軸だけなので、費用の伸びる場所が 1 つに絞れています。

1 ユーザーの Storage 使用量は `保存できる件数 × 1 件あたりの画像枚数` で上限が決まります（Free は 20×1、Pro は 1,000×3）。**プランの上限がそのままコストの上限になる**設計にしました。

## GitHub への双方向同期

USP の中核なので、ここが一番手がかかりました。各バケットを GitHub の Issue + Projects v2 のアイテムにミラーします。アプリと GitHub のどちらから編集しても反映されます。

- 操作したときにプッシュ、起動したときにプル
- 衝突は last-write-wins（`updatedAt` の新しい方を採用）
- 負けた側は捨てずに退避して、ホームの通知から差分を見て復元できる

### 同期キューは永続化しないと取りこぼす

最初は in-memory のキューだけで作りましたが、**フラッシュ前にアプリが kill されると同期操作が消えます**。そこでキューをローカルに永続化しました。

```
enqueue → 永続ストアに 1 レコード追加
       → flush（並列度 1・最大 3 回・指数バックオフ）
       → 成功したときだけ永続レコードを削除
起動時 → 前回やり切れなかったレコードを復元して再フラッシュ
```

バックグラウンドタスクは使っていません。操作したその場でフラッシュして、失敗したら次の起動で再開する方式です。個人開発でバックグラウンド処理を増やすと、動作確認できない経路が増えるだけだと考えました。

### 新規作成だけ冪等でない

同期の中で 1 か所だけ厄介なところがあります。**新規バケットの初回同期は Issue の POST なので冪等ではありません。** POST が通ってから Issue 番号をローカルに保存するまでの窓で kill されると、素朴な再送で Issue が二重に作られます。

3 段構えで防ぎました。

1. POST の直後に Issue 番号を**その場で保存**して、危険な窓を最小化する
2. Issue の本文に**不可視のマーカー**（`<!-- blc-bucket-id: UUID -->`）を埋め込む
3. 起動時の復元で「実行開始済みだが Issue 番号が未確定」のものだけ、POST の前にマーカーで既存 Issue を探して、あれば PATCH に倒す

3 で通信に失敗して確証が得られないときは、**再送を止めます**。二重作成より取りこぼしのほうがマシだと決めました。取りこぼしは次回に回せますが、二重作成はユーザーのリポジトリを汚します。

### 消す操作だけは絶対に捨てない

3 回失敗した同期操作は、無限に再試行せず退避してログに落とします。ただし**削除系の操作だけは例外**にしました。

公開をやめたのに Firestore のドキュメントが残る、画像を消したのに Cloud Storage に残る、という状態はプライバシーの契約を破ります。**カバー画像はトークン付き URL を知っていれば Security Rules を迂回して見られる**ので、実害があります。

そこで削除系の操作は別のストアへ退避して、**起動のたびにキューへ戻します**。成功するまで諦めません。同時に、恒久的に失敗している状態はホームの通知に 1 件出して、ユーザーに見えるようにしました。

## 認証は OS ごとに非対称

ここは意図的に揃えませんでした。

|  | サインイン |
|---|---|
| iOS | **Apple Sign-In のみ** |
| Android | **Google Sign-In のみ** |

理由は、それぞれのプラットフォームで最も摩擦が少ない手段が違うからです。iOS でガイドライン上の要件でもある Apple サインインを入れるなら、それ 1 つで足ります。Android で Apple サインインを出すと Web の OAuth 経路になって体験が落ちます。

そして重要な整理として、**Google と GitHub はサインイン手段ではありません**。

|  | 役割 |
|---|---|
| Google OAuth | カレンダーとスプレッドシートの連携専用 |
| GitHub Device Flow | データ用のプライベートリポジトリの同期専用 |

「身元」と「機能ごとのトークン」を分けて、設定画面もその 2 段構造にしました。連携を解除しても身元は残るし、サインアウトしても連携トークンの意味が変わらない、という整理です。

クロスプラットフォームの移行では、Apple と Google のアカウントを 1 つの Firebase UID にリンクします。

## 課金の権利判定はサーバーに寄せた

ここは作り直しました。最初はオンデバイスの購入情報だけで判定していましたが、iOS と Android で同じ Pro を扱う必要が出てきて破綻しました。

最終形はこうです。

```
クライアント: 購入の「証拠」を送る
   iOS     → StoreKit の署名付きトランザクション（JWS）
   Android → Play の purchaseToken
                    ↓
Cloud Function（callable）: 権利の「結論」を出す
   iOS の JWS   → Apple のルート証明書チェーンでオフライン検証
   Android のトークン → Play Developer API に問い合わせ
                    ↓
Firestore: users/{uid}/entitlement/current
   Rules は「本人の read と delete」だけ。create / update を持たない
```

判定は `オンデバイスの購入が有効 OR サーバーの権利ドキュメントが有効` です。

設計で効いたところを 4 つ。

- **Rules から書き込みを外したのが対策の本体です。** App Check だけでは塞げません。正規アプリから取り出したトークンを添えれば REST は通るからです。値域制限やレート制限では代替になりません（**攻撃は 1 回の書き込みで完了する**ので）
- サーバーは購入ハンドルを保管して、**どちらの OS から呼ばれても両ストアに問い合わせて最新化します**。これで「購入した OS のアプリを開かないと失効が反映されない」という初期の挙動が消えました
- ストア API を引けなかったプラットフォームは**前回の状態を温存**します。障害で権利を消してしまうのが一番まずいので
- 同じレシートを別のアカウントで使い回す横流し対策として、**実在が確認できた購入ハンドルを最初の UID に束縛**します。別 UID からの申告は「その購入は無い」ものとして落とします（エラーにはしません。端末ローカルの権利は無傷のままにしたいので）

## Cloud Functions はほぼ使わない

サーバーコードを増やさない方針にしました。運用を個人で背負えないからです。

| やり方 | 使う場面 |
|---|---|
| **Security Rules** | アクセス制御のほぼ全部 |
| **App Check** | クライアントの正当性の担保 |
| **クライアント側の計算** | 集計・ランキング・統計 |
| Cloud Functions | **本当に必要なところだけ** |

Cloud Functions を使ったのは、課金の権利検証と、予算の上限に達したときのキルスイッチだけです。どちらも「クライアントに置けない」ものです。

ランキングや統計をサーバーで集計しないので、クライアントの読み取り回数は増えます。ただ Firestore の無料枠が 1 日あたり読み取り 50,000 件なので、この規模なら収まります。**運用コストを増やさないほうが、個人開発では効きます。**

Security Rules は iOS と Android で**同じファイルを共有**しています。片方だけ緩いという事故が構造的に起きません。

## リポジトリの構成

```
App/         エントリポイント / ルーティング / アプリ状態
Core/        サービス層。Firebase / GitHub / Google / Sync /
             Persistence / Keychain / Store / Theme / Markdown /
             ImageProcessing / Export / Account / Models / UI
Features/    画面ごとの SwiftUI
Packages/    ローカル SPM（6 モジュールに切り出した共通処理）
Widgets/     WidgetKit
firebase/    Security Rules と Rules のテスト
Tests/       単体テストと UI テスト
android/     Kotlin / Compose の別モジュール
docs/spec/   仕様書（15 章）
```

### `.xcodeproj` はコミットしない

XcodeGen で `project.yml` から生成する運用にしました。プロジェクトファイルのコンフリクトが起きないのが目的です。

副作用が 1 つあります。**`.xcodeproj` を gitignore すると、そこに置かれる `Package.resolved` もリポジトリに入りません。** つまりロックファイルを持てません。

`upToNextMajor` のような範囲指定だと、**誰がいつ生成したかで解決されるバージョンが変わります**（手元と別のワークツリーで別バージョンがビルドされうる）。そこで **`project.yml` に厳密なバージョンを直接書く**ことにしました。更新は数値を書き換える = レビューできる差分になります。

### 共通処理はローカル SPM に切り出した

ネットワーク / テーマ / 画像 / Keychain / 触覚 / SwiftUI の補助を 6 モジュールに分けて、アプリ本体からは薄いファサード経由で使っています。ウィジェットや拡張ターゲットからも同じコードを使えるようにするためです。

## テストと検証のゲート

`make verify` に全部まとめました。

```
1. XcodeGen でプロジェクト生成
2. SwiftLint
3. iOS ビルド
4. iOS 単体テスト
5. Shared パッケージのテスト
6. Firestore / Storage Rules のテスト
7. 表示名・禁止語・埋め込みブラウザ不在などの構造ガード
8. Android ビルド
9. Android 単体テスト
10. Android lint
```

CI（GitHub Actions）では SwiftLint・Shared パッケージ・Rules テスト・Android の lint とテストを回しています。**macOS ランナーの iOS ビルドと単体テストは分単価が高いので手動実行だけ**にして、そのぶんは手元の `make verify` で担保する構成です。

### 構造ガードを 30 本以上置いた

「仕様に書いたのに、いつのまにか実装がズレる」を防ぐために、シェルスクリプトのガードを CI に並べました。

| ガードの例 | 何を凍結しているか |
|---|---|
| 表示名の検査 | 全プラットフォームの表示面でアプリ名の表記を統一 |
| ブランド名の翻訳禁止 | 日本語 UI でもブランド名を音写しない |
| 埋め込みブラウザ不在 | 外部リンクが OS のブラウザに渡ることを保証 |
| 新規の日本語文言の検出 | 英訳の付いていない文言を落とす |
| ハードコードされた色・フォントの検出 | テーマ経由でしか色を使えないようにする |

ここで 1 つ学びがありました。**ガードを置いたら「壊して落ちること」を必ず測る必要があります。** 実バグを再発させても合格を出していたガードが 3 本ありました。原因は、検索する前にコメントを剥がしていなかったことです（コードを消してもコメントに当たって通ってしまう）。

## Android を後から足して分かったこと

iOS を作り切ってから Android を書きました。順序としてこれは正解でした。仕様の判断がすべて済んでいるので、移植は手が止まりません。

一方で、**同じロジックを 2 回書くところは避けられませんでした**。

| 共有できたもの | 2 回書いたもの |
|---|---|
| Firestore の Security Rules | 画面（SwiftUI / Compose） |
| Firestore のコレクション設計 | ローカル永続化（SwiftData / Room） |
| GitHub のデータ形式 | ネットワーク層 |
| 定番バケットのシードデータ | Markdown の整形 |
| 課金の権利判定（Cloud Function） | 課金のクライアント側 |

Room のエンティティは SwiftData のモデルを **1 対 1 でミラー**しました（フィールド名・型・列挙を一致させる）。バックアップと移行のスナップショットを無変換で写せるようにするためです。ここを揃えておかないと、移行のたびに変換コードを書くことになります。

**iOS を UI の正典**と決めたのも効きました。差分が出たら Android を iOS に合わせる、と決めておくと「どちらが正しいか」の議論が発生しません。

## まとめ

構成の判断で効いたものを 5 つ。

1. **ローカルを正にする。** オフラインで完結し、どのサービスが落ちてもデータは残ります
2. **公開範囲を 1 軸にして、クラウドへの書き込みトリガーをそこだけに絞る。** プライバシーの方針とコストの上限が同じ設計で決まります
3. **課金の結論はサーバーが出す。** Rules から書き込みを外すのが本体で、App Check は補助です
4. **Cloud Functions を増やさない。** Security Rules + App Check + クライアント計算で足りるところは足ります
5. **Room を SwiftData の 1 対 1 ミラーにする。** 移行とバックアップの変換コードが要らなくなります

実装で踏んだ罠は [SwiftUI と SwiftData 編](/engineering/app/ios-swiftui-swiftdata-pitfalls/) / [iOS の課金とストア連携編](/engineering/app/ios-storekit-widget-pitfalls/) / [Android 編](/engineering/app/android-compose-room-billing-pitfalls/)に分けて書きました。

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

### このシリーズの他の記事

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

### ストアに出すまでのシリーズ

- [① 個人開発でストアに登録する前の事務手続き（屋号・開業届・D-U-N-S）](/engineering/app/indie-app-store-registration-paperwork/)
- [② 個人開発アプリの周辺インフラを用意して踏んだ罠（サイト・メール・OAuth・バックエンド）](/engineering/app/indie-app-domain-site-email-oauth-setup/)
- [③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/app-store-connect-review-submission-notes/)
- [④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/google-play-console-review-submission-notes/)
- [⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程](/engineering/app/personal-app-cross-store-release-full-journey/)
