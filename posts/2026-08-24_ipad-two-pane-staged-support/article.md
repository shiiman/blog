---
title: iPad 対応を段階的に入れた（2ペイン化の判断とやらなかったこと）
slug: ipad-two-pane-staged-support
date: '2026-08-24T14:00:00+09:00'
categories:
  - app
tags:
  - ios
  - swiftui
  - indie-dev
  - app-development
draft: false
id: 2264
excerpt: iPhone 向けに作ったアプリを iPad に対応させるとき、幅が広くなるだけで済むわけではありません。個人開発アプリで2ペイン化を4段階に分けて入れた記録です。最大幅の制限を入れて後から取り下げた話、汎用の2ペインコンテナを1つ作って4タブに適用した話、そして horizontalSizeClass が使えなかった理由を書きました。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** は iPhone 向けに作りました。ただ iPad でも動くので、そのままだと**幅いっぱいに 1 カラムのリストが伸びた画面**になります。

これを段階的に直しました。順番に入れて、途中で 1 つ取り下げています。この記事はその判断の記録です。

> シリーズ「アプリの中身と実装」全 12 本の 10 本目です。最初から読むなら [① アプリの概要](/engineering/app/indie-bucket-list-app-overview/) から。ストアに出すまでの手続きは[別シリーズ](/engineering/app/personal-app-cross-store-release-full-journey/)にまとめました。

---

## 4 段階に分けた

まず全体像です。

| 段階 | やったこと | 結果 |
|---|---|---|
| A | 読みやすい最大幅で中央寄せ | **後から取り下げた** |
| B | リストタブだけ 2 ペイン化 | 入れた |
| C | 残り 3 タブも 2 ペイン化 | 入れた |
| D | 横向きを許可 | 入れた |

小さく入れて確かめる形にしたのは、**iPad をメインの想定端末にしていない**からです。大きく作り替えて iPhone 側が壊れると、被害のほうが大きくなります。

## 段階 A: 最大幅の制限を入れて、取り下げた

最初にやったのが「**幅を 700pt で止めて中央に寄せる**」でした。iPad の広い画面で 1 行が横に伸びすぎるのを防ぐ、よくある手法です。

入れてみて、そのあと**取り下げました**。

理由は段階 B と C で 2 ペインを入れたことです。2 ペインにすると、右側の詳細ペインはもともと適度な幅に収まります。そこにさらに最大幅の制限がかかると、**ペインの中でさらに中央に寄って、左右に無駄な余白ができます**。

つまり A と B/C は**同じ問題への別解**でした。両方入れると効果が重複して、見た目が悪くなります。

### 呼び出し箇所だけ残した

取り下げるときに、**関数の呼び出しは全部残しました**。中身を何もしない実装に変えただけです。

```
readableContentWidth()  ← 呼び出しは各タブに残っている。実装は no-op
```

こうしたのは、幅の制限を復活させたくなったときに**1 か所を直せば全部に効く**ようにしたかったからです。呼び出しを消してしまうと、次に入れるときにまた全タブを触ることになります。

「やめた機能のフックだけ残す」は普段はやらないんですが、**判断が揺れそうなところ**だったのでこうしました。

## 段階 B / C: 2 ペイン化

### どのタブに入れたか

| タブ | 2 ペイン | 理由 |
|---|---|---|
| リスト | 入れた | リスト → 詳細の構造がある |
| 設定 | 入れた | 同じ |
| フレンド | 入れた | 同じ |
| 探す | 入れた | 同じ |
| **ホーム** | **入れない** | **ダッシュボードなので、詳細に相当するものが無い** |

判断の基準は「**リスト → 詳細の構造があるか**」の一点でした。

ホームは進捗のカードや今日の予定を並べたダッシュボードです。左に一覧、右に詳細という形にする対象がありません。ここを無理に 2 ペインにすると、**左側に何を出すかを新しく考える**必要が出ます。それは iPad 対応ではなく新機能です。

なので**幅いっぱいの 1 カラムのまま**にしました。

### 汎用のコンテナを 1 つ作った

段階 B ではリストタブに直接書きました。段階 C で残り 3 タブに広げるときに、**汎用のコンテナに切り出しました**。

```swift
SplitNavigation<Selection>  // 選択の型だけをパラメータにした
```

タブごとに選択する対象の型が違うので（バケット / 設定の項目 / フレンド / 公開リスト）、そこをジェネリクスにしました。

**4 タブぶんの 2 ペイン実装を 1 つにまとめた**のがここの効果でした。分岐の条件や選択が消えたときの復帰処理を 4 回書くと、必ずどこかが違う挙動になります。

### horizontalSizeClass が使えなかった

ここでハマりました。

2 ペインを出すかどうかは「**iPad で幅が広いとき**」で判定したい。素直に書けば `horizontalSizeClass == .regular` です。

