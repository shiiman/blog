---
title: SwiftUI と SwiftData で踏んだ罠（落ちる・効かない・反映されない）
slug: ios-swiftui-swiftdata-pitfalls
date: '2026-08-24T13:00:00+09:00'
categories:
  - app
tags:
  - swiftui
  - swiftdata
  - swift
  - ios
draft: false
id: 2254
excerpt: 個人開発アプリをSwiftUIとSwiftDataで書いて踏んだ罠をまとめました。NavigationLinkが行内のボタンを吸う、Menuに.tintを足すとSIGSEGV、@Environmentの書き方を1文字間違えると表示した瞬間にSIGTRAP、scaledToFillが兄弟のレイアウトを壊す、.constantのfullScreenCoverが配下のsheetを全部止める、増分ビルドの型情報が古いままクラッシュする、など。すべてビルドとlintは通ります。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** を SwiftUI と SwiftData で書いて踏んだ罠をまとめます。

並べてみて共通していたのは、**ビルドも SwiftLint も通ること**でした。型は存在して補完も効くのに、その画面を実際に表示した瞬間だけ落ちる。あるいは落ちずに、ボタンが押せない・値が反映されないという形で出ます。**静的に気づけない**のが質でした。

> このシリーズの他の記事です。[アプリの概要](/engineering/app/indie-bucket-list-app-overview/) / [技術構成](/engineering/app/indie-app-technical-architecture/) / [iOS の課金とストア連携の罠](/engineering/app/ios-storekit-widget-pitfalls/) / [Android 版で踏んだ罠](/engineering/app/android-compose-room-billing-pitfalls/)

---

## 落ちる系

### Menu に `.tint` を足すと SIGSEGV する

アプリの標準ボタンを「薄い accent 塗りのカプセル」に揃えたくて、`.buttonStyle(.bordered)` に `.tint(theme.color(.accent))` を一括で足しました。**シミュレータが起動直後に落ちました。**

```
swift_retain
  → initializeWithCopy for Menu
  → FriendsView.listBody のタプルコピー中
```

`Button` では起きません。**`Menu` だけ**で再現しました。ViewBuilder のタプルの型が変わって、値ウィットネスのコピーで不正なアドレスを触っているように見えます。

対処は「一括で tint を足すなら Menu は除外する」にしました。Menu の見た目を変えたいなら label 側で背景を作ります。

ここで学んだのは、**見た目を変えたら起動して生存を確認する**必要があるということでした。ビルドが通ることは起動できることの証拠になりません。

```sh
# 起動して、クラッシュレポートが増えていないかと、生きているかを見る
ls ~/Library/Logs/DiagnosticReports/ | tail -3
xcrun simctl spawn <UDID> launchctl list | grep <bundle-id>
```

### `@Environment` の書き方を 1 文字間違えると表示した瞬間に SIGTRAP

新しいシートを追加したときに `@Environment(ThemeProvider.self)` と書いたら、**シートを開いた瞬間にクラッシュ**しました。

```
_assertionFailure
  → EnvironmentValues.subscript.getter
  ...
  → -[UISheetPresentationController presentationTransitionWillBegin]
```

原因は、このアプリがテーマを **EnvironmentKey**（`\.themeProvider`）で配っていて、`@Observable` なインスタンスを `.environment(_:)` で注入していなかったことです。

`ThemeProvider` は `@Observable` なので **`@Environment(ThemeProvider.self)` は文法として正しく、ビルドも lint も通ります**。落ちるのは「その View を表示したとき」だけです。

数えてみたら、**アプリ全体で EnvironmentKey 経由が 62 ファイル、Observable 直参照は自分が書いたその 1 か所だけ**でした。既存の書き方から外れていることに気づけていませんでした。

> **教訓**: 新しい View を足すときは、同じものを既存が何と書いているかを数えてから書く。
> `grep -rn 'Environment(' | wc -l` で主流を確認して、単独の書き方は疑う。

対処として、プロパティ宣言の行だけを見る構造ガードを CI に追加して凍結しました。ここでも 1 つ落とし穴があって、**素朴に文字列 grep すると「なぜ書いてはいけないか」を説明したコメントに当たって永久に落ちます**。

### 増分ビルドの型情報が古いままクラッシュする

これが一番タチが悪かったものです。

