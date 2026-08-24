---
title: iOS と Android を日英対応にした運用（String Catalog と values-en）
slug: ios-android-string-catalog-i18n
date: '2026-08-24T18:00:00+09:00'
categories:
  - app
tags:
  - i18n
  - swift
  - kotlin
  - indie-dev
draft: false
id: 2259
excerpt: 個人開発アプリをiOSとAndroidの両方で日英対応にした運用の記録です。キーが日本語原文なのでハードコード検出スクリプトは設計上ゼロにならない、xcodebuildはカタログを更新しない、カタログが3つある、values-enを追加すると既存テストが落ちる、そしてiOSとAndroidの文言一致は機械照合しないと目視では見つからない話をまとめました。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** を、iOS と Android の両方で日英対応にしました。配信地域が日本と米国なので、英語版が必要でした。

やってみて分かったのは、**翻訳そのものより運用のほうが厄介**だということでした。とくに「何をもって未翻訳ゼロと言うか」を決めるところで一度つまずいています。

> このシリーズの他の記事です。[アプリの概要](/engineering/app/indie-bucket-list-app-overview/) / [技術構成](/engineering/app/indie-app-technical-architecture/) / [SwiftUI と SwiftData の罠](/engineering/app/ios-swiftui-swiftdata-pitfalls/) / [Android 版で踏んだ罠](/engineering/app/android-compose-room-billing-pitfalls/)

---

## 方針: アプリ内に言語切り替えを作らない

最初に決めたのがこれです。**アプリ内に言語切り替えの UI は作りません。** OS の言語設定に従います。

理由は 2 つです。

1. OS の設定と別に持つと、2 つの真実ができる
2. 切り替え UI を作ると、再起動なしで全画面を再描画する仕組みが必要になる

英語対応そのものは実装しますが、切り替えは OS に任せる、という整理です。

## iOS: キーが日本語の原文になる

String Catalog（`.xcstrings`）を使いました。ここで一番大きな性質がこれです。

**キーが日本語の原文になります。**

```swift
Text("保存")                    // キーは "保存"
String(localized: "保存しました")  // キーは "保存しました"
```

つまり**ソースコードの中には日本語が残り続けます**。これが運用上の落とし穴になりました。

### ハードコード検出スクリプトは設計上ゼロにならない

「ソースに日本語が残っていたら未ローカライズだ」という発想で検出スクリプトを書いていました。**これは指標になりません。**

私のところでは、ローカライズ着手前が 2,003 件で、作業後が **2,013 件に増えました**。文字列補間の分割などで、むしろ増えることがあります。

正しい検証は「ソースの日本語の数」ではなく、**カタログ側**を見ることでした。

| 見るもの | 判定 |
|---|---|
| ソースの日本語の数 | **指標にならない**（キーが日本語なので減らない） |
| カタログの `en` が欠落しているキーの数 | **これが 0 なら未翻訳ゼロ** |
| カタログの `state: new` の数 | 抽出されたが未対応のもの |

検証スクリプトはこちらを見るように書き直しました。

### `xcodebuild` はカタログを更新しない

もう 1 つ、CLI 運用で必ず踏む罠です。

**文字列の抽出（`.stringsdata` の生成）にはビルド設定が必要です。**

```yaml
# project.yml
settings:
  base:
    SWIFT_EMIT_LOC_STRINGS: YES
```

これが無いと抽出そのものが走らず、カタログにキーが流れ込みません。

そして `xcodebuild`（CLI）は、**ソースの `.xcstrings` を更新しません**。これは Xcode の IDE 限定の挙動です。CLI で運用するなら、ビルド後に明示的に同期します。

```sh
xcrun xcstringstool sync <catalog>.xcstrings \
  --stringsdata <DerivedData>/.../<Target>.build/Objects-normal/arm64/*.stringsdata
```

`xcstringstool` には他に `compile`（検証と `en.lproj` の生成確認）、`print`、`extract` があります。

### カタログは 3 つある

ターゲットごとにバンドルが分かれるので、カタログも分かれます。

| カタログ | 対象 |
|---|---|
| `Resources/Localizable.xcstrings` | アプリ本体 |
| `Widgets/Localizable.xcstrings` | ウィジェット |
| `Features/QuickAdd/Localizable.xcstrings` | App Intents（ショートカット） |

