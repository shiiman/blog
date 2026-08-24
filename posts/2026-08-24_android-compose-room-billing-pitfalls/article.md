---
title: Android 版で踏んだ罠（Compose・Glance・Play Billing・エミュレータ検証）
slug: android-compose-room-billing-pitfalls
date: '2026-08-24T19:00:00+09:00'
categories:
  - app
tags:
  - kotlin
  - jetpack-compose
  - android
  - room
draft: false
id: 2256
excerpt: iOS版を作り切ってからAndroid版を移植して踏んだ罠をまとめました。Glanceがウィジェットのセッション中はデータを読み直さない、キーボードでウィンドウがリサイズされないのでiOSの問題が起きない代わりに別の問題が起きる、Kotlinのtrailing lambdaがdefault引数で壊れる、values-enを追加するとjaフォールバック前提のテストが落ちる、Play配信で検証したつもりがサイドロード版だった、など。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** の Android 版で踏んだ罠をまとめます。

iOS 版を作り切ってから移植しました。ロジックは同じでも、**プラットフォームの前提が違うところで必ず割れます**。しかも iOS で問題になった箇所が Android では問題にならず、代わりに別の場所が壊れる、という非対称が多くありました。

> シリーズ「アプリの中身と実装」全 13 本の 5 本目です。最初から読むなら [① アプリの概要](/engineering/app/indie-bucket-list-app-overview/) から。ストアに出すまでの手続きは[別シリーズ](/engineering/app/personal-app-cross-store-release-full-journey/)にまとめました。

---

## Glance はセッション中データを読み直さない

これが Android で一番タチの悪い不具合でした。

症状は「**バケットを変更しても、ウィジェットは常にちょうど 1 手前の内容を描き続ける**」でした。アプリを再起動すると必ず直ります。

真因は `provideGlance` の構造でした。

```kotlin
// 壊れていた形
override suspend fun provideGlance(context: Context, id: GlanceId) {
    val snapshot = loadBuckets()      // ← provideContent の外側で読んでいた
    provideContent { WidgetBody(snapshot) }
}
```

**Glance はウィジェットのセッションが生きている間、`provideGlance` を再実行しません。** `updateAll()` は既存のコンポジションの再コンポーズを促すだけです。`provideContent` の外で読んだ値はそこに閉じ込められて、更新されません。

しかも**同じ理由でテーマ・カテゴリ・Pro 状態の反映も全部死んでいました**。どれも `updateAll()` で反映させる設計にしていたからです。

修正は `produceState` を使ってコンポジションの中から購読する形に変えました。

```kotlin
provideContent {
    val state by produceState(initial) {
        combine(buckets, theme, categories, isPro) { ... }.collect { value = it }
    }
    WidgetBody(state)
}
```

### 決め手になったのはログの非対称

仮説を **5 回外しました**。最後に効いたのは、ログの出方の違いでした。

```
アプリ起動時    : updateAll returned → provideGlance のログが出る → 正しい値
バケット変更時  : updateAll returned → provideGlance のログが出ない → 1 手前の値
```

「常に 1 手前」「再起動すれば直る」「`updateAll()` は成功を返している」の 3 つを**同時に説明できる仮説がこれだけ**でした。

外した仮説はこういうものでした。

| 外した仮説 | なぜ違ったか |
|---|---|
| 追加だけ更新されない | 同期が終わる前から測っていた計測ミス |
| Flow の collect が例外で恒久停止している | release でも毎回 collect のログが出ていた |

**「3 つの症状を同時に説明できるか」で仮説を絞る**のが一番速かったです。

## キーボードでウィンドウがリサイズされない

iOS で 2 つ問題になっていた箇所が、Android では**そもそも起きませんでした**。

| iOS の症状 | Android |
|---|---|
| キーボード回避でレイアウトが押し上がり、検索バーがタイトルに被る | 起きない（1px も動かない） |
| FAB が入力 UI のクリアボタンに重なる | 起きない（FAB がキーボードの裏に完全に隠れる） |

理由は `AndroidManifest.xml` に `windowSoftInputMode` の指定が無く、Compose は単一の `ComposeView` なのでスクロールできるかを判定できず、pan 相当になるためです。フォーカス中の入力欄は元から見えているので pan も起きません。API 34 と API 36 の両方で実測しました。

### 直そうとして悪化させた

「空状態の説明文とクリアボタンがキーボードの裏に隠れる」のを直そうと `imePadding()` を足したら、**空状態が丸ごと見えなくなりました**。

計算するとこうでした。

```
端末の高さ 731dp
  - キーボード 335dp
  - 上部のクローム 130dp
  = 残り 約 94dp

空状態の必要な高さ 約 250dp  → 上下にはみ出して全部クリップされる
```

