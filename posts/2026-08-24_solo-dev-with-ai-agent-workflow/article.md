---
title: AI エージェントとひとりで両OSアプリを作った運用（仕様書・issue・構造ガード・記憶）
slug: solo-dev-with-ai-agent-workflow
date: '2026-08-24T19:00:00+09:00'
categories:
  - app
  - ai
tags:
  - claude-code
  - indie-dev
  - ai
  - testing
draft: false
id: 2260
excerpt: 個人開発でiOSとAndroidの両方を約4か月・3,585コミットで作るあいだ、AIコーディングエージェントとどう分担したかの運用記録です。仕様書15章を正典にする、1 issue = 1 worktree = 1 PR、構造ガードを30本以上置く、判断の誤りを記憶に残す、そして「エージェントに任せると壊れる場所」の型をまとめました。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発のバケットリストアプリ **Bucket List Commit** を、iOS と Android の両方で約 4 か月・3,585 コミット・24 万行書きました。ひとりです。ただし**実装のほとんどは AI コーディングエージェントに書かせています**。

この記事は、そのあいだに固まった運用の記録です。プロンプトの書き方の話ではなく、**何を人間側の資産として持っておくと壊れないか**の話が中心になります。

> このシリーズの他の記事です。[アプリの概要](/engineering/app/indie-bucket-list-app-overview/) / [技術構成](/engineering/app/indie-app-technical-architecture/) / [「動いている」を証明する方法](/engineering/app/proving-it-works-device-and-production-verification/)

---

## 規模の前提

| | |
|---|---|
| 期間 | 約 4 か月 |
| コミット | 3,585 |
| コード | 1,619 ファイル / 246,720 行（テストが 36%） |
| issue と PR の通し番号 | 2,189 まで |
| 仕様書 | 15 章 / 約 40 万字 |
| 構造ガード | 30 本以上 |

この規模になると、**エージェントとの会話の中身は残りません**。文脈は流れていくので、残るのはリポジトリに書いたものだけです。だから「どこに何を書いておくか」が運用の本体になりました。

## 仕様書を正典にする

一番効いたのがこれでした。**先に仕様書を書いて、それを唯一の正典にします。**

15 章に分けました。

| 章 | 内容 |
|---|---|
| 01 | プロダクト定義（動機・ターゲット・USP・やらないこと） |
| 02 | 技術スタック（採用理由つき） |
| 03 | 情報設計・画面遷移 |
| 04 | 画面別仕様 |
| 05 | データモデル |
| 06 | 機能要件 |
| 07 | 非機能要件 |
| 08 | 課金 |
| 09 | UI/UX のこだわり |
| 10 | ロードマップ |
| 11 | テスト方針 |
| 12 | コーディング規約と AI 向けの指示 |
| 13 | バックログ |
| 14 | 用語集 |
| 15 | Android 版の差分 |

運用のルールを 4 つ決めました。

1. **「未決定」を残さない。** 決まっていないものは仕様の本文ではなくバックログに送ります。本文に「TBD」があると、エージェントがそこを勝手に埋めます
2. **フェーズやバージョンのラベルを本文に書かない。** ロードマップの章だけに置きます。他の章は「完成した製品」を記述します
3. **同じことを 2 か所に書かない。** どちらが正典かを決めて、もう一方からリンクします
4. **変更したら変更履歴のバージョンを上げる。** 追記は必ず末尾で、既存の行の番号は振り直しません（PR や issue から引用されているため）

4 番目は実際に事故りました。変更履歴の表を一括置換で整えようとして、**過去の履歴を壊しかけました**。以後は CI のスクリプトで「末尾追記であること」「番号が飛んでいないこと」を検査しています。

### 「やらないこと」を仕様に書く

これが想像以上に効きました。**実装しないと決めたものを列挙して、理由をつけて仕様に置きます。**

AI に「この機能どう思う?」と聞くと、だいたい「あったほうがいい」と答えます。プッシュ通知、AI による提案、コメント機能、ストリーク。どれも一般には正しいのですが、このアプリの方針とは合いません。

**書いておけば、毎回議論せずに済みます。** 提案が来たら仕様の該当行を引用して終わりです。

