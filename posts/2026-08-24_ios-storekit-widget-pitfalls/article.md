---
title: iOS の課金・ウィジェット・App Check で踏んだ罠（Debug では一度も通らない経路）
slug: ios-storekit-widget-pitfalls
date: '2026-08-24T14:00:00+09:00'
categories:
  - app
tags:
  - storekit
  - swift
  - ios
  - firebase
draft: false
id: 2255
excerpt: 個人開発アプリのiOS側で、課金・ウィジェット・App Check・App Intentsまわりに踏んだ罠をまとめました。App Checkが90日間1件もトークンを出していなかった話、Releaseの実機で初めて落ちた2件、ローカルビルドが通るのにアップロードが必ず失敗していた話、Siri発話が動かないと判定した経緯、Sandbox課金の寿命など。共通点は「Debugでは別実装に差し替わるので一度も実行されていなかった」ことです。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** の iOS 側で、課金・ウィジェット・クライアント保護・App Intents まわりに踏んだ罠をまとめます。

[SwiftUI と SwiftData の罠](/engineering/app/ios-swiftui-swiftdata-pitfalls/)は「ビルドは通るのに表示した瞬間に落ちる」型でしたが、この記事のものは全部別の型です。**Debug では別の実装に差し替わるので、本番の経路が一度も実行されていなかった**というものが並びます。

> このシリーズの他の記事です。[アプリの概要](/engineering/app/indie-bucket-list-app-overview/) / [技術構成](/engineering/app/indie-app-technical-architecture/) / [SwiftUI と SwiftData の罠](/engineering/app/ios-swiftui-swiftdata-pitfalls/) / [Android 版で踏んだ罠](/engineering/app/android-compose-room-billing-pitfalls/)

---

## App Check が 3 か月間 1 件もトークンを出していなかった

これが一番効いた発見です。

仕様書に「App Check は day 1 から入れる。任意ではない」と書いて、`AppDelegate` に App Attest のプロバイダファクトリを install していました。**それが一度も成立していませんでした。**

真因は App Attest の entitlement が無かったことでした。attestation が成立せず、**トークンが 90 日間 1 件も発行されていませんでした**。

```
firebaseappcheck.googleapis.com/services/verification_count  → 0 件
firebaseappcheck.googleapis.com/services/verdict_count       → 0 件
```

### なぜ全部のゲートを通り抜けたか

| 検出手段 | 実際 |
|---|---|
| ビルド / SwiftLint | **通る**（entitlement が無いだけで、コードは正しい） |
| 実行時 | **落ちない**。失敗はログに出るだけ |
| アプリの挙動 | **完全に正常に見える** |
| Debug での確認 | **不可能**（Debug プロバイダに差し替わる） |

この 4 行が、この型の本質を全部言っています。**「動いていないことに気づけない」経路**は、放っておくと永遠に気づけません。

対処として、Debug でモックに落ちる経路は **Release ビルドで値が無ければ即座に落ちる**ようにしました。動かないまま出荷されるより、ビルドの段階で止まるほうがマシです。

計測の注意も 1 つあります。**メトリクス名を間違えると 404 になります**（`request_count` は存在しません）。まず利用できるメトリクスの一覧を引いてから、名前を確定させました。

## Release の実機で初めて落ちた 2 件

OAuth の審査用にデモ動画を撮るために、**初めて Release ビルドを実機に入れたら起動直後に落ちました**。2 件出ました。

Debug・シミュレータ・単体テスト 1,507 件は全部通っていました。**根本原因は 2 件とも同じで、Debug が別実装に差し替わるので本番の経路が一度も実行されていなかったこと**です。

### 1 件目: Firebase の構成前にクライアントを作っていた

`FIRIllegalStateException` で落ちました。

**SwiftUI の `App.init()` は `didFinishLaunchingWithOptions` より前に走ります。** `FirebaseApp.configure()` は AppDelegate 側にあるので、`App.init()` で作られる Live クライアントは**必ず configure より前**に生成されていました。

しかも初期化の順序は動かせません。**App Check のプロバイダファクトリは configure より前に設定する必要がある**からです。そこでクライアント側を遅延解決に変えました（`final class` は計算プロパティに、`actor` は `lazy var` に）。

この修正で 1 つ気づいたことがあります。**同じ型の不具合を過去に 3 件直していて、その取り残しでした。** 「隣の同型が未配線」というパターンの再発です。

grep にも罠がありました。`= Firestore.firestore()` で検索すると、**先頭のドットを省略した `= .firestore()` を取りこぼします**。

```sh
# 取りこぼす
grep -rn '= Firestore.firestore()'
# 取りこぼさない
grep -rnE '= *\.(firestore|storage|auth|functions)\(\)'
```

### 2 件目: 2 つのストア設定の名前が両方 `default` だった

`Can't assign an object to a store that does not contain the object's entity.` で落ちました。

SwiftData のモデル設定を 2 つに分けたときに、どちらにも名前を付けていませんでした。両方が `default` になっていて、CoreData 上は「1 つの configuration に 2 つのストアがぶら下がる」形になり、オブジェクトの所属ストアを解決できなくなっていました。