`verticalScroll` と `heightIn(min = maxHeight)` を足しても改善しませんでした。そもそも `Column(Modifier.fillMaxSize())` を親の `Column` の子に `weight` なしで置いていたので、**残りではなく最大高さを取る**形になっていたのも効いていました。

ここで 1 つ想定が外れました。**`enableEdgeToEdge` を使っていないから `WindowInsets.ime` が 0 になるのでは、という懸念は外れました。** `imePadding()` はちゃんと効いてレイアウトが変わりました。つまり、**インセットが取れることと、それで直せることは別**でした。

## Kotlin の trailing lambda は default 引数で壊れる

関数型のパラメータの前に、後から default 付きのパラメータを挿入したら、**既存の呼び出しの一部だけがコンパイルエラーになりました**。

```kotlin
// こう変えた
class SignInIdentityGate(
    lastUidStore: Store,
    localProfileUid: () -> String? = { null },   // ← 後から挿入
    wipeLocalData: () -> Unit,
)
```

呼び出し側で起きたことがこれです。

| 呼び出し方 | 結果 |
|---|---|
| `SignInIdentityGate(store) { wipe }`（trailing lambda） | 最後の `wipeLocalData` にマップされる。通る |
| `SignInIdentityGate(store, onWipe)`（位置引数 2 つ） | `onWipe` が**中間の `localProfileUid` にマップ**されて型不一致 |

**`f(a) { }` は `f(a, { })` と等価で、lambda は「次の位置引数」にマップされます。default を読み飛ばして末尾には行きません。** 一方 trailing lambda は「最後のパラメータ」にマップされます。この 2 つは両立しないことがあります。

対処は、壊れる側を named 引数に直すことでした（`f(store, wipeLocalData = onWipe)`）。新しいパラメータを末尾に置くか中間に置くかは、**どの呼び出しパターンが多いか**で決めます。

そしてこの型のエラーは、**Robolectric のテストのコンパイル（`:app:compileDebugUnitTestKotlin`）で初めて落ちます**。本体のビルドだけでは通ってしまうので、Android を触ったら必ずテストまで通す運用にしました。

## values-en を追加すると既存テストが落ちる

英語対応で `values-en` を新設したら、**それまで通っていた Robolectric のテストが落ちました**。

追加前は、英語のロケールでも ja のリソースしか無いので ja に解決されていました。追加後は Robolectric の既定ロケール（ほぼ en）が en を解決するので、`getString` の日本語の本文を `contains("3 年前")` のように検証していたテストが落ちます。

修正は ja のコンテキストに固定することでした。

```kotlin
val jaCtx = ctx.createConfigurationContext(
    Configuration(ctx.resources.configuration).apply { setLocale(Locale.JAPANESE) }
)
```

`values-en` を追加するときは、**翻訳キーを `getString` しているテストを先に洗い出す**必要があります。

```sh
grep -rn "R.string.<key>" app/src/test/
```

### lint に止められるもの

`values-en` を作ると、既定（ja）の translatable なキー全部に en が必要になります（無いと `MissingTranslation` でエラー）。`translatable="false"` にしたもの（クライアント ID など）は不要です。

もう 1 つ、Compose の中で `LocalContext.current.getString(...)` を呼ぶと **`LocalContextGetResourceValueCall` の lint に止められます**（error 扱い）。`val ctx = LocalContext.current` を経由しても同じです。

対処は、Composable のスコープで `stringResource` を先に解決してから lambda に渡すことでした。動的な引数があるときはテンプレートを取って `String.format` します。

### KDoc に `*/` を書くとファイルが壊れる

これは笑ってしまった罠です。KDoc のコメントの中に `` `values*/strings.xml` `` のようにグロブを書くと、**`*/` でブロックコメントが閉じます**。以降が全部コードとして解析されて「Expecting a top level declaration」が大量に出ます。

パスやグロブをコメントに書くときは `*/` の連続を避けます。検査はこれで足ります。

```sh
grep -nE '[^ ]\*/' <file>
```

### 日付の書式は ICU のスケルトンで揃う

iOS と Android で日付の表示を一致させるのに、少し探しました。

| 使うもの | 結果 |
|---|---|
| `DateFormat.getBestDateTimePattern(locale, "yMd")` | en は `5/18/2026` / ja は `2026/5/18`。**iOS と一致** |
| `DateTimeFormatter.ofLocalizedDate(SHORT)` | en が 2 桁年（`5/18/26`）になって iOS と不一致 |

スケルトンの単文字 `M` / `d` が非ゼロ埋めのロケール順を返すので、iOS の `.month(.defaultDigits).day(.defaultDigits)` と揃います。