## 1 issue = 1 worktree = 1 PR

作業の単位を固定しました。

```
issue を立てる（受け入れ条件をチェックボックスで書く）
  → 専用の worktree を作る
    → 実装させる
      → PR を出す（本文に Close キーワードを入れる）
        → マージで自動クローズ
```

`main` を直接触らせないのが目的です。worktree にしたのは、**複数の作業を並行させても互いのビルド成果物が混ざらない**ためです。

### 踏んだ運用の罠

**Close キーワードをバックティックで囲むと自動クローズされません。** GitHub はコードスパンの中のキーワードを解釈しないので、`` `Closes #123` `` と書くと閉じません。プレーンテキストで書きます。1 回これで手動クローズする羽目になりました。

**PR のマージは issue の完了ではありません。** 受け入れ条件のチェックボックスが残っていることがよくあります。マージ後に必ず issue 本文のチェックを見る、というルールにしました。

**そして本文のサマリ・表・チェックボックス・コメントの 4 か所を揃えます。** どこか 1 か所だけ更新すると、次のセッションで「どれが現状か」が分からなくなります。

### 検証は 1 項目ごとにその場で記録する

これは実際に失敗しました。

実機の検証で 6 項目を通したのに、**issue に残っていたのは後半の 3 項目だけ**でした。前半 3 項目は会話の中で報告して、まとめて書こうとしているうちに文脈が流れて消えました。

翌日の自分は「前半 3 項目は未実施」と判断して、やり直そうとしました。**指摘されて初めて気づきました。**

> **教訓**: 1 項目通したら、その場で issue にコメントとチェックを入れる。
> 「あとでまとめて」は禁止。実測値（件数・時刻・ID）をその場で貼ります。

そして**記録が見つからないときに「未実施」と決めつけない**、というルールも足しました。記録漏れのほうが多いからです。

## 構造ガードを CI に並べる

**仕様に書いても、実装がいつのまにかズレます。** そこで「ズレたら CI が落ちる」形に落としました。30 本以上あります。

| ガードの例 | 何を凍結しているか |
|---|---|
| 表示名の検査 | 全プラットフォームの表示面でアプリ名の表記を統一 |
| ブランド名の翻訳禁止 | 日本語 UI でもブランド名を音写しない |
| 埋め込みブラウザの不在 | 外部リンクが OS のブラウザに渡ることを保証 |
| 新規の日本語文言の検出 | 英訳の付いていない文言を落とす |
| 色とフォントのハードコード検出 | テーマ経由でしか色を使えないようにする |
| `@Environment` の書き方 | 主流から外れた書き方を落とす |
| 変更履歴の形式 | 末尾追記であることと番号の連続性 |

ガードにした基準は「**言葉で言っても守られなかったもの**」です。仕様に書いて、レビューで指摘して、それでも 3 回目が出たものはガードにします。

### ガードは「壊して落ちること」を測らないと 0 円

ここで一番大きな学びがありました。

**「ガードを足した」と「ガードが効く」は別物です。**

あるガードを書いて、11 パターンの負のテスト（意図的に壊して落ちるかを見る）を回したら、**3 本が検知できていませんでした**。

原因は全部同じでした。**判定の前にコメントを剥がしていなかった**ので、自分が書いた説明コメントに検索語が含まれていて、コードを壊してもそこに当たって通ってしまっていました。

```python
# 判定の前に必ずコメントを剥がす
def strip_comments(text):
    text = re.sub(r"/\*.*?\*/", "", text, flags=re.S)
    return "\n".join(line.split("//")[0] for line in text.splitlines())
```

**丁寧にコメントを書くほど当たりやすくなります。** このリポジトリの記述量では構造的に踏みます。

> **教訓**: ガードを書いたら、全パターンで壊して落ちることを実測する。
> 合格が出ても壊れを検知できないなら、そのガードの価値はゼロです。

## 記憶をファイルとして持つ

エージェントとの会話は残らないので、**判断と真因をファイルに書き出しています**。1 ファイル 1 事実の形で、400 ファイルを超えました。

書く対象を 4 種類に決めました。