`@ViewBuilder` の `if / else` の**どれかの分岐の中で View の子要素数を変える**と（`HStack` の子を 2 個から 4 個に増やす、`Spacer` を足す、など）、増分ビルドが古い型情報を再利用してクラッシュします。

クラッシュレポートのレジスタが指していた型がこれでした。

```
HStack<TupleView<(
    ModifiedContent<Image, _ForegroundStyleModifier<Color>>,
    VStack<TupleView<(Text, Text)>>
)>>
```

これは「2 要素の HStack」です。でも実際のコードは `Spacer` と chevron を足して **4 要素**にしていました。古い型レイアウトで新しい実体を触って、不正なポインタを retain して落ちていました。

**シミュレータの消去やアプリの削除では治りません。** `xcodebuild clean` からのフルリビルドで初めて治りました。型情報はアプリのバンドルではなく `.dylib` に焼き付くので、端末側をリセットしても意味がないためです。

同じ症状は **main をマージした直後**にも出ました。切り分けの結果はこうでした。

| 状態 | 結果 |
|---|---|
| 自分のブランチ単体 | 通る |
| main 単体 | 通る |
| **マージ直後** | **クラッシュ** |
| ブランチを切り替えて全再ビルドした後の同じコード | 通る |

見分け方があります。**`initializeWithCopy for〈型〉` が不正なアドレスを読んでいる形は、型レイアウトの食い違いの典型**です。クラッシュ場所が「絶対に NULL 参照しないはずの SwiftUI 標準のコード」なら、ほぼキャッシュ汚染です。

| 変更の種類 | clean が必要か |
|---|---|
| HStack / VStack の子要素数を変えた | **必要** |
| `if let` / `if` の分岐を足した・消した | **必要** |
| `Group { }` のラップを足した・消した | **必要** |
| Gesture の modifier を足した | **必要** |
| 色を変えた・文字を変えた・`.foregroundStyle` を変えた | 不要 |

## 効かない系

### NavigationLink が行の中のボタンを吸う

リストの行に複数のボタン（達成チェック / 追加 / いいね）を置いていて、行全体をタップしたら詳細に飛ばしたい、という要件でした。

素直に `NavigationLink { Detail } label: { RowWithButtons }` にすると、**ラベル領域全体が NavigationLink のタップターゲットになって、中のボタンが反応しません**。「ハートが押せない」「追加ボタンが押せない」とユーザーから指摘されました。`.buttonStyle(.borderless)` を足しても完全には直りません。

行に内側のボタンがあるときは、NavigationLink を使わずにこうしました。

```swift
@State private var detailTarget: TargetType?   // Hashable が必要

ForEach(items) { item in
    BucketRowView(...)
        .contentShape(Rectangle())
        .onTapGesture { detailTarget = item }
}
.navigationDestination(item: $detailTarget) { item in
    DetailView(item: item)
}
```

### `.constant` の fullScreenCover が配下の sheet を全部止める

症状は「FAB の `+` を押しても何も出ない」でした。フラグは true になっていて、シートの body も構築されているのに、視覚的に出てきません。

真因は別の場所にありました。ルート側でスターターピッカーを `fullScreenCover(isPresented: .constant(shouldShowStarterPicker))` で出していたのですが、**`.constant` のバインディングを渡した fullScreenCover は、値が false でも提示ホストを恒久的に占有します**。その配下のシートが一切出せなくなっていました。

**`.sheet` は false のとき占有しないのに、`.fullScreenCover(.constant(...))` は占有する**という非対称でした。

対処は、スターターピッカーを fullScreenCover ではなく**全画面の if / else 分岐**にすることでした。挙動は同じで、ホストを占有しません。

`.constant()` を提示のバインディングに使うのはアンチパターンです。条件付きの提示は state のバインディングか分岐でやります。

切り分けの手順も残しておきます。シミュレータはプログラムからタップできないので工夫が必要でした。

- ディープリンクでフラグを立てて、UI 操作なしで再現させる
- DEBUG 限定で表示から 10 秒後にフラグを立てるトリガーを仕込み、ログで「フラグ到達」と「シートの body 構築」を観測する
- **`print` はログストリームに載らないので `NSLog` を使う**（ViewBuilder の中では `let _ = NSLog(...)`）
- 親のモーダルを 1 つずつ無効化して犯人を絞る

