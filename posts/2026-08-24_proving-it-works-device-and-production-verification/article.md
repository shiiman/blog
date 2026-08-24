---
title: 「動いている」を証明する方法（実機・本番・計測の作法）
slug: proving-it-works-device-and-production-verification
date: '2026-08-24T17:00:00+09:00'
categories:
  - app
tags:
  - testing
  - indie-dev
  - ios
  - android
draft: false
id: 2258
excerpt: 個人開発アプリのリリース前に、実機と本番でしか出ない不具合をどう測ったかの記録です。HTTP 200が返ることは反映の証拠にならない、Storageの403は実体がなくても返る、実機の状態はplist回収で機械的に測れる、断定する前に代替軸を1つ潰す。同じ型の誤りを1セッションで7回踏んだ表も載せています。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** をリリースする前に、**実機と本番でしか出ない不具合**を何度も踏みました。

そのたびに困ったのは、不具合そのものより「**動いていることをどう証明するか**」でした。ビルドが通る、テストが green、画面に出ている、HTTP 200 が返る。どれも「動いている」の証明にはなりません。

この記事はその作法をまとめたものです。個別の不具合は [SwiftUI と SwiftData の罠](/engineering/app/ios-swiftui-swiftdata-pitfalls/) / [iOS の課金とストア連携](/engineering/app/ios-storekit-widget-pitfalls/) / [Android 版](/engineering/app/android-compose-room-billing-pitfalls/) / [Security Rules](/engineering/app/firestore-security-rules-production-only-bugs/) に分けて書きました。

---

## 証拠にならないものの一覧

まず、繰り返し騙されたものを並べます。

| 見えるもの | 証明していること | 証明していないこと |
|---|---|---|
| ビルドが通る | コンパイルできた | 起動できる |
| テストが green | テストが通った | 本番の経路が動く |
| HTTP 200 が返る | 何かが返った | 期待した内容が返った |
| 画面に「削除しました」と出る | UI が遷移した | サーバーのデータが消えた |
| `UPLOAD SUCCEEDED` | 転送が終わった | ストアが受理した |
| API が「成功」を返す | 呼び出しが通った | 状態が変わった |
| PR をマージした | コードが main に入った | 本番に反映された |
| エミュレータのテストが green | ロジックが正しい | 本番で通る |

これを毎回自分に言い聞かせるようにしました。**「◯◯だから動いている」と書きそうになったら、◯◯が本当にその証明になっているかを確認する。**

## デプロイ直後の 1 回目は旧版が 200 で返る

静的サイトの `assetlinks.json` を更新したときの実測です。

```
[試行 1] HTTP 200 / 430 bytes  ← 旧版（指紋 2 件）
[試行 2] HTTP 200 / 537 bytes  ← 新版（指紋 3 件）
```

ワークフローが success でも、直後の 1 回目は旧版が返ります。しかも **HTTP 200 なのでエラーとして気づけません**。

判定は「**期待する内容が来るまでリトライ**」にしました。`size_download` の変化も手がかりになります。**1 回で「未反映」と判断すると誤報になります。**

そして `curl` の使い方でも 1 回外しました。**`-L` を付けないと、末尾スラッシュへの 308 リダイレクトの空ボディを読んでしまいます**。これで「反映されていない」と誤認しました。

## Storage の 403 は実体がなくても返る

Cloud Storage の画像を消したことを確認しようとして、こう考えました。「消えたなら 404 になるはずだ」。

**なりません。** Cloud Storage は**未認証の呼び出し元に存在の有無を明かさない**ので、実体があってもなくても 403 を返します。

```
削除前: 403
削除後: 403   ← 区別できない
```

正しい判定は「**200 以外**（実際に返るのは `403 Permission denied.`）**かつ、オブジェクトが実際に削除されていること**」でした。後半は管理側の API で一覧を取って確認します。

ここでも一度誤断定しました。プレフィックスでクエリして 0 件だったので「実体は 0 件」と結論したのですが、**クエリの書き方が間違っていて、そもそも列挙できていませんでした**。