### 文字列のパリティは機械照合するしかない

iOS と Android で同じ文言になっているかを目視で確認するのは無理でした。ダッシュの種類（`–` U+2013）、波ダッシュ（`〜` U+301C）、マイナス（`−` U+2212）、ASCII のチルダの取り違えは目では見えません。

Android の XML と iOS のカタログを機械的に読んで、**書式指定子の方言を正規化してから逐語比較**するスクリプトを書きました。

```
%lld / %1$lld  ↔  %d / %1$d
%@             ↔  %s
```

これをやったら、**日本語でも iOS と Android で文言が割れていた箇所**が出てきました（シート名のスペース、割合の表記、状態のラベル）。英語対応の作業で炙り出た形です。過去のレビューは 5 巡していてロジックの一致は見ていましたが、**文字列の完全一致までは追えていませんでした**。

## Apple 連携は iOS が動いていても Android だけ落ちる

クロスプラットフォームの移行で、Android から Apple アカウントをリンクする必要がありました。ここで落ちました。

**Firebase の Apple プロバイダは、iOS ネイティブ用と web OAuth（code flow）用で設定が分かれています。** iOS のネイティブな Apple サインインは code flow を使わないので、**Services ID / Team ID / Key ID / .p8 が空でも iOS だけは動きます**。

つまり **iOS が動いていることは何の保証にもなりません**。

切り分けは Identity Toolkit の `createAuthUri` を直接叩くのが最速でした。

```sh
curl -s -X POST "https://identitytoolkit.googleapis.com/v1/accounts:createAuthUri?key=<web API key>" \
  -H 'Content-Type: application/json' \
  -d '{"providerId":"apple.com","continueUri":"https://example.firebaseapp.com/__/auth/handler"}'
```

| 応答 | 意味 |
|---|---|
| `400 OPERATION_NOT_ALLOWED : Code flow is not enabled for Apple.` | code flow の設定が無い |
| `authUri=https://appleid.apple.com/auth/authorize?response_type=code&...` | 設定済み |

**エラーの文言が「プロバイダが無い」ではなく「code flow が無い」のが要点**でした。前者ならプロバイダ自体が無効という意味になります。対照群として `google.com` を同じ手順で叩くと切り分けが早くなります。

ここで使う `.p8` は、App Store Server API のキーや App Store Connect API のキーとは**別種で使い回せません**。再ダウンロードもできません。

## Play 配信で検証したつもりがサイドロード版だった

これは検証の前提が崩れていた話です。

エミュレータに入っていたビルドが、**`installerPackageName=null` = `adb install` で入れたサイドロード版**でした。Play 配信版は Play のアプリ署名鍵で署名されるので、ローカル鍵で入れた APK への上書き更新は**構造的に不可能**です。

```
INSTALL_FAILED_UPDATE_INCOMPATIBLE: Existing package signatures
do not match newer version; ignoring!
```

**Play ストアの画面には `Can't install` としか出ません。** logcat を見ないと原因に辿り着けません。

判定はこれで確定します。

```sh
adb -s <serial> shell dumpsys package <package> | grep -E 'versionCode|installerPackageName'
#  installerPackageName=com.android.vending  ← これが無ければ Play 配信ではない
```

エミュレータのスナップショットが巻き戻ったあとは特に疑う必要があります。巻き戻り先がサイドロード時代なら、**以後の「Play 配信で検証した」は全部前提が崩れています**。

もう 1 つ、機内モードでの検証で踏んだものがあります。**機内モードでの Firebase Auth のサインインは失敗しません。** クライアント側のタイムアウトを持たないので、9 分以上ハングします。「オフラインで失敗する経路」を測るつもりが、単に固まっているだけでした。

## エミュレータの内部状態を読むには環境を選ぶ

検証で内部状態まで見たいことがよくあります。ただし**環境によって読めるものが変わります**。

| 環境 | 読めるか |
|---|---|
| Play イメージ + release ビルド | **一切読めない**（`adb root` が拒否され、release は `run-as` も不可） |
| Google API イメージ + debug ビルド | 読める |

読める場所は `shared_prefs/` の下です。

| ファイル | 中身 |
|---|---|
| アカウントの prefs | 最後にサインインした UID などが平文で入る |
| 設定の prefs | 各種設定値が平文 |
| Firebase Auth の store | **暗号化されていて UID は読めない** |

ただし、**書き換わった瞬間は md5 で読めます**。サインインやサインアウトのたびにこのファイルが書き直されるので、**md5 の変化を「サインインが成立した瞬間」のトリガーに使えます**。これがあると「ある操作がセッションの窓の内側に着弾した」ことを推定ではなく断定できました。