### iPad の 2 ペインは selection 駆動でないと詳細が切り替わらない

サイドバーの行を `Button` で叩いて外部の `@State` を変え、詳細を `NavigationStack { detail }.id(selected)` にする方式は**効きませんでした**。

NavigationSplitView は detail カラムをキャッシュして、外部の state 変化や `.id()` の変更を尊重しません。詳細を深く push した状態で別の項目を選ぶと、**左の選択ハイライトは動くのに右のペインが切り替わらない**という状態になります。

正解は NavigationSplitView 自前の selection 機構で駆動することでした。

- サイドバーは **regular のときだけ `List(selection: $selected)`**
- 行は `NavigationLink(value:)`（chevron・全幅タップ・選択連動を標準で提供してくれる）
- 詳細は `NavigationStack { detailColumn }`。`.id()` は不要

compact では selection なしの素の `List` にします。`.constant(nil)` を渡すと行が選択モードになって `NavigationLink` の push を殺すので、**`List` 自体を `horizontalSizeClass` で分岐**させました。

### `role: .destructive` の赤が祖先の `.tint` に負ける

削除ボタンが青くなっていました。`Button(role: .destructive)` と書いているのに、です。

原因はタブ全体に `.tint(theme.color(.accent))` を掛けていたことでした。そして**負ける文脈と負けない文脈があります**。実測した結果です。

| 文脈 | 実測した色 | 判定 |
|---|---|---|
| `swipeActions` 内の destructive | accent の青 | **負ける** |
| `.borderedProminent` + destructive | accent そのもの | **負ける** |
| Form / List 行の plain な destructive | システムの赤 | 負けない |
| `Menu` の中 / `.alert` の中 | 赤 | 負けない |
| `.onDelete` が生成する削除ボタン | 赤 | 負けない（UIKit 側が描くので tint を拾わない） |

「`role: .destructive` を書いたから赤」という思い込みが危険でした。負ける 2 つの文脈では `.tint(theme.color(.danger))` を明示する必要があります。

しかもこの誤解は**コメントとして固定化されていました**。Android 側のコードに「iOS の swipeActions は削除 = システムの赤」と書かれていて、実態と違う前提が残っていました。色は記憶で判断せず、スクリーンショットから画素を測って確定させるようにしました。

### Form 内で `simultaneousGesture` を足すとボタンが発火しない

キーボードを閉じたくて Form 全体に `.simultaneousGesture(TapGesture().onEnded { ... })` を付けたら、**Form 内のボタンをタップしてもアクションが発火しなくなりました**。TapGesture がボタンのタップを横取りします。

この組み合わせは使わず、`.scrollDismissesKeyboard(.immediately)` だけにしました。他のボタンや Picker をタップしたときは SwiftUI の標準挙動でフォーカスが外れるので、追加処理は不要でした。

編集モードの解除は `@FocusState` の `onChange` に寄せています。

```swift
.onChange(of: isTitleFocused) { _, focused in
    if !focused { isTitleEditing = false }
}
```

## レイアウトが壊れる系

### `scaledToFill` が兄弟のレイアウトを押し出す

カードの全面背景に画像を置いたら、**同じ ZStack にいるテキストが見えなくなりました**。

ユーザーからは「コントラストが足りない」「同化して見えない」と報告されたのですが、実態は違いました。`Image.resizable().scaledToFill()` を frame の制約なしに ZStack に置くと、**画像の自然サイズ（例: 800×600）をレイアウト提案として上流に伝えてしまい、兄弟の VStack がカードの外へ押し出されていました**。

`.frame(maxWidth: .infinity, maxHeight: .infinity).clipped()` を後置しても、親への提案を引きずる残効が消えません。

対処は、画像を直接 ZStack の子にせず `Color.clear.overlay(image)` でラップすることでした。

```swift
Color.clear
    .overlay(
        image.resizable().scaledToFill()
    )
    .clipped()
```

`Color.clear` は自然サイズを持たないので親の提案に従い、その overlay として画像が強制的に親の領域に収まります。

> **教訓**: 「コントラストが足りない」と言われたら、まずレイアウトの overflow を疑う。
> ビルドキャッシュを消しても再現するので、前述の clean ビルド問題とは別物です。

## 反映されない系

### SwiftData の配列プロパティは init 経由で渡す