**これが動きません。**

原因は、**サイドバーの中に入ると horizontalSizeClass が compact に上書きされる**ことでした。左ペインの中身は狭いので、SwiftUI が compact だと報告してきます。すると「2 ペインを出すか」の判定が、2 ペインの中で false になります。

`UIDevice.current.userInterfaceIdiom` で端末の種類を直接見る形に変えました。

```
判定に使うもの: 端末が iPad かどうか（+ 設定のトグル）
判定に使わないもの: horizontalSizeClass
```

**環境から取れる値が、階層のどこで読むかによって変わる**ケースでした。レイアウトの分岐に使う値は、階層に依存しないものを選ぶ必要があります。

### 設定でオフにできるようにした

2 ペインは全員が好むレイアウトではありません。特に**リストをできるだけ広く見たい人**には邪魔になります。

設定に「2 ペイン表示（iPad）」のトグルを置きました。既定はオンです。

個人開発でトグルを増やすのは基本的に避けたいんですが、ここは**好みが分かれる見た目の変更**なので入れました。判断の基準はこうしています。

- **正しい答えが 1 つに決まるもの** → トグルにしない
- **好みが分かれるもので、両方の実装コストが同じ** → トグルにする

2 ペインをオフにすると、もともとの 1 カラムの実装に戻るだけです。追加のコードはほぼ要りませんでした。

## 段階 D: 横向きを許可した

ここまで**アプリ全体が縦向き固定**でした。

そのため、2 ペインが**iPad の縦向きでしか全画面表示されません**。iPad を横に持つ人は多いので、これは実質的に半分しか効いていない状態でした。

そこで向きの許可を OS ごとに分けました。

|  | 許可する向き |
|---|---|
| iPhone | **縦のみ**（変更なし） |
| iPad | **4 方向すべて** |

iPhone を縦固定のままにしたのは、iPhone を横にして使う場面が想定になかったからです。横向きを許可すると、**全画面のレイアウトを横でも確認する必要が出ます**。効果に対してコストが合いません。

iPad だけ許可する設定は、プロジェクトの生成設定に書いています（このプロジェクトは `.xcodeproj` を生成する構成なので、Info.plist を直接編集しません。詳しくは[技術構成](/engineering/app/indie-app-technical-architecture/)に書きました）。

## やらなかったこと

iPad 対応の話でよく出てくるもののうち、入れなかったものを書いておきます。

|  | 判断 |
|---|---|
| Split View / Slide Over の最適化 | やらない（マルチタスクは OS 任せ） |
| ポインタとキーボードの対応 | やらない |
| Mac Catalyst | やらない |
| ホームタブの 2 ペイン化 | やらない（前述） |
| 最大幅の制限 | 入れたが取り下げた |

どれも「**iPad をメインの想定端末にしていない**」という一点から降りてきています。

個人開発だと、**端末ごとに最適化する労力は端末ごとに掛け算で増えます**。iPad を「使えるようにする」のと「iPad 向けに作る」の間には、大きな差があります。前者で止めました。

## まとめ

iPad 対応を段階的に入れて残ったことを 5 つ。

1. **段階に分けると、途中で取り下げられる。** 最大幅の制限は、2 ペインを入れた時点で不要になりました
2. **やめた機能のフックだけ残すのは、判断が揺れそうなところでは有効。** 復活させるときに 1 か所で済みます
3. **2 ペインにするかは「リスト → 詳細の構造があるか」で決まる。** ダッシュボードには無理に入れません
4. **汎用のコンテナに切り出すと、4 回ぶんの実装が 1 つになる。** 選択の型だけをパラメータにしました
5. **レイアウトの分岐に horizontalSizeClass を使えない場面がある。** サイドバーの中では compact に上書きされます

そして一番大きい判断は、**「iPad で使える」で止めた**ことでした。ここを線引きしないと、対応する端末が増えるたびに全画面を作り直すことになります。

iOS 側で踏んだ他の罠は [SwiftUI と SwiftData の記事](/engineering/app/ios-swiftui-swiftdata-pitfalls/)に書きました。

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
- [⑥ Security Rules が本番だけ壊れる（エミュレータでは再現しない話）](/engineering/app/firestore-security-rules-production-only-bugs/)

#### 仕様と設計の判断

- [⑦ 個人開発アプリの課金とプラン設計（上限・買い切り・解約戻り・キルスイッチ）](/engineering/app/indie-app-pricing-and-entitlement-design/)
- [⑧ iOS と Android を往復できるようにした（アカウントリンクと単一端末リース）](/engineering/app/cross-platform-account-migration-ios-android/)
- [⑨ 報告とブロックを後から入れた話（ストアのUGCポリシー対応）](/engineering/app/adding-ugc-report-block-for-store-policy/)

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