同じ名前のカタログが複数あるとビルドが「複数のコマンドが同じものを生成する」で落ちるので、アプリ本体のソースのグロブから拡張側のカタログを除外しています。

**ここが後で効きました。** ウィジェット専用の文字列は本体のカタログに無いので、**アプリ内のデバッグ用ギャラリーで見ると日本語のまま表示されます**（実際のウィジェットは自分のカタログを見るので英語になる）。「英語が出ていない」と焦って調べる時間を使いました。

### 使い分けのルール

`LocalizedStringKey` が自動で効く場所と、明示的に書く場所を分けました。

| 文脈 | 書き方 |
|---|---|
| `Text` / `Label` / `Button` / `Section` / `.navigationTitle` / `.alert` / `.accessibilityLabel` のリテラル | そのまま書く（自動で抽出される） |
| `String` 型の文脈（`return "..."` / 補間 / `LocalizedError` / `String` 引数） | `String(localized:)` にする |
| 翻訳しない値 | `Text(verbatim:)` |

### 共有コンポーネントの引数型は変えない

これは実際に破綻させかけました。

18 ファイルから共有されているコンポーネントの `title` 引数を `String` から `LocalizedStringKey` に変えようとしたら、**全部の呼び出し側が壊れます**。

対処は、**`String?` のままにして、呼び出し側で `String(localized:)` に統一する**ことでした。共有コンポーネントの公開引数の型は変えない、というルールにしました。

## Android: `values-en` を追加すると既存テストが落ちる

Android 側で最初に踏んだのがこれです。

`values-en` を新設したら、**それまで通っていた Robolectric のテストが落ちました**。

追加前は、英語のロケールでも日本語のリソースしか無いので日本語に解決されていました。追加後は Robolectric の既定ロケール（ほぼ英語）が英語を解決するので、`getString` の日本語の本文を `contains("3 年前")` のように検証していたテストが落ちます。

修正は日本語のコンテキストに固定することでした。

```kotlin
val jaCtx = ctx.createConfigurationContext(
    Configuration(ctx.resources.configuration).apply { setLocale(Locale.JAPANESE) }
)
```

`values-en` を追加する前に、**翻訳キーを `getString` しているテストを洗い出します**。

```sh
grep -rn "R.string.<key>" app/src/test/
```

### lint が止めてくるもの

`values-en` を作ると、既定（日本語）の translatable なキー全部に英語が必要になります。無いと `MissingTranslation` でエラーです。`translatable="false"` にしたもの（クライアント ID など）は不要です。

もう 1 つ、Compose の中で `LocalContext.current.getString(...)` を呼ぶと **`LocalContextGetResourceValueCall` の lint に止められます**（error 扱い）。`val ctx = LocalContext.current` を経由しても同じです。

対処は、Composable のスコープで `stringResource` を先に解決してから lambda に渡すことでした。

```kotlin
// Composable のスコープで解決してから渡す
val message = stringResource(R.string.saved)
LaunchedEffect(Unit) { showToast(message) }
```

動的な引数があるときは、テンプレートを取って `String.format` します。

### 純 Kotlin の層に文字列を注入する

計算ロジックを Android のフレームワークに依存させたくないので、`Context` を持たない純粋な層に置いていました。そこに文言が必要になったときの形がこれです。

**解決済みのバンドルを、構築側で作って渡します。**

```kotlin
// Composable / DI コンテナ側で resources から解決
val strings = ExportStrings.from(LocalContext.current.resources)
// ViewModel には () -> ExportStrings のプロバイダとして注入
```

これで整形処理は純粋な関数のまま保てて、テストからは実際の `strings.xml` を読んで検証できます。**ハードコードした期待値ではなく実リソースを読むのが要点**でした。

数量に依存する文字列は `<plurals>` と `getQuantityString` にして、同じくバンドルに `(Int) -> String` のラムダとして持たせました。日本語は `other` 1 件、英語は `one` / `other` の 2 件です。

## 日付の書式を iOS と揃える

地味に探しました。