Preview やデバッグメニュー用に `@Model` のインスタンスを ModelContext に insert せず standalone で作るとき、**post-init での配列プロパティの代入が読み出し時に空配列を返す**ことがありました。

```swift
// 反映されないことがある
let bucket = BucketItem(title: ..., categoryID: ...)
bucket.coverImageDataList = [data]

// 確実に反映される
let bucket = BucketItem(title: ..., categoryID: ..., coverImageDataList: [data])
```

`@Model` は内部で backing data 経由でプロパティを管理していて、ModelContext に紐づいていないインスタンスでは特定の型（配列や Optional）の post-init での書き込みが効かないケースがあるようです。「写真があるのに表示されない」という症状で出ました。

`status` のような単一の値は post-init の代入でも反映されるので、余計に気づきにくいです。**配列系は init 経由**に統一しました。

### テストで `ModelContainer` を保持しないと SIGTRAP

これはテスト側の罠です。

```swift
// クラッシュする
private func makeContext() throws -> ModelContext {
    try AppSchema.makePreviewContainer().mainContext  // container が即解放される
}

// 正しい
let container = try AppSchema.makePreviewContainer()
let context = container.mainContext   // container をテストのスコープで保持
```

`ModelContext` はコンテナを強参照で持たないので、生成元の `ModelContainer` が解放されると context が無効になり、ストアに触れた瞬間に `EXC_BREAKPOINT (SIGTRAP)` で落ちます。

コンパイルは通ります。しかも**アサーション失敗ではなくクラッシュでスイートが落ちる**ので、テスト名を見ても原因が分かりません。ヘルパーから `.mainContext` だけ返さない、というルールにしました。

## 画像のメモリ対策はバイト数では止まらない

レビューで「入力データが 50MB を超えたら弾く」という対策を提案されたのですが、これは二重に的外れでした。

1. 普通の写真（1〜10MB）ではまず引っかからないので緩すぎる
2. そして**高圧縮画像は数 MB でもデコードすると巨大なビットマップに展開されます**（8000×8000 = 64 メガピクセル ≈ RGBA で 256MB）。バイト数では止められません

正しい対策は**ダウンサンプリング**でした。`CGImageSourceCreateThumbnailAtIndex` に `kCGImageSourceThumbnailMaxPixelSize` を渡すと、**デコードと縮小を一発でやってフル解像度のビットマップを一切メモリに載せません**。入力が何ピクセルでも、メモリ使用量は出力サイズに比例して有界になります。

バイトの上限も併設はしました。ただし主たる対策ではなく、粗いサニティチェックとして 20MB です。そして**エラーには超過した上限を MB で明示**します（「画像が大きすぎます（20MB まで）」）。書かないとユーザーは何で引っかかったのか分かりません。

> **教訓**: サイズ系の指摘が「バイト数チェック」を提案してきたら、それが本当にメモリ枯渇を防げるか確認する。
> 寸法ベースの対策を優先します。

## まとめ

SwiftUI と SwiftData で踏んだ罠に共通していたのは 3 つでした。

1. **ビルドと lint は通る。** 型は存在して補完も効くので、書いた時点では何も警告されません
2. **落ちるのは「その画面を実際に表示したとき」だけ。** 画面に到達しないと分かりません
3. **見た目の症状と原因が一致しない。** 「コントラストが足りない」の真因がレイアウトの overflow、「ボタンが反応しない」の真因が親のモーダル、というふうにズレます

やり方として効いたのは、**既存が同じことを何と書いているか数えてから書く**ことでした。単独の書き方をしているときは疑うようにして、確定したものは構造ガードで凍結しています。

課金とストア連携の罠は[こちら](/engineering/app/ios-storekit-widget-pitfalls/)、Android 版の罠は[こちら](/engineering/app/android-compose-room-billing-pitfalls/)に書きました。

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

<a class="link-card" href="/engineering/app/ios-android-string-catalog-i18n/">
<span class="link-card__thumb"><img src="/eyecatch/ios-android-string-catalog-i18n.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑧ iOS と Android を日英対応にした運用（String Catalog と values-en）</span>
<span class="link-card__excerpt">キーが日本語原文なので検出スクリプトはゼロにならない、カタログが 3 つある、文言の一致は機械照合しないと目視では見つからない。</span>
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