> **教訓**: 空の結果は「0 件である証明」ではなく「0 件を返した証明」です。
> 列挙の手段が正しく動いていることを、既知の 1 件で先に確かめます。

同じ型で、本番データの棚卸しでも **3 回誤報しました**。原因は毎回同じで、「推測でパスを叩いていた」ことでした。対処は簡単でした。**まずルートのコレクションを列挙してから、その一覧を回る**。推測でパスを組み立てない。

## 実機の状態は機械的に測れる

実機の検証を「目視で確認して報告してもらう」から「機械で測る」に変えられました。iOS の場合の手段が 3 つあります。

### 1. Release ビルドのアプリコンテナから plist を回収する

```sh
# 一覧
xcrun devicectl device info files --device <UDID> \
  --domain-type appDataContainer --domain-identifier <bundle-id>

# 取り出す
xcrun devicectl device copy from --device <UDID> \
  --domain-type appDataContainer --domain-identifier <bundle-id> \
  --source Library/Preferences/<bundle-id>.plist --destination "$PWD/prefs.plist"

plutil -p prefs.plist | grep -iE "account\.|terms\."
```

**Release ビルドでも読めます。** Xcode の自動署名が Development プロファイルで署名するので、そのための entitlement が立つためです。App Group 側のコンテナも読めます（`Library/Preferences/*.plist` だけ）。

ただし **SwiftData のストアは取れません**。App Group のルート直下にあるファイルは、パスの制約でコピーが拒否されます。エラーメッセージ自体は実体の場所を示しているので、**ファイルは存在していて、ツール側の制約**でした。ローカル DB の中身は plist と画面で代替しました。

### 2. 強制終了のタイミングを制御する

「同期の途中でアプリが kill されたらどうなるか」を測るには、狙ったタイミングで殺す必要があります。プロセスに SIGKILL を送れば再現できました。

### 3. サーバー側のポーリングを kill のトリガーにする

サーバーのデータが特定の状態になった瞬間にアプリを殺したい、というときは、Firestore をポーリングして条件が成立したら kill します。「操作の途中」を人間の反射で狙うのは無理なので、ここは自動化が必要でした。

### Android は環境で読めるものが変わる

| 環境 | 内部状態 |
|---|---|
| Play イメージ + release ビルド | **一切読めない**（`adb root` が拒否され、`run-as` も不可） |
| Google API イメージ + debug ビルド | 読める |

読める場所は `shared_prefs/` の下です。ただし **Firebase Auth のストアは暗号化されていて UID は読めません**。

そこで使ったのがこれです。**「書き換わった瞬間」は md5 で読めます。** サインインやサインアウトのたびにこのファイルが書き直されるので、**md5 の変化を「サインインが成立した瞬間」のトリガーに使えます**。これがあると「ある操作がセッションの窓の内側に着弾した」ことを、推定ではなく断定できました。

データを消す検証のカナリアには、`.clear()` でキーごと消える設定値を使います。既定値と違う値を入れておけば、UI でも prefs でも判定できます。

## 「Play 配信で検証した」は前提を確認するまで言えない

Android で一番効いた確認がこれでした。

```sh
adb -s <serial> shell dumpsys package <package> | grep -E 'versionCode|installerPackageName'
#  installerPackageName=com.android.vending  ← これが無ければ Play 配信ではない
```

エミュレータに入っていたビルドが、**`installerPackageName=null` = `adb install` で入れたサイドロード版**でした。Play 配信版は Play のアプリ署名鍵で署名されるので、ローカル鍵で入れた APK への上書き更新は構造的に不可能です。

しかも **Play ストアの画面には `Can't install` としか出ません**。logcat を見ないと原因に辿り着けません。

エミュレータのスナップショットが巻き戻ったあとは特に疑う必要があります。巻き戻り先がサイドロード時代なら、**以後の「Play 配信で検証した」は全部前提が崩れています**。

## 計測ハーネスを信じる条件

ウィジェットの追従を測ったとき、**計測ミスを 5 回踏みました**。そこから整理した条件が 2 つです。