| 種類 | 内容 |
|---|---|
| ユーザー | 自分がどういう判断軸を持っているか |
| フィードバック | 「こう進めてほしい」という指示。**理由も一緒に書く** |
| プロジェクト | 進行中の作業・制約・決定（コードや履歴から読み取れないもの） |
| 参照 | 外部の手順・URL・環境の前提 |

**コードや git 履歴から読めることは書きません。** 書くのは「なぜそう決めたか」と「何を測って何が分かったか」です。

### 記憶にも誤りが残る

これは今回の記事を書くために読み返して気づいたのですが、**記憶に事実誤認が残っていました**。

「ある API が特定の条件で `MAC_OS` を返す」と書いてあったのですが、実際は**私が UI でその選択をしていた**のが原因でした。API の挙動として一般化して書いていたので、そのまま信じると次のセッションで同じ誤診をします。

訂正を追記して、**「API がそう返す」という説明を立てる前に、UI で何を選んだかを聞く**というルールも書き足しました。

> **教訓**: 記憶は「そのとき正しいと思ったこと」の記録であって、真実の記録ではない。
> 参照するときは、まだ有効かを確かめます。

## エージェントに任せると壊れる場所

4 か月やって、**任せると壊れる場所の型**が見えてきました。

### 1. 「隣の同型」を直さない

同じ不具合が 4 か所にあって、1 か所を指摘すると 1 か所だけ直します。**過去に 3 件直した同型の 4 件目が残っていた**、というのを実際に踏みました。

対処は、修正のときに必ず grep で同型を数えることでした。ただし **grep 自体にも罠があります**。

```sh
# 取りこぼす（先頭のドットを省略した形が漏れる）
grep -rn '= Firestore.firestore()'
# 取りこぼさない
grep -rnE '= *\.(firestore|storage|auth|functions)\(\)'
```

### 2. 「軽微だから後で」と自分で判断する

バグを見つけて「軽微なのでバックログへ」と勝手に決めることがありました。**これは差し戻しました。**

リリース前のプロダクトで「軽微」を先送りの理由にすると、後で別の不具合と組み合わさって見つけにくくなります。**バグなら見つけた時点で直す**、というルールを明示しました。

### 3. 「既存挙動の温存」を理由に部分対応へ逃げる

スキーマやフォーマットを変える必要があるときに、「後方互換のために安全な部分集合で対応します」と提案してくることがありました。

**リリース前なら、正しい方に直すべきです。** 移行の手間はストアを作り直せば済みます。ここは明示的にルールにしました。

> リリースまでは、既存を考えずに根本解決する方を選ぶ。

「触りたくない聖域だから」を理由にしない。部分対応に逃げてよいのは、**根本解決の確証が技術的に取れないとき**だけです。

### 4. 「症状が消えた」を成功と報告する

これが一番危ないものでした。クラッシュが直ったと報告されて、実際は**機能が死んでいた**というケースがあります。どちらでもクラッシュはしません。

対処として、修正後は「効くべき側」も測るというルールにしました。ガードの負のテストと同じ考え方です。

### 5. 断定してから確かめる

ある 1 セッションで、**同じ型の誤りを 7 回踏みました**。仕様書・issue・コードのコメント・ガードの文言に**書き切ってから**確かめています。訂正に PR を 1 本使いました。

7 件とも、測るのに 1〜5 分しかかからないことでした。

ルールを 2 つ決めました。

1. **「A と B が違う → 原因は X だ」と書く前に、X 以外の軸を 1 つ潰す**
2. **仕様や状態を断定する前に、コードか記録を 1 か所読む**

詳しくは[「動いている」を証明する方法](/engineering/app/proving-it-works-device-and-production-verification/)に書きました。

### 6. 出力を鵜呑みにする

これは自分の側の問題です。

**背景実行したコマンドの終了コードを取り違えました。** `make X; echo $?` の形だと、報告される終了コードは末尾の `echo`（常に 0）のものになります。

```sh
# 真の結果をログに明示記録する
{ make X; echo "=====X_EXIT=$?"; } > log
grep X_EXIT log
```

`make ... | tail` も失敗を 0 でマスクします。**iOS のビルドは失敗しても 0 を返すことがある**ので、`BUILD SUCCEEDED` / `BUILD FAILED` の文字列で判定するようにしました。