| 使うもの | 結果 |
|---|---|
| `DateFormat.getBestDateTimePattern(locale, "yMd")` | 英語は `5/18/2026` / 日本語は `2026/5/18`。**iOS と一致** |
| `DateTimeFormatter.ofLocalizedDate(SHORT)` | 英語が 2 桁年（`5/18/26`）になって iOS と不一致 |

ICU のスケルトンの単文字 `M` / `d` が非ゼロ埋めのロケール順を返すので、iOS の `.month(.defaultDigits).day(.defaultDigits)` と揃います。

パリティのテストではタイムゾーンを UTC に固定して、日付の暦日を決定論的にしました。

## 文言の一致は機械照合しかない

ここが一番効きました。

**iOS と Android で同じ文言になっているかを目視で確認するのは無理です。** 実際に見つかった違いがこれです。

| 見た目 | 実体 |
|---|---|
| `–` | U+2013 en ダッシュ |
| `〜` | U+301C 波ダッシュ |
| `−` | U+2212 マイナス |
| `~` | U+007E ASCII チルダ |

これらの取り違えは目では見えません。

そこで、Android の XML と iOS のカタログを機械的に読んで、**書式指定子の方言を正規化してから逐語比較**するスクリプトを書きました。

```
%lld / %1$lld  ↔  %d / %1$d
%@             ↔  %s
```

### 日本語側も正規化しないと見逃す

ここで 1 回外しました。最初は英語だけを比較していたのですが、**日本語のキーも指定子を正規化しないと照合できません**。

Android の日本語が `%1$d`、iOS の日本語が `%lld` だと、**同じ文なのに完全一致しません**。すると「Android 固有の文言」と誤判定されて、対応する英語のずれを見逃します。

正規化してから引き直したら、検出される不一致が **5 件から 7 件に増えました**。

### 日本語でも既に割れていた

この作業で炙り出たのが、**英語以前に日本語で iOS と Android が割れていた箇所**でした。

- シート名のスペースの有無
- 割合の表記（`(3/4)` の形）
- 状態のラベル（「進行中」「未達成」の使い分け）

過去のレビューを 5 巡していて、ロジックの一致は見ていました。でも**文字列の完全一致までは追えていませんでした**。英語対応をやると必ず出てくるので、そこで日本語ごと iOS 側に揃えました。

### 計算層に解決済みのラベルを注入する

もう 1 つ、分割ミスをしました。

統計の区間ラベルのように「UI とエクスポートの両方で共有される計算結果の文字列」を、エクスポート側だけ英語化してしまい、**UI が日本語のまま割れました**。

対処は、**計算層に解決済みのラベルを注入する**ことでした。そうすれば両方が同時に追従します。

## テストとリリースの前に測ったもの

最後に、リリース前の確認手段です。

### iOS でシミュレータを英語で起動する

```sh
xcrun simctl launch <sim> <bundle-id> -AppleLanguages "(en)" -AppleLocale "en_US"
```

バンドルに `en.lproj/Localizable.strings` が生成されていることを合わせて確認します。

ただし**新規インストールだとオンボーディングから始まるので、タブバーまで到達できません**（シミュレータはプログラムからタップできない）。深い画面は、英語のバンドルを直接読んで確認しました。

```sh
plutil -convert json en.lproj/Localizable.strings
```

### 英語で出やすい崩れ

**英語のほうが文字列が長くなる**ので、レイアウトが崩れる箇所が出ます。

| 崩れた場所 | 対処 |
|---|---|
| ウィジェットの小サイズのフッターが 2 行になってタイトルを圧迫 | `lineLimit(1)` + 短い英訳 |
| ウィジェットのヘッダー（`Completed N/M`）が省略される | 達成数に `layoutPriority(1)` + `minimumScaleFactor` |

ここは**英語で 1 回全画面を見る**しかありませんでした。

### iOS のカタログを手で編集するとき

新しいキーを既存のカタログへ手で挿入したいことがあります。そのときは **Python でラウンドトリップさせるのが安全**でした。

```python
json.dumps(catalog, ensure_ascii=False, indent=2) + "\n"
```

この形式が Xcode の出力とバイト単位で一致することを確認しました。キーの順序は Xcode 独自のソートなので、自然な隣接キーの直前に挿入すれば、差分は新規エントリだけで済みます。

## まとめ

日英対応をやって残ったものを 5 つ。