### 1. 「何が起きたか」を推測せず、アプリ自身の表示と突き合わせる

最初は「トグルを押したから期待値は 1 のはず」と推測して測り、トグルの意味を取り違えて全滅しました。

正しくは **アプリのホームが出す達成数 A と、ウィジェットの値 W を比較**します。これなら操作の意味に依存しません。

### 2. OK / NG / 計測不能 を分けて数える

一緒にすると、**ハーネス自身の不具合を「不具合が再現した」と誤読します**。実際、正規表現が `text="0 / 24"` にマッチしないせいで 6 連続のタイムアウトが出て、アプリの不具合だと誤読しかけました。

### 踏んだ細かい罠

| 罠 |
|---|
| **前景復帰でウィジェットが更新される**ので、アプリに戻ってから読むと先に直ってしまう |
| 同期していない状態から測ると「2 秒で追従」と誤認する（期待値がたまたま現在値と一致する） |
| ハーネス内の戻る操作が画面遷移を壊して、以降の試行が全部無効になる |
| `adb shell am start --es ... "複数語"` は端末側のシェルで再分割され、先頭 1 語しか渡らない |
| **日本語 IME が `adb shell input text` を変換する**（メールアドレスが全角に化ける） |
| Compose の `clickable(enabled = false)` は uiautomator の `enabled` 属性に出ない（`clickable="false"` を見る） |

## 「症状が消えた」を成功と誤読しない

これが一番危ない型でした。

**「クラッシュしなくなった」は「ガードが効いた」証拠になりません。** 「機能が死んだ」でも同じ結果になります。

同じ話で、CI に置いた構造ガードにも問題がありました。**実バグを再発させても合格を出していたガードが 3 本ありました。** 原因は、検索する前にコメントを剥がしていなかったことです（コードを消してもコメントに当たって通ってしまう）。

**ガードを置いたら「壊して落ちること」を必ず測ります。** 11 パターンの負のテストを回して、3 本が検知できていないことが分かりました。

> **教訓**: 合格が出ても壊れを検知できないなら、そのガードの価値はゼロです。

## 断定する前に代替軸を 1 つ潰す

ある 1 セッションで、**同じ型の誤りを 7 回踏みました**。いずれも測るのに 1〜5 分しかかからないことを、仕様書や issue やコードのコメントに**書き切ってから**確かめています。訂正に PR を 1 本使いました。

| 断定した内容 | 実際 | 潰すべきだった軸 |
|---|---|---|
| in-memory だから削除が素通りした | ファイル永続でも成功した。真因は別 | **3 つ目の構成を試す** |
| いま退会すると権利を失う | 既にその権利は無かった | **自分が取ったログ** |
| この 3 件は実測不能と判断済み | 判断されていない | **issue のコメント本文** |
| デバッグのトグルを ON にした | 設定の書き込みが届いていない | **アプリが読む経路での読み出し** |
| 漏れは 2 個 | 9 個 | **機械で数える** |
| この条件で直る | 初期化の順序で無効になる | **修正後に「効くべき側」も測る** |
| サインアウトで解放される | 解放は全消去だけ | **該当関数を grep する** |

危険なのは 4 番目と 6 番目です。**どちらも「症状が消えた」ことを成功と誤読していました。**

そこでルールを 2 つ決めました。

1. **「A と B が違う → 原因は X だ」と書く前に、X 以外の軸を 1 つ潰す**
2. **仕様や状態を断定する前に、コードか記録を 1 か所読む**

どちらも数分で済みます。書いてから確かめると、訂正のコストが数十倍になります。

## 仮説の絞り方

不具合の原因を外し続けたときに効いたのは、**症状を全部同時に説明できるか**で絞ることでした。

Android のウィジェットが「常に 1 手前の内容を描く」不具合では、**5 回仮説を外しました**。最後に効いた判定はこうです。

```
症状 1: 常に 1 手前の値
症状 2: アプリを再起動すれば直る
症状 3: 更新の API 呼び出しは成功を返している
```