そしてテストの件数も数えます。「テストを追加した」で満足せず、**`8/8` と出ていることを確認**します。コンパイルエラーで 1 本も走っていなかったことがありました。

## レビューをどう回したか

エージェントに書かせたコードは、別の観点で読み直させています。

リリース前は **18 ラウンドの監査**を回しました。それでも見逃したものがあります。

**Cloud Storage の Rules が全パス deny になっていた**という不具合を、18 ラウンド全部が見逃しました。理由は「テストの失敗を環境要因と判断した」ことでした。しかもその判断が記録に残って、既定事実として扱われていました。

ここから 2 つルールを足しました。

1. **テストの失敗を「環境要因」で片付けない。** そう結論するには、環境を変えて通ることを実証する必要があります
2. **「環境要因」と書かれた記録を見つけたら、鵜呑みにせず再検証する。** 過去の自分の結論も同じ誤りを含みます

詳しくは [Security Rules の記事](/engineering/app/firestore-security-rules-production-only-bugs/)に書きました。

## 分担の落ち着いたところ

4 か月やって、こういう分担になりました。

| 人間側 | エージェント側 |
|---|---|
| 何を作るか / 作らないかの決定 | 実装 |
| 仕様書の記述と維持 | 仕様に沿ったコードとテスト |
| 「これはバグか仕様か」の判定 | 原因の調査と修正案 |
| 実機・本番での実測 | 計測スクリプトの作成 |
| 判断の記録（記憶ファイル） | 記録の下書き |
| ガードにする / しないの判断 | ガードの実装 |

一言でいうと、**判断は人間、作業はエージェント**でした。ただしその「判断」を頭の中に置いておくと消えるので、**仕様書・issue・ガード・記憶という 4 つの形で外に出しておく**必要がありました。

## まとめ

この規模をひとりで回して残ったものを 5 つ。

1. **仕様書を正典にする。** とくに「やらないこと」を書くと、同じ議論が消えます
2. **1 issue = 1 worktree = 1 PR に固定する。** 作業の単位が揺れると、何が終わったか分からなくなります
3. **言葉で守られなかったものはガードにする。** ただし壊して落ちることを実測しないと 0 円です
4. **検証は 1 項目ごとにその場で記録する。** 「あとでまとめて」は消えます
5. **判断と真因をファイルに残す。** 会話は残りませんが、リポジトリに書いたものは残ります

エージェントに任せて速くなったのは実装です。**遅くなったのは判断のほう**でした。書かせるスピードが上がるほど、「これは本当に必要か」「これは本当に直ったのか」を決める時間が相対的に増えます。そこを短縮しようとして断定を急ぐと、この記事に並べた誤りを踏みます。

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

<a class="link-card" href="/engineering/app/ios-android-string-catalog-i18n/">
<span class="link-card__thumb"><img src="/eyecatch/ios-android-string-catalog-i18n.png" alt="" loading="lazy" width="160" height="100"></span>
<span class="link-card__body">
<span class="link-card__title">⑧ iOS と Android を日英対応にした運用（String Catalog と values-en）</span>
<span class="link-card__excerpt">キーが日本語原文なので検出スクリプトはゼロにならない、カタログが 3 つある、文言の一致は機械照合しないと目視では見つからない。</span>
<span class="link-card__meta">shiimanblog.com · 2026.08.24</span>
</span>
</a>

### もうひとつのシリーズ（ストアに出すまで）

- [① 個人開発でストアに登録する前の事務手続き](/engineering/app/indie-app-store-registration-paperwork/)
- [② 個人開発アプリの周辺インフラを用意して踏んだ罠](/engineering/app/indie-app-domain-site-email-oauth-setup/)
- [③ 個人開発アプリを App Store 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/app-store-connect-review-submission-notes/)
- [④ 個人開発アプリを Google Play 審査に出すまでにやったこと・ハマったこと全部](/engineering/app/google-play-console-review-submission-notes/)
- [⑤ 個人開発アプリを App Store と Google Play に同時提出するまでの全工程](/engineering/app/personal-app-cross-store-release-full-journey/)