1. **キーが日本語原文なので、ソースの日本語の数は指標にならない。** カタログの `en` 欠落を数えます
2. **CLI ではカタログが自動更新されない。** ビルド設定と明示的な同期の 2 つが必要です
3. **`values-en` の追加は既存テストを壊す。** 日本語フォールバックに依存していたテストを先に洗い出します
4. **文言の一致は機械照合しかない。** ダッシュ 1 文字の違いは目で見つかりません
5. **英語対応をやると日本語のずれが炙り出る。** ロジックのレビューでは文字列の完全一致まで追えていませんでした

そして一番効いたのは、**「未翻訳ゼロ」の定義を先に決めること**でした。「ソースに日本語が無い」を目標にしていたら、永遠に終わりませんでした。

## 関連記事

**このシリーズ（アプリの中身と実装）**

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

<a class="link-card" href="/engineering/app/ios-swiftui-swiftdata-pitfalls/">
<span class="link-card__thumb"><img src="/eyecatch/ios-swiftui-swiftdata-pitfalls.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">③ SwiftUI と SwiftData で踏んだ罠（落ちる・効かない・反映されない）</span>
<span class="link-card__excerpt">Menu に .tint で SIGSEGV、@Environment の書き方で SIGTRAP、NavigationLink が行内ボタンを吸う、増分ビルドの型情報が古いまま落ちる。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

<a class="link-card" href="/engineering/app/ios-storekit-widget-pitfalls/">
<span class="link-card__thumb"><img src="/eyecatch/ios-storekit-widget-pitfalls.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">④ iOS の課金・ウィジェット・App Check で踏んだ罠（Debug では一度も通らない経路）</span>
<span class="link-card__excerpt">App Check が 3 か月間トークンを出していなかった話、Release の実機で初めて落ちた 2 件、アップロードが必ず失敗していた話。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

<a class="link-card" href="/engineering/app/android-compose-room-billing-pitfalls/">
<span class="link-card__thumb"><img src="/eyecatch/android-compose-room-billing-pitfalls.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑤ Android 版で踏んだ罠（Compose・Glance・Play Billing・エミュレータ検証）</span>
<span class="link-card__excerpt">Glance がセッション中データを読み直さない、キーボードでウィンドウがリサイズされない、Play 配信で検証したつもりがサイドロード版だった。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

<a class="link-card" href="/engineering/app/firestore-security-rules-production-only-bugs/">
<span class="link-card__thumb"><img src="/eyecatch/firestore-security-rules-production-only-bugs.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑥ Security Rules が本番だけ壊れる（エミュレータでは再現しない話）</span>
<span class="link-card__excerpt">request.resource は create / update だけ。372 件のエミュレータテストが green のまま本番のバックアップ復元が一度も成功していませんでした。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

<a class="link-card" href="/engineering/app/proving-it-works-device-and-production-verification/">
<span class="link-card__thumb"><img src="/eyecatch/proving-it-works-device-and-production-verification.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑦ 「動いている」を証明する方法（実機・本番・計測の作法）</span>
<span class="link-card__excerpt">HTTP 200 は反映の証拠にならない、Storage の 403 は実体がなくても返る、実機の状態は plist 回収で測れる、断定の前に代替軸を 1 つ潰す。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

<a class="link-card" href="/engineering/app/solo-dev-with-ai-agent-workflow/">
<span class="link-card__thumb"><img src="/eyecatch/solo-dev-with-ai-agent-workflow.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑨ AI エージェントとひとりで両OSアプリを作った運用（仕様書・issue・構造ガード・記憶）</span>
<span class="link-card__excerpt">仕様書 15 章を正典にする、1 issue = 1 worktree = 1 PR、構造ガード 30 本、そして「任せると壊れる場所」の型。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

### もうひとつのシリーズ（ストアに出すまで）

- [① 個人開発でストアに登録する前の事務手続き](/engineering/app/indie-app-store-registration-paperwork/)
- [② 個人開発アプリの周辺インフラを用意して踏んだ罠](/engineering/app/indie-app-domain-site-email-oauth-setup/)
- [③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/app-store-connect-review-submission-notes/)
- [④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/google-play-console-review-submission-notes/)
- [⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程](/engineering/app/personal-app-cross-store-release-full-journey/)