**エンティティの割り当ては正しく排他だった**のが罠でした。設計は合っていて、名前だけが抜けていました。

再現テストにも罠がありました。

- この破綻は **CloudKit ミラーが有効なときだけ顕在化します**。ミラーを無効にすると CoreData が「そのエンティティを持つストアが 1 つだけならそれを使う」経路で救ってしまい、**シミュレータでは原理的に再現しません**
- そして **Debug のアプリは常にインメモリの単一設定で起動する**ので、2 設定の経路は Release の実機が初回実行でした

だから回帰を防ぐガードは、insert のスモークテストではなく**構造検査**に置きました。「configuration 名が一意であること」「エンティティの割り当てが排他かつ全モデルを覆っていること」を検査します。挙動を試すテストでは、そもそも壊れた経路に到達できません。

## ローカルビルドの成功は「アップロードできる」を意味しない

ビルド番号の方式を整えて、**初めて App Store Connect へアップロードを試したときに発覚しました**。

**このアプリは、それまで一度もアップロードできる状態になったことがありませんでした。**

`altool --validate-app` が `91179 Invalid extension` で必ず失敗します。しかも `make build` も archive も export も ipa 生成も、**全部成功します**。サーバー側のバリデーションでしか検出できませんでした。

原因は App Intents の拡張ターゲットの構成でした。`type: app-extension` として作っていたのですが、Apple は App Intents 拡張を **ExtensionKit の拡張**として要求します。

対処は**拡張を削除して本体 1 本にする**ことでした。ExtensionKit 化ではありません。理由は、その拡張には `@main` のエントリポイントが無く、**実行時に何も担っていなかった**からです（実機での検証も一度もされていませんでした）。**App Intents は本体ターゲットに入っていれば Shortcuts と Spotlight で動きます。**

削除の前に 1 つ作業が必要でした。**拡張専用の文字列カタログにあった 2 キーを本体のカタログへマージする**ことです。拡張が消えると、インテントは本体プロセスで本体のカタログを引くようになります。しかも**キーが日本語の原文なので、ja は動いて en だけ壊れます**。気づきにくい壊れ方でした。

そしてこの経験から、archive の最後に validate を自動実行するようにしました。

```
make ios-archive → archive → export → ipa → validate（自動）
```

**`--validate-app` はアップロードしないのでビルド番号を消費しません。** アップロードの前に必ず通す運用にしました。

## Siri の発話は動かないと判定した

App Intents を入れて、ショートカットアプリと Spotlight からの Quick Add は動きました。**しかし Siri に話しかけての実行は、実機でフレーズが一度もマッチしませんでした。**

アプリ名の解決（`INAlternativeAppNames`）は効いています。それでも発話では起動しません。2 回にわたって検証して、**「対応しない」と判定して受け入れました**。

タップでの実行（ショートカットアプリ / Spotlight）とオートメーションは動くので、`AppShortcutsProvider` はそのまま残しています。仕様書にも「Siri 発話での実行は未対応」と明記して、そこで止めました。

個人開発だと、こういう「原因が自分側にあるとは限らないもの」に時間を使い続けるのが一番危ないと思っています。**2 回測って同じ結果なら、諦めて仕様に書く**という判断にしました。

## ウィジェットは検証手段を先に作る

WidgetKit のウィジェットは、**シミュレータではプログラムからホーム画面に追加できません**。タップもできないので、見た目の確認手段がありませんでした。

そこで**アプリ内に Debug 用のギャラリー画面**を作りました。

- 描画ビューと Entry を `Core/Widget/` に置いて、**本体アプリとウィジェット拡張の両方でコンパイル**する
- `@Environment(\.widgetFamily)` ではなく **`family` を明示的な引数**にする（ウィジェット以外の文脈では environment を設定できないので）
- ウィジェット拡張側は「environment を読んで渡すだけ」の薄いラッパーにする
- ギャラリーは実寸に近いフレームでサンプルを描く

これで見た目の確認ができるようになりました。ただし、そのギャラリーにも罠がありました。

**テーマ背景のアセットはウィジェット側のカタログにしかありません。** ギャラリーで実際のテーマ ID を渡すと画像が空になります。既定背景で描いて、その旨を画面に注記しました。

英語表示の検証にも別の罠がありました。**ギャラリーは本体アプリのカタログを参照します。** ウィジェット専用の文字列は本体のカタログに無いので、言語を en にして起動しても **ja のまま表示されます**（実際のウィジェット拡張は自分のカタログを見るので en になる）。検証するには一時的にカタログをマージして、確認したら revert する必要がありました。

### ウィジェットの追従を測るときに 5 回間違えた

ウィジェットの値が本体に追従するかを測ったとき、**計測ミスを 5 回踏みました**。

| 踏んだ罠 |
|---|
| 「トグルを押したから期待値は 1 のはず」と推測して測り、トグルの意味を取り違えて全滅した |
| **前景復帰でウィジェットが更新される**ので、アプリに戻ってから読むと先に直ってしまう |
| 正規表現が `text="0 / 24"` にマッチせず、6 連続のタイムアウトをアプリの不具合と誤読しかけた |
| 同期していない状態から測って「2 秒で追従」と誤認した（期待値がたまたま現在値と一致していた） |
| ハーネスの中の戻る操作が画面遷移を壊して、以降の試行が全部無効になっていた |