この 3 つを同時に説明できる仮説が 1 つだけ残りました（フレームワークがデータを読み直していない）。**症状を 1 つずつ説明する仮説は無限に作れます**が、全部を同時に説明できるものは少ないです。

もう 1 つ効いたのが**ログの非対称**でした。

```
アプリ起動時    : 更新 → データ読み込みのログが出る → 正しい値
データ変更時    : 更新 → データ読み込みのログが出ない → 1 手前の値
```

「出る / 出ない」の差が、そのまま原因を指していました。

## Debug でモックに落ちる経路は構造的に止める

最後に、一番効いた対策です。

**Debug では別の実装に差し替わる経路は、本番の経路が一度も実行されません。** クライアント保護のトークンが 3 か月間 1 件も発行されていなかったのも、Release の実機で初めて 2 件落ちたのも、全部これでした。

対策は 3 つ置きました。

1. **Release で値が無ければ即座に落とす。** 動かないまま出荷されるより、ビルドで止まるほうがマシです
2. **回帰のガードは挙動のテストではなく構造検査に置く。** 壊れた経路には、そもそもテストが到達できません
3. **Debug で確認できないことを一覧にしておく。** 「ここは実機の Release でしか測れない」と分かっていれば、リリース前にまとめて測れます

## まとめ

やってみて、証明の作法として残ったものを 5 つ。

1. **「◯◯だから動いている」の ◯◯ が本当に証明になっているかを確認する。** HTTP 200 も画面表示もマージも、証明ではありません
2. **空の結果は「0 件である証明」ではない。** 列挙の手段が動いていることを既知の 1 件で確かめます
3. **計測ハーネスは「推測した期待値」で組まない。** アプリ自身の表示と突き合わせます
4. **ガードを置いたら壊して落ちることを測る。** 合格が出ても検知できないなら価値はゼロです
5. **断定の前に代替軸を 1 つ潰す。** 数分で済み、訂正のコストは数十倍です

個人開発だと、間違いを指摘してくれる人がいません。**測り方を先に決めておくことが、レビュアーの代わりになる**と感じました。

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
- [④ iOS の課金・ウィジェット・App Check で踏んだ罠（Debug では一度も通らない経路）](/engineering/app/ios-storekit-widget-pitfalls/)
- [⑤ Android 版で踏んだ罠（Compose・Glance・Play Billing・エミュレータ検証）](/engineering/app/android-compose-room-billing-pitfalls/)
- [⑥ Security Rules が本番だけ壊れる（エミュレータでは再現しない話）](/engineering/app/firestore-security-rules-production-only-bugs/)
- [⑦ 個人開発アプリの課金とプラン設計（上限・買い切り・解約戻り・キルスイッチ）](/engineering/app/indie-app-pricing-and-entitlement-design/)
- [⑧ iOS と Android を往復できるようにした（アカウントリンクと単一端末リース）](/engineering/app/cross-platform-account-migration-ios-android/)
- [⑨ 報告とブロックを後から入れた話（ストアのUGCポリシー対応）](/engineering/app/adding-ugc-report-block-for-store-policy/)
- [⑩ iPad 対応を段階的に入れた（2ペイン化の判断とやらなかったこと）](/engineering/app/ipad-two-pane-staged-support/)
- [⑫ iOS と Android を日英対応にした運用（String Catalog と values-en）](/engineering/app/ios-android-string-catalog-i18n/)
- [⑬ AI エージェントとひとりで両OSアプリを作った運用（仕様書・issue・構造ガード・記憶）](/engineering/app/solo-dev-with-ai-agent-workflow/)

### ストアに出すまでのシリーズ

- [① 個人開発でストアに登録する前の事務手続き（屋号・開業届・D-U-N-S）](/engineering/app/indie-app-store-registration-paperwork/)
- [② 個人開発アプリの周辺インフラを用意して踏んだ罠（サイト・メール・OAuth・バックエンド）](/engineering/app/indie-app-domain-site-email-oauth-setup/)
- [③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/app-store-connect-review-submission-notes/)
- [④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/google-play-console-review-submission-notes/)
- [⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程](/engineering/app/personal-app-cross-store-release-full-journey/)