データを消す検証のカナリアには、`.clear()` でキーごと消える設定値を使います。既定値と違う値を入れておけば、UI でも prefs でも判定できます。

### 自動操作で踏んだ 2 つ

- **Compose の `clickable(enabled = false)` は uiautomator の `enabled` 属性に出ません。** `clickable="false"` を見ます。ここを間違えると「処理中」の状態を一度も検出できず空振りします
- **日本語 IME が `adb shell input text` を変換します**（メールアドレスが全角に化けます）。キーイベントを直接送っても composing に入るので、IME の変換候補に出る正しい ASCII をタップして確定させました

## 偽の UID ではソーシャル機能は動かない

デバッグメニューに認証状態を差し替える仕組みを入れていました。**これは observable な状態だけを差し替えるので、実際の認証は変わりません。**

つまり偽の UID を入れた状態では、共有・QR・招待リンクのようにサーバーと通信する機能は動きません。UI の分岐を確認するには使えますが、機能の検証には使えませんでした。ここを取り違えて「機能が壊れている」と誤診しかけました。

## iOS を UI の正典にした

パリティの監査を何巡もやりました。そこで決めておいてよかったのが、**iOS を UI の正典とする**という原則です。

差分が出たら Android を iOS に合わせます。「どちらが正しいか」の議論が発生しません。Android 側の過去の独自の逸脱を「意図的な設計」とみなしてスコープから外さない、というところまで決めておきました。

そのうえで、**iOS で指摘・修正した点は Android にも先回りで反映する**ようにしました。同じ指摘を Android でもう一度受けるのは、単に手戻りです。

ただし逆のパターンもありました。[削除ボタンの色](/engineering/app/ios-swiftui-swiftdata-pitfalls/)は Android が最初から正しく赤で、**iOS だけが壊れていました**。Android は `MaterialTheme.colorScheme.error` で最初から赤なのに、iOS はタブ全体の tint に負けて青くなっていました。正典を決めても、**実測しないと「どちらが壊れているか」は分かりません**。

## まとめ

Android 版を移植して分かったことを 4 つ。

1. **プラットフォームの前提が違うところで必ず割れる。** キーボードの挙動のように、iOS の問題が消える代わりに別の問題が出ます
2. **「更新されない」系は、フレームワークがいつ再実行するかを疑う。** Glance の件は API の呼び出し自体は成功していました
3. **検証の前提を毎回確認する。** 「Play 配信で検証した」は `installerPackageName` を見るまで言えません
4. **文字列のパリティは機械照合しかない。** 目視ではダッシュ 1 文字の違いを見つけられません

そして一番効いたのは、**iOS を作り切ってから移植した**ことでした。仕様の判断が全部済んでいるので、手が止まりません。同じ機能を 2 回書く手間はありますが、**判断を 2 回する手間よりはずっと安い**と感じました。

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
- [⑥ Security Rules が本番だけ壊れる（エミュレータでは再現しない話）](/engineering/app/firestore-security-rules-production-only-bugs/)

#### 仕様と設計の判断

- [⑦ 個人開発アプリの課金とプラン設計（上限・買い切り・解約戻り・キルスイッチ）](/engineering/app/indie-app-pricing-and-entitlement-design/)
- [⑧ iOS と Android を往復できるようにした（アカウントリンクと単一端末リース）](/engineering/app/cross-platform-account-migration-ios-android/)
- [⑨ 報告とブロックを後から入れた話（ストアのUGCポリシー対応）](/engineering/app/adding-ugc-report-block-for-store-policy/)
- [⑩ iPad 対応を段階的に入れた（2ペイン化の判断とやらなかったこと）](/engineering/app/ipad-two-pane-staged-support/)

#### 運用と検証

- [⑪ 「動いている」を証明する方法（実機・本番・計測の作法）](/engineering/app/proving-it-works-device-and-production-verification/)
- [⑫ iOS と Android を日英対応にした運用（String Catalog と values-en）](/engineering/app/ios-android-string-catalog-i18n/)
- [⑬ AI エージェントとひとりで両OSアプリを作った運用（仕様書・issue・構造ガード・記憶）](/engineering/app/solo-dev-with-ai-agent-workflow/)

### ストアに出すまでのシリーズ

- [① 個人開発でストアに登録する前の事務手続き（屋号・開業届・D-U-N-S）](/engineering/app/indie-app-store-registration-paperwork/)
- [② 個人開発アプリの周辺インフラを用意して踏んだ罠（サイト・メール・OAuth・バックエンド）](/engineering/app/indie-app-domain-site-email-oauth-setup/)
- [③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/app-store-connect-review-submission-notes/)
- [④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/google-play-console-review-submission-notes/)
- [⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程](/engineering/app/personal-app-cross-store-release-full-journey/)