信頼できるハーネスの条件は 2 つに整理できました。

1. **「何が起きたか」を推測せず、アプリ自身の表示と突き合わせる。** アプリのホームが出す達成数とウィジェットの値を比較すれば、操作の意味に依存しません
2. **OK / NG / 計測不能を分けて数える。** 一緒にすると、ハーネス自身の不具合を「不具合が再現した」と誤読します

## Sandbox 課金は寿命が検証の順序を決める

課金のテストで一番効いた知識がこれです。

|  | 更新の周期 | 打ち切り |
|---|---|---|
| iOS Sandbox の月額 | 5 分 | 12 回 |
| iOS Sandbox の**年額** | **1 時間** | 12 回 |
| Play のテスト定期購入 | 5 分 | **最大 6 回で自動キャンセル（約 30 分）** |

つまり **長く Pro の状態を保ちたい検証は、先に年額へ上げてから回します**。逆に「Pro を消した状態を見たい」ときは Play 側を解約すると次の 5 分境界で切れます。

実課金がゼロであることの確認も、Console を見に行く必要はありませんでした。**購入シート自身が言い切ります。**

- iOS: Sandbox アカウントで動いているので購入は常に無料
- Play: 「Test card, always approves」「You will not be charged.」と価格が `¥280/5 min` のように出る

**タップの前にこれを読むのが、唯一必要な安全確認**でした。

## 「購入したのに Free」の真因は非対称だった

Android 側で「Pro をテスト購入して Play は決済成功なのに、アプリは Free のまま」という不具合が出ました。「購入を復元」を押して初めて Pro になります。

原因は 2 つの重なりでした。

1. 購入イベントの購読開始が、単一端末リースのチェックより**後ろ**にあり、ロック中は early return で到達していなかった
2. 購入イベントの Flow が `MutableSharedFlow(extraBufferCapacity = 8)` で **replay 0** だった

2 つ目が本質でした。**`extraBufferCapacity` は emit を suspend させないためのバッファで、後から購読した相手には再配信しません。** 取りこぼしを防ぐのは `replay` の役目です。ここを取り違えていました。

面白かったのは、**iOS は構造的に安全だった**ことです。

```swift
// 生成時にコールバックを登録し、権利の「変化」をトリガにする
store.onEntitlementChange = { [weak self] _ in self?.verifyEntitlementWithServer() }
```

iOS は生成時にコールバックを登録して、権利の変化をトリガにしているので、購読開始のタイミング問題が起きません。**Android だけが「後から購読する Flow」方式**でした。

同じ機能を 2 回書くと、こういう非対称がどこかに残ります。**片方で見つけた不具合は、必ずもう片方の同じ場所を見る**というルールにしました。

## プライバシーの線引きを先に決めた

レビューで「ウィジェットや Siri で非公開のバケットを出すのは PII の漏洩だ」という指摘を受けました。これは**却下しました**。

このアプリは単一ユーザーの個人データアプリです。**「ユーザー自身のデータを、ユーザー自身の端末やユーザーが明示的に連携した外部サービスに出す」ことは第三者への漏洩ではなく、機能の目的そのもの**です。

この線引きを先にドキュメントに書いておいたのが効きました。同じ種類の指摘が来るたびに、そこを引用して判断できます。ロック画面への露出や Siri の読み上げというシナリオも検討したうえで、**オーナー基準で許容**と決めました。ウィジェットはタイムラインでロック画面とホーム画面を区別できないので、「出す / 出さない」の二択しかないという事情もありました。

一方で、**他のユーザーが書き込んだデータには防御的な検証を入れました**。通知に含まれる表示名は 50 文字で切り、アバターの URL は https のみを通します。ここは「自分のデータ」ではないので線の反対側です。

## まとめ

この記事の罠に共通していたのは 1 つでした。**Debug では別の実装に差し替わるので、本番の経路が一度も実行されていなかった**ことです。

| 経路 | Debug では |
|---|---|
| App Check | Debug プロバイダに差し替わる |
| Firebase のクライアント生成 | 遅延解決の問題が起きない |
| SwiftData の 2 設定 | インメモリの単一設定で起動する |
| 外部連携 | モックに差し替わる |

対策として置いたのは 3 つです。

1. **Release で値が無ければ即座に落とす。** 動かないまま出荷されるより、ビルドで止まるほうがマシです
2. **回帰ガードは挙動のテストではなく構造検査に置く。** 壊れた経路には、そもそもテストが到達できません
3. **archive の最後に validate を自動実行する。** ローカルビルドの成功は、アップロードできることの証明になりません

Android 側の罠は[こちら](/engineering/app/android-compose-room-billing-pitfalls/)に書きました。

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

- [③ SwiftUI と SwiftData で踏んだ罠（落ちる・効かない・反映されない）](/engineering/app/ios-swiftui-swiftdata-pitfalls/)
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
