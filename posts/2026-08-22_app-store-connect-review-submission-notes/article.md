---
title: 個人開発アプリをApp Store審査に出すまでにやったこと・ハマったこと全部
slug: app-store-connect-review-submission-notes
date: '2026-08-22T10:00:00+09:00'
categories:
  - app
  - ios
tags:
  - app-store-connect
  - testflight
  - xcode
  - swift
  - in-app-purchase
  - store-review
  - duns
draft: false
id: 2246
excerpt: 個人開発のiOSアプリ「Bucket List Commit」をApp Store審査に提出するまでの全工程を記録しました。Individual加入・D-U-N-S・W-8BEN・署名・TestFlight・IAP登録・年齢レーティング・審査ノート・提出まで、実際に踏んだ罠を全部載せています。
eyecatch: ./assets/eyecatch.png
---

## はじめに

個人開発の iOS / Android アプリ **Bucket List Commit** を、App Store と Google Play へ同時に提出しました。この記事はその Apple 側（App Store Connect）で提出までにやったことの記録です。

正直に言うと、アプリのコードを書き終えてからが本番でした。実装は「動くものを作る」だけで済みますが、ストアに出すには契約・税務・署名・課金商品の登録・年齢レーティング・プライバシー申告・審査ノートという、コードとまったく関係ない山を越える必要があります。しかもこの領域には、こういう罠が異常に多いです。

- 公式ドキュメントに載っていない中間ステップがある
- 成功レスポンスが返っているのに実際は何も起きていない
- 画面の表示が「完了」でも実際は未完了

ここで何度も「終わったつもり」になって、後から巻き戻しました。

この記事は、同じ道を通る個人開発者が同じ穴に落ちないことだけを目的に書いています。手順の紹介より「事前に知っていれば消耗しなかった点」に重みを置いています。

> アプリの実装やバグ修正の話は含みません。ストア審査のためだけに追加で実装が必要になった機能については、[全工程まとめの記事](/engineering/app/personal-app-cross-store-release-full-journey/)で扱います。

### この記事の表記について

固有の識別子は伏せています。以下は読み替えてください。

| 記事内の表記 | 実際は |
|---|---|
| `com.example.app` | Bundle ID |
| `1234567890` | Apple の App ID（App Store Connect のアプリ ID） |
| `ABCDE12345` | Team ID |
| `review@example.com` | 審査用アカウント |

アプリ名 `Bucket List Commit` と公式サイト `blc.shiiman.dev` は実名です。

### 前提

- 個人事業主（法人ではない）
- macOS + Xcode。CI は GitHub Actions だが、ストア関連の操作はほぼローカルから
- 課金あり（サブスク 2 種 + 買い切り 6 種）
- 配信地域は**日本 + 米国**の 2 か国のみ
- iOS と Android を**同日に提出**する方針（片方だけ先に出さない）

---

## 提出までの全体像

最初に全体像を出します。**この順序が重要**で、入れ替えると外部待ちが直列につながって数週間損します。

```
Step 0   名義・アカウント所有の確定（誰の名義で出すか）
Step 0.5 屋号の決定
Step 1   開業届に屋号を記載して出し直す → 控えを入手
Step 2   D-U-N-S 番号の取得                    ← 外部待ち（数日〜1.5か月）
Step 3   Apple Developer Program に加入
Step 4   D&B への反映を待って Google Play を組織アカウントで登録
Step 5   Firebase プロジェクト作成 + 従量課金へのアップグレード
Step 6a  サポート用メールアドレス + 事業者サイトの公開
Step 6b  法務ページ（利用規約 / プライバシーポリシー / サポート / アカウント削除）
Step 7   Google OAuth 同意画面の本番公開 + 機密スコープ検証   ← 外部待ち
Step 8   アップロード鍵生成 → release 署名 → AAB → 内部テスト
Step 9-A 【Apple】有料アプリ契約 / 税務情報 / 銀行情報
Step 9-B 【Play】販売・配布契約 / 税務 / 銀行情報
Step 10  課金商品の登録 → 実機・クロス OS 検証 → 提出
```

### リードタイムが律速になるもの

| 項目 | 待ち時間 | 備考 |
|---|---|---|
| D-U-N-S 番号の発行 | 数日〜30営業日 | Apple 経由なら無料・数日で出た |
| D&B グローバルデータへの反映 | 数日〜2週間 | 反映前に Play へ入力すると弾かれる |
| Google Play の本人確認 | 数時間〜数日 | 一度却下されるとやり直し（試行回数に上限あり） |
| Google OAuth の機密スコープ検証 | 数日〜数週間 | Calendar / Sheets を使う場合 |
| Apple Developer Program のアクティベーション | 数時間 | 私の場合は支払いから約11時間 |

**外部待ちの間に何をやるか**を先に決めておくのがコツです。D-U-N-S 待ちの裏で Firebase 構築・法務ページ公開・署名周りを並行させました。

---

## 先に確認すべき「期限のある公式要件」

提出の直前にまとめて調べたのですが、これは着手時に洗い出しておくべきでした。ストアの要件には**期限つきのもの**があり、間に合っていないと提出そのものが弾かれます。

| 要件 | 期限 | 性質 |
|---|---|---|
| Xcode / SDK の最低バージョン | 毎年更新される（私のときは「以降は最新メジャーの Xcode + SDK 必須」） | 提出ブロック |
| 年齢レーティング質問票の新フィールド | 未回答だと新規提出もアップデートもブロックされる | 提出ブロック |
| プライバシーマニフェスト | 必須（required-reason API の申告を含む） | 提出ブロック |
| アクセシビリティの情報表示 | 現時点では任意（「将来必須化」と言うだけで期限は示されていない） | 任意 |

年齢レーティングの質問票は途中で項目が増えます。私が確認したときは新しいフィールドが 10 個ほど追加されていて（健康や医療の話題を扱うか / 暴力表現の程度 / メッセージ機能 / ペアレンタルコントロール / 年齢保証 / ルートボックス / 年齢制限のあるソーシャル機能 / 無制限のウェブアクセス など）、未回答のままだと提出できません。

プライバシーマニフェストは「収集するデータ」だけでなく「使った API の理由コード」も要ります。私の場合は UserDefaults / ファイルのタイムスタンプ / システム起動時間の 3 種で、それぞれ理由コードを申告しました。

> **教訓**: 期限つき要件は「今の自分が満たしているか」を提出前ではなく着手時に調べる。
> 私は提出直前の総ざらいで確認して**たまたま全部満たしていました**が、満たしていなければそこで数日〜数週間止まっていました。

---

## アカウント種別 — Individual だと本名が公開される

ここが最初の分岐で、しかも**やり直しが効きにくい**部分です。

### Apple は Individual か Organization か

|  | Individual（個人） | Organization（組織） |
|---|---|---|
| D-U-N-S | 不要 | 必須 |
| デベロッパー表示名 | 本名固定・変更不可 | 法人名・屋号 |
| 変換 | Individual → Organization は可能 | Organization → Individual は不可 |

Apple の加入フローには、はっきりこう書かれています。

> 法人として登録されていない個人事業主／個人経営者の場合は**正式な個人名で本プログラムに登録されます**。ダウンロードページには**この正式な個人名が表示されます**。

つまり Individual で加入したら、App Store のアプリページに本名が出ます。屋号は使えません。回避手段はありません。

私は「Apple = Individual / Google Play = 組織アカウント」という**非対称な構成**を選びました。理由は次のとおりです。

- Play の組織アカウントには 12人×14日のクローズドテスト要件が免除されるという実利がある
- 一方 Apple 側は Organization にしても D-U-N-S が要るだけで、得られるのは表示名の変更だけ
- Play のデベロッパー名は種別に関係なく自由設定できるので、「組織化すれば本名を隠せる」は Apple 側では成立しない

結果として両 OS でデベロッパー表示名が揃いません（Apple は本名 / Play は屋号）。これは受け入れました。

Individual にはもう 1 つ制約があります。証明書 / Identifier / プロビジョニングプロファイルにアクセスできるのは **Account Holder 本人だけ**です。App Store Connect のユーザーは 50 人まで追加できますが、Developer Program のチームメンバーにはならないので、証明書まわりを他人に任せられません。

 Account Holder の Apple ID は、後から変更するのが極めて困難です。加入する前に、その Apple ID の 2 ファクタ認証と信頼できるデバイスの登録を済ませておくべきでした（この話は[事務手続き編](/engineering/app/indie-app-store-registration-paperwork/)に詳しく書きました）。

> **個人情報の公開範囲を先に確認する**
> Apple の Individual では、**EU 配信時のみ** DSA のトレーダー申告で氏名・住所・電話・メールが App Store 上に公開されます。日本 + 米国配信に絞ればこれは発生しません。詳しくは後述の「EU DSA のバナー」で。

### D-U-N-S が要るのは Organization のときだけ

Organization を選ぶなら D-U-N-S 番号が必須です。これは Apple の加入手続き専用ページから無料・数日で取得できます（書類提出なし）。

- 申請 URL: `https://developer.apple.com/enroll/duns-lookup/`（Apple ID でのサインインが必要）
- 画面に「検索で組織が見つからない場合は、D&Bに情報を送信して、**無料で**D-U-N-S番号の取得を申請することができます」と明記されている
- 入力は半角英数字のみ。日本語も中黒（`・`）も通らない

D-U-N-S は事業体に 1 つ割り当てられる識別番号で、用途に紐づきません。Apple 経由で取得しても Google Play の組織アカウント登録に使えます。当初これを知らず、日本の窓口の有料コース（16,500円）を勧めてしまいました。

**教訓: 有料サービスを検討する前に「専用ページ・無料ルートの有無」を確認する。**

> 屋号の決定 / 屋号入り開業届の取り直し（税務署）/ D-U-N-S の 3 ルート比較 / 名義の可逆性 / 2FA と復旧手段は両ストア共通の前段なので、[事務手続き編](/engineering/app/indie-app-store-registration-paperwork/)にまとめました。Apple の各画面で入力するローマ字表記（屋号・住所）を 4 か所で完全一致させる話もそちらです。

### Apple Developer Program の加入

- 入口は `https://developer.apple.com/programs/enroll/` の「Start your enrollment」。Web で完結する
- 「個人情報を確認してください」画面は **日本語入力不可・半角英数字のみ**。D-U-N-S と同じローマ字表記を使う
- **「続ける」を押した後は氏名の変更が不可**。政府発行の身分証と一致していること
- 料金は **¥12,980 / 年**
- アクティベーション完了の判定は「メール 2 通のペア」

```
- Apple Developer Programへようこそ
- Welcome to App Store Connect.     ← 有効な有料メンバーシップがないと送られない
```

後者が来ていれば確実に有効化されています。前者だけだと判断できません。

登録 ID / 注文番号と Team ID は別物です。登録 ID（`9YQV…` 形式）と注文番号（`W15…` 形式）は Apple への問い合わせ用の参照番号で、Team ID は承認後に Membership 詳細で確認します。最初これを混同しました。

### identifier はグローバルに一意 — App Group で衝突した

App ID を登録する前に、App Group と iCloud Container の identifier を先に作る必要があります（App ID の Capabilities 画面では既存 identifier を選ぶだけで、インラインで新規作成できない）。

ここで踏んだのがこれです。

```
group.com.example.app  →  not available
```

App Group ID は全 Apple Developer アカウントを通じてグローバルに一意です。自分のアカウントの App Groups 一覧は空なのに拒否される = **外部の誰かが押さえている**。Apple のグローバル名前空間には照会手段がなく、理由は特定できません。

`group.com.example.app.shared` に変えたら通りました。iCloud Container は別名前空間なので当初の名前で通っています。

> **教訓**
> リネームの可否判断を「一覧に無い」で終わらせず、別文字列を 1 回試して名前空間チェックが正常に動いていることを確認する。同じ操作が診断にも修正にもなります。未リリースなら共有コンテナのデータ移行が不要なので、リネームは安いうちに済ませます。

---

## 有料アプリ契約・税務・銀行情報

ここが有効になるまで IAP は 1 つも作れず、Sandbox 課金テストもできません。App Store Connect の API で `inAppPurchasesV2` が 0 件・`subscriptionGroups` が 0 件だったのを「API の使い方が悪いのか」と疑って時間を使いましたが、原因はこれでした。

有効化は規約のダイアログを読んで同意するだけです。**添付ファイル（Exhibits）も規約の一部**なので、ここでダウンロードして保存しておくといいと思います。

![App Store Connect の「有料アプリ契約」ダイアログ。利用規約の本文と「添付ファイル」の折りたたみが表示され、下部に「上記の利用規約を読み、それに同意します」のチェックボックスがある](./assets/asc-paid-apps-agreement.png)

### 場所が分かりにくい

App Store Connect → 「ビジネス」です（旧称「契約/税金/口座情報」）。直リンクは `https://appstoreconnect.apple.com/business`。**メニュー名が改名されていて見つけにくい**うえ、アカウント所有者でサインインする必要があります。

未整備の状態はこう見えます。契約 / 銀行口座 / 納税フォームの 3 ブロックが縦に並び、上部にバナーが出ます。

![App Store Connect の「ビジネス」ページ。上部に EU のトレーダーステータスを求める赤いバナーと、銀行口座・納税フォームの追加を促す青いバナーが並び、その下に契約・銀行口座・納税フォームの 3 ブロックが縦に並んでいる](./assets/asc-business-page.png)

一番上の赤いバナーが EU の**トレーダーステータス申告**です。「有料アプリ契約」のステータスが「ユーザ情報を保留中」のままだと、この後の課金商品を 1 つも作れません。

揃うとこうなります。

```
- 無料アプリ契約                                          有効
- 有料アプリ契約                                          有効
- U.S. Certificate of Foreign Status of Beneficial Owner  有効
- U.S. Form W-8BEN                                        有効
- 銀行口座                                                有効・使用中
```

### 納税フォームは 2 つある

1. U.S. Certificate of Foreign Status of Beneficial Owner（Apple 独自）
   - `Type of Beneficial Owner` は `Individual/Sole proprietor` が自動で入る
   - **Title** に入力が必要 → **`Owner`**（個人事業主が本人署名する場合の標準表記）
2. **U.S. Form W-8BEN**（IRS 様式）

### W-8BEN の Line 10 は空欄が正解

これは私が誤った案内をして、Apple 自身の資料に覆された箇所です。

最初は「Line 10 に租税条約 12条1項 / 0% / Income from the sale of applications を書く」と案内しました。しかし Apple が配布している `Download Form W-8BEN Tips Sheet` にこう書かれています。

> **Line 10**: Complete this item only if you are eligible to claim any applicable treaty benefits that require you to meet conditions not covered by the representation on line 9. It is expected that Line 10 would not normally be applicable.
> **Line 5 and 6**: A U.S. TIN or Foreign TIN is required in order to receive any applicable benefit of the reduced tax treaty rate.

実際の入力画面がこれです。Line 9 のチェックだけ入れて、Line 10 の 2 つの入力欄は `Optional` のまま空にしておくのが正解でした。

![W-8BEN の Part II: Claim of Tax Treaty Benefits。Line 9 の居住国のチェックだけが入り、Line 10 の Article と paragraph、税率の入力欄はどちらも Optional のまま空になっている](./assets/asc-w8ben-treaty.png)

Line 9（居住者証明）+ Line 6a（外国 TIN）が揃えば、Apple が条約表から標準税率を適用します。日本の使用料は日米租税条約 第12条第1項で標準 0%。**空欄が 0% 適用の正規ルート**です。

さらに罠があります。

ラジオボタン（`Income from the sale of applications`）は一度選ぶと解除できません。誤って選んだら**ページをリロードして初期状態に戻す**。半端に埋まった Line 10 は内部矛盾した申告になります。

確定した入力値はこうなりました。

| 欄 | 値 |
|---|---|
| 5. U.S. TIN | 空欄（EIN / SSN のラジオも選ばない） |
| 6.a. Foreign TIN | マイナンバー 12 桁（必須） |
| 7. Reference Number(s) | 空欄（Individual なので屋号は不要） |
| 8. Date of Birth | `MM-DD-YYYY` 形式 |
| ○9. | チェック |
| 10. | 完全に空欄・ラジオも未選択 |
| ○Part III + capacity to sign | チェック |

> Tips Sheet 自体は 2015 年版（「iTunes Connect」表記のまま）ですが、Line 10 の扱いは IRS の指示そのままなので現在も有効です。

### 銀行口座の罠 2 つ

| 欄 | 罠 |
|---|---|
| 国または地域 | 既定が「アメリカ合衆国」。ただし「法人と同じ」をチェックすると住所ごと自動で日本に直る |
| 口座名義人名 | 「法人と同じ」では入らない。銀行の登録どおりに手入力が必要 |

口座名義人名は「銀行口座に表示される名前と句読点も含め完全に一致」が要求されます。日本円口座は英字・カタカナ・数字 + `-)(` が許可されており、半角カタカナ（`ﾔﾏﾀﾞ ﾀﾛｳ` のような形式）で受理されました。半角の濁点が弾かれるのではという懸念は不要でした。

- 口座名義人の種類 = **個人**
- **口座の種類 = 普通**（Apple は「普通」が既定。Play は既定が「貯蓄」になっているので注意 —— [Google Play Console 編](/engineering/app/google-play-console-review-submission-notes/)参照）

### ロイヤルティ通貨 USD は「集計単位」で、着金は JPY

「すべての銀行口座」の行に `銀行通貨 JPY / ロイヤルティ通貨 USD` と出て、JPY のロイヤルティ通貨行が存在せず、`+` を押しても通貨の選択項目が出ません。

私は「日本の売上が USD → JPY の二重換算を通るのでは」と懸念しましたが、公式ドキュメントで否定されました。

> Apple's bank consolidates proceeds for each currency in which you have App Store sales, resulting in a single payment to your bank per fiscal month … the total proceeds owed to you displayed in the currency of your bank account

| 項目 | 挙動 |
|---|---|
| `ロイヤルティ通貨 USD` | 集計・レポート上の通貨。振り込まれる通貨ではない |
| 着金 | 登録口座の通貨（JPY）で月 1 回にまとめて入金 |
| 期限 | 会計月末から 45 日以内 |

通貨の選択肢が出ないのも整合しています（公式に「複数口座への分割払いはサポートされない・支払いは primary 口座に行われる」とあり、口座が 1 つなら Apple が自動割当する）。設定を触る必要はありません。

### EU DSA のトレーダー申告バナーは触らない

「ビジネス」ページ上部に赤いバナー（トレーダー申告）が出ます。これに答えると 氏名・住所・電話・メールが EU の App Store 上で公開されます。Individual アカウントなので、**自宅住所が公開される**ことになります。

配信地域から EU を外せば申告自体が不要です。

```
アプリ → <アプリ名> → サイドバー「価格および配信状況」
  → 「アプリの配信状況」→ 「配信状況の設定」
      → ラジオ 3 択で 「特定の国または地域」を選ぶ
```

`availableInNewTerritories` を必ず `false` にします。これを放置すると、今後 Apple が追加した地域に自動で入り、EU が混ざります。

「ビジネス」ページに表示される「175個の国または地域」は契約の適用地域で、アプリの配信設定とは別物です。ここは混同しやすいので気をつけてください。

 API の落とし穴。

```bash
# 未設定のときは NOT_FOUND が返る
GET /v1/apps/1234567890/appAvailabilityV2
→ 404 "There is no resource of type 'appAvailabilities'"
```

これを「API 非対応」と誤読しました。設定後は普通に読めます。

```bash
GET /v1/apps/1234567890/appAvailabilityV2
GET /v2/appAvailabilities/1234567890/territoryAvailabilities?limit=200
```

実測で「総 175 件 / `available: true` は USA + JPN の 2 件のみ / `availableInNewTerritories: false`」を確認しました。

こういう全数確認は、スクリーンショットではなく API でやります。フルページのスクショは 2950×15503 → 381×2000 に縮小され、文字が読めなくなります。

---

## アカウント削除は「消えないもの」を先に決める

Guideline 5.1.1(v) の**アプリ内でのアカウント削除**は避けられない要件です。実装で決めた設計を残しておきます。

**削除の順序**: サーバーのデータ → 端末の資格情報 → 認証ユーザー

順序に意味があります。個人情報を先に消したいので、まずユーザー本体のドキュメントから消します。さらに **途中で失敗したら認証ユーザーは消しません**（消してしまうと、残ったサーバーデータに本人が二度とアクセスできず、削除も不可能になります）。

しかも「完全に消えない部分」があることを、仕様として先に明記しました。

| 残るもの | 理由 |
|---|---|
| 相手側から自分へ向いている関係のデータ | アクセス制御上、自分だけでは消せない（相手側の掃除に委ねる） |
| 他人の受信箱に入った、自分が送った通知 | 受信者だけが削除できる |
| 外部サービス側のデータ（連携先のリポジトリなど） | 本人の所有物なので削除せず、連携解除だけ行う |

これを公開ページの「アカウント削除について」にも書きました。「削除します」とだけ書いて実際は残る部分があると、ストアの申告とも実装とも食い違います。「いいねの参照は匿名で残る」「削除してもサブスクリプションは自動解約されない」まで明記しました。

「アプリのデータを削除」と「アカウント削除」は別要件です。前者は端末のローカルデータのリセット、後者はサーバー側の身元の削除で、**両方 UI に必要**でした（それぞれ二重確認つき）。

## アプリレコードとアプリ名の確保

### API ではアプリを作れない

```
POST /v1/apps
→ 403 FORBIDDEN_ERROR
   The resource 'apps' does not allow 'CREATE'.
   Allowed operations are: GET_COLLECTION, GET_INSTANCE, UPDATE
```

アプリレコードの作成 = アプリ名の確保は Web UI 専用です。一度「API で作成できる」と誤って提案しました。

API でできるのは 存在確認（GET）/ メタデータ更新（PATCH）/ ビルドのアップロード / TestFlight 管理です。ここを押さえておくと、以降「これは API か画面か」の判断が速くなります。

### App Store 名は先着制

App Store のアプリ名は先着制で、App Store Connect でアプリレコードを作った時点に予約されます。アカウントが有効化されたら最優先で押さえるタスクです（Play は名称の一意性を要求しないので、この論点は Apple だけ）。

私は #362（リリース要件のトラッキング Issue）に「アプリ名の確保」項目そのものが無く、指摘されて追記しました。しかもその後、**5 日間チェックが倒れていない状態で放置**していました。実際にはコメントが「次にやること」として挙げた **8 分後に作成が完了していた**のです。

> **教訓**: API / 実サーバで検証できる項目は、Issue の記載よりも実測を信じる。

### 商標スクリーニングの機械的な経路

アプリ名を決める前に、最低限のスクリーニングはやっておきました。調査手段の当たり外れを共有します。

| 調査先 | 結果 |
|---|---|
| App Store 既存アプリ | iTunes Search API で日本・米国とも完全一致を確認できる |
| 米国 USPTO | Justia / uspto.report は 403。`https://tsdr.uspto.gov/statusview/sn{serial}` は curl で 200 が取れる（本家・唯一の機械的経路） |
| 日本 J-PlatPat | SPA で API 照会不可 = 手作業必須。商標検索の `称呼（単純文字列検索）` が最も網が広い（英字商標にも称呼が付く）。3欄は AND なので 1 欄だけ埋める |

第9類（ソフトウェア）に同名の登録があっても、指定商品が重ならなければ効力は及びません。私が見つけた同名登録は指定商品がアイウェア専門（類似群 10B01 / 23B01）で、ソフトウェアの 11C01 を含んでいませんでした。「区分だけ見て諦める」のは早計です。

### CFBundleName は識別子ではなく表示面

これは見落としやすい罠です。

`CFBundleName` を未設定にすると `$(PRODUCT_NAME)` にフォールバックします。私のプロジェクトでは `PRODUCT_NAME` がスペースなしの内部名だったので、iOS のサインインダイアログ（`ASWebAuthenticationSession`）にその内部名が表示されました。

```
「BucketListCommit」がサインインのために…
```

`CFBundleName` は `PRODUCT_NAME` とは独立に、表示名として明示設定する必要があります。同様の表示面は他にもあります。

- plist の表示系キー（`CFBundleDisplayName` / `CFBundleName`）
- String Catalog（`InfoPlist.xcstrings`）のローカライズ
- Android の `AndroidManifest` の `android:label` / strings

私はこれらすべてを CI のガードスクリプトで凍結しました（表示面が 1 つでもズレたら落ちる）。

> おまけ: iOS のホーム画面のアイコンラベルは短縮される
> SpringBoard は**スペースを削除してから**、1 行に収まらない場合に省略します。シミュレータで実際に確認すると `Bucket List Commit` → `BucketListCom…` になります。`simctl listapps` / `devicectl device info apps` はフルネームを返し、設定 / Siri / サインインダイアログ / App Store もフル表示なので、**設定は正しくてもホーム画面だけ短くなる**のは仕様です。回避手段はありません。

---

## 署名と配布証明書 — CLI だけで完結する

Portal をブラウザで触る必要はありませんでした。自動署名前提（`project.yml` に `CODE_SIGN_STYLE` を書かない）で、次の 2 コマンドです。

```sh
xcodebuild archive -project MyApp.xcodeproj -scheme MyApp \
  -configuration Release -destination 'generic/platform=iOS' \
  -archivePath <dir>/MyApp.xcarchive -derivedDataPath /tmp/archive \
  -allowProvisioningUpdates

xcodebuild -exportArchive -archivePath <dir>/MyApp.xcarchive \
  -exportPath <dir>/export -exportOptionsPlist <dir>/ExportOptions.plist \
  -allowProvisioningUpdates
```

`ExportOptions.plist` の要点はこうです。

| キー | 値 | 意味 |
|---|---|---|
| `method` | `app-store-connect` |  |
| `teamID` | `ABCDE12345` |  |
| `signingStyle` | `automatic` |  |
| `destination` | `export` | App Store へ送らず発行だけで止める |
| `uploadSymbols` | `true` |  |
| `manageAppVersionAndBuildNumber` | `false` | 番号を Xcode に触らせない |

### 配布証明書が作られるのは archive ではなく export の時点

これは 2 回踏みました。

```sh
# archive 成功直後にこれをやると…
codesign -dv <app>
→ Authority=Apple Development:     ← 「発行されていない」と誤読する
```

さらに `security find-identity -v -p codesigning` にも配布 identity は現れません。これも「失敗した」と誤読しやすいところです。本当の確認は ipa 側を見ます。

```sh
unzip -q "<export>/MyApp.ipa" -d <dir>/ipax
codesign -dv --verbose=4 <dir>/ipax/Payload/*.app   # Authority=Apple Distribution: を確認
security cms -D -i <app>/embedded.mobileprovision   # ProvisionedDevices 0 件 = 配布用
find <app> -name "*.appex" -maxdepth 3              # 拡張も同様に確認
```

`.appex` は `PlugIns/` にあります（`Extensions/` は存在しないので、zsh の glob だと no-match で行ごと落ちる → `find` を使う）。

### Cloud Managed 証明書なので .p12 バックアップは不要

発行される証明書は `Cloud Managed Apple Distribution` で、**ローカルに秘密鍵がありません**（Apple が保持）。

実際に確認: `security find-identity -v` に配布 identity なし / `security find-certificate -c "Apple Distribution"` も 0 件 / `/v1/certificates` にも現れない。

要するに、こういう性質です。

- 「紛失」の概念がない。マシンを変えても Apple ID のサインインで再取得できる
- 逆に **`.p12` を書き出すことはできない**

Android のアップロード鍵（紛失すると Play サポートにリセット依頼・パスワードマネージャ保管が必須）とは性質がまったく違います。最初「.p12 をバックアップしましょう」と提案しかけました。同じ扱いを勧めると間違いです。

### `No provider associated with App Store Connect user` は無害なログ

export のログには毎回これが出ます。エラーではありません。

```
IDEDistribution: App Store Connect request for store configuration failed for account (null)
  ... "No provider associated with App Store Connect user"
** EXPORT SUCCEEDED **       ← ちゃんと出る
```

私はこれを致命的エラーと誤読し、3 つの別事象を 1 つの原因に括りました。実際に起きていたのは以下の 3 件です。

1. 本当の失敗は `error: exportArchive The request expected results but none were found`
   → **もう一度そのまま叩いたら成功した**（= 一時的な App Store Connect 側の不調。同時刻に GitHub Actions も `Service Unavailable` を返していた）。**まず素直にリトライする**
2. `Cloud signing permission error` / `No signing certificate "iOS Distribution" found` は
   自分が足した `-authenticationKeyPath` が原因。`altool`（validate / upload）は API キーで動くが、`xcodebuild -exportArchive` の Cloud signing に API キーを渡すと逆に壊れる。素の `-allowProvisioningUpdates` だけにする
3. `defaults read com.apple.dt.Xcode | grep Account` の出力を見て「Xcode のアカウントが空」と判断したが、**実際は 2 つサインイン済み**だった（grep が拾った行が展開されていないだけ）

> GUI の状態を CLI の grep 断片で判定しない。これで不要な再認証をユーザーに頼んでしまいました。

なお、後日別の理由で本当に再サインインが必要になったケースもあります。export が `Unable to log in with account` で失敗したときは、Xcode の有料チーム認証キャッシュ切れが原因で、リトライでは直らず GUI での再サインインが必要でした。再サインイン後も plist には有料チームが現れないのに export は成功します。やはり CLI の断片では判定できません。

### export が落ちたときの手順

archive は成功しているので export からやり直せます（フルビルド不要・数分）。**ビルド番号も消費されません**（消費するのは `altool --upload-app` のみ）。

```sh
# 1. 素のオプションでリトライ（API キーは渡さない）
xcodebuild -exportArchive -archivePath <dir>/MyApp.xcarchive \
  -exportPath <dir>/export -exportOptionsPlist <dir>/ExportOptions.plist \
  -allowProvisioningUpdates
# 2. ipa の署名を実測（Authority=Apple Distribution / get-task-allow=false）
# 3. validate   ← 番号を消費しない
# 4. upload     ← ここで初めて消費
```

---

## Release ビルドの検証は「設定を読む」では足りない

提出物が本当に Release 相当かは、ビルド設定を読むのではなく、archive から出した ipa の実バイナリに対して確認しました。設定が正しいことは「実際に入っていない」ことの証明にならないからです。

| 検証 | 結果 |
|---|---|
| `#if DEBUG` コードの混入 | ○構造的に不可能。`xcodebuild -showBuildSettings` で確認すると Release は `SWIFT_ACTIVE_COMPILATION_CONDITIONS` が未設定（Debug だけ `DEBUG`）= コンパイルすらされない |
| デバッグ用ログ | ○`print(` / `NSLog` が 0 件（`#if` のネストを追跡して全ソースを走査） |
| Debug 画面の文言残存 | ○0 件（Debug 画面のソースから抽出した 8 文言を Release バイナリ内で検索） |
| 最適化 | ○`-O` / `wholemodule` / `ENABLE_TESTABILITY = NO` / dSYM 生成 |

**ここが一番大事**でした。

「検査自体が機能していること」も確認しました。本番の文言（キャッチコピーやアプリ名）が**ちゃんと検出される**ことを先に確かめています。これが無いと「0 件」と「検査が壊れている」を区別できません。

偽陽性も 1 件踏みました。Debug 専用だと思っていた文言が 1 件ヒットしたのですが、実体は**本番データ側の文言**で、Debug のサンプルデータがそれを流用していたためでした。この種の検査は、ヒットした文言の出自を確認しないと誤検出します。

iOS と Android で保証の強さが違う点も押さえておくべきでした。

|  | 仕組み | 強さ |
|---|---|---|
| iOS | `#if DEBUG` のコンパイル時除去 | 構造的保証（そもそもバイナリに存在しない） |
| Android | R8 の到達不能コード削除 | 最適化への依存（到達可能と判断されれば残る） |

### Release 実機だけで出る不具合は「ログが読めるか」で調査時間が変わる

Release の実機でしか再現しない不具合が複数出ました（起動直後のクラッシュ / 外部連携の権限不足 / プッシュが一度も届かない）。そこで詰まったのがログです。

ログの出力を全部「非公開」扱いにしていたので、実機ログでカテゴリと重大度しか読めませんでした。本文が読めないと、実機だけで再現する不具合の切り分けが著しく遅くなります。

対処はこの 2 つです。

- **Debug ではログ本文を出す**ようにした（`#if DEBUG` で切り替え）
- ただしそれでは Release 実機だけで出るものには無力なので、「過去に実機だけで再現した不具合が通った経路」に限定して診断ログを撒いた（撒きすぎると読めなくなる）
- 診断ログは消してもビルドも lint もテストも通る（診断能力が黙って落ちる）ので、**CI のガードで配線を凍結**した

**実機ログの回収には root 権限が必要**でした（`log collect`）。リリース前に「実機のログを本文まで読める経路」を作っておくと、後半の調査が明らかに速くなります。

輸出コンプライアンスは自動回答にできます。`ITSAppUsesNonExemptEncryption` を `false` に設定しておくと、提出時に聞かれず API 側も `usesNonExemptEncryption: False` になります（実際に確認しました）。毎回手で答える項目を 1 つ減らせます。

## アップロードは「転送 → 受理 → 配信」の 3 段

ここは個人開発で最もハマると思います。3 段それぞれで別の罠を踏みました。

|  | 転送 | 受理 | 配信 |
|---|---|---|---|
| 判定 | `UPLOAD SUCCEEDED` | `processingState=VALID` | ベータグループの `builds` に載る |
| 何が起きる | ipa が Apple に届く | Apple が中身を検査して受理 | TestFlight のテスターに見える |

### 第0段: そもそも「アップロードできる状態」ではなかった

ビルド番号を整えて初めてアップロードを試したとき、アプリは一度もアップロードできる状態になったことがありませんでした。

```
altool --validate-app
→ 91179 Invalid extension
```

原因は、App Intents の拡張ターゲットが ExtensionKit の要件（`Extensions/` 配下 + 専用の Info.plist キー）を満たしていなかったことです。旧来の `app-extension` + `NSExtension` の形で定義していました。

`make build` も `xcodebuild archive` もすべて成功します。これはサーバ側のバリデーションでしか検出できません。

調べてみると、その拡張にはエントリポイントが無く、実行時に何も担っていませんでした（実機での検証も一度もされていなかった）。App Intents は本体ターゲットに入っていればショートカットや Spotlight から動くので、拡張を削除して本体一本にしました。

`--validate-app` はアップロードしないのでビルド番号を消費しません。だからこれをローカルのゲートにできます。**ゲートが空振りでないことも確認しました**（修正前の ipa は `91179` で落ち、修正後は通る）。

### 第1段: `UPLOAD SUCCEEDED` は「受理された」ではない

`altool` が `UPLOAD SUCCEEDED` を返しても、Apple が受理していないことがあります。

私のビルドは `ITMS-90626` で弾かれました。文面はこうです。

```
ITMS-90626: Invalid Siri Support - App Intent description
'... Siri から達成扱いにします。' cannot contain 'siri'
```

App Intent の説明文に "siri" という語を含められないという制約でした（ja / en の両方に入っていた）。

さらに重要なのが、**この検査は前段をすり抜けること**です。

| 段階 | 結果 |
|---|---|
| `make build` / `archive` | ○通る |
| `--validate-app`（ローカルゲート） | 通ってしまう（この検査を含まない） |
| `--upload-app` 後の Apple 側処理（ITMS-*） | ×ここで落ちる |

API のビルド一覧がいつまでも 0 件だった真因がこれで、Apple の通知メールを読んで初めて分かりました。`UPLOAD SUCCEEDED` の後は、API にビルドが出現して `processingState` が確定するまで追う必要があります。

弾かれてもビルド番号は消費されます。1 本を却下で失いました。

却下はメールで通知されます。転送成功のログだけ見ていると永久に気づけません。しかも再発防止として、説明文に禁止語が入らないことを CI で凍結しました（ローカライズ値も検査対象にする）。

### 第2段の罠: `/v1/apps/{id}/builds` は順不同

受理の確認は `processingState` を見ますが、ビルド一覧の取り方に致命的な罠があります。

```
# 総13件のときの limit=10 の実測
['3332', '3369', '3279', '3268', '3294', '3270', '3231', '3257', '3264', '3249']
→ 最新の 3384 が入っていない。そもそも「最新10件」ですらない（昇順でも降順でもない）
```

`sort=-version` / `sort=-uploadedDate` は HTTP 400（`PARAMETER_ERROR`）なので、「新しい順に取る」ことができません。

正しい引き方はこれです。

```
/v1/builds?filter[app]={appId}&filter[version]={番号}&fields[builds]=version,processingState
```

`/v1/apps/{id}/builds` ではなく `/v1/builds` + `filter[app]` です。全件が必要なときだけ `limit=200` を使い、**探索には使わない**。

これが起こした実害はこうでした。受理判定スクリプトが `limit=10` で書かれていたため、ビルド総数が 10 件を超えた時点で新しいビルドが永久に「存在しない」判定になりました。

```
UPLOAD SUCCEEDED → 「API に未出現」×30分 → タイムアウト → make が失敗
→ その後の TestFlight グループ割当が実行されない
```

症状は「TestFlight には上がっているのにテスターに配信されない」です。しかも**ビルド番号は消費済み**なので同じ番号でやり直せません。

さらに私自身、`limit=8` の一覧を見て「このビルドは App Store Connect に存在しない」と報告しました。実際は `VALID` で存在していました。

> `limit` を小さくした一覧の「無い」は、「無い」ことの証拠になりません。
> 存在確認は必ず `filter[version]`、または `limit=200` で全件取ってから判定する。

スクリプトのログ末尾に出る成功表示も鵜呑みにできません。あるビルドではログが「アップロードの転送に失敗しました」で終わっていたのに、App Store Connect 上は `VALID` + 配信済みでした。結論は必ず実データで出しましょう。

### 第3段: 受理されてもテスターには配信されない

`processingState=VALID` まで通っても、TestFlight アプリにビルドは出てきません。ベータグループへの割当が別 API で、そこが自動化されていませんでした。

グループの画面は「テスター / ビルド / 設定」の 3 タブに分かれていて、**ビルドの割当はこの「ビルド」タブの中**です。アップロードした一覧とは別管理になっています。

![App Store Connect の TestFlight 内部テストグループの画面。テスター / ビルド / 設定 の 3 タブが並び、設定タブにグループ名や「Apple シリコン搭載 Mac で iPhone および iPad アプリをテストする」「ビルドの配信」の項目が表示されている](./assets/asc-testflight-internal.png)

```
build 3268  VALID    ← Apple は受理している
グループ「内部テスト」の割当: 3199 / 3205 / 3231 / 3249 / 3257 / 3264   ← 3268 が無い
→ TestFlight には 3264 までしか出ない
```

判定と修復はこうです。

```bash
# グループ一覧（hasAccessToAllBuilds=false なら手動割当が必要）
GET /v1/apps/{appId}/betaGroups?limit=10
# 割当済みビルド
GET /v1/betaGroups/<groupId>/builds?limit=20
# 割当（HTTP 204 が成功・空ボディ）
POST /v1/betaGroups/<groupId>/relationships/builds
  {"data":[{"type":"builds","id":"<buildId>"}]}
```

### 自動配信は後から ON にできない

- **API**: `PATCH /v1/betaGroups/<id>` に `hasAccessToAllBuilds` を含めると **HTTP 409**
  → `The attribute 'hasAccessToAllBuilds' can not be included in a 'UPDATE' operation`（**作成時専用**）
- **UI**: TestFlight → グループ → 設定 に「ビルドの配信 / Xcodeビルドを手動配信」と**テキスト表示があるだけでトグルが無い**（同じ画面の他の行は編集リンクを持つので、無いのは仕様）

残る手は 2 つだけです。

1. アップロード処理の末尾に割当工程を組み込む
2. `hasAccessToAllBuilds: true` でグループを**新規作成**してテスターを移す（旧グループのビルド履歴・フィードバックの紐づきが切れる）

私は 1 を採りました。HTTP 204 を信用せず、割当後に `/v1/betaGroups/<id>/builds` を再取得して確認する実装にしています。

なお TestFlight まわりは、グループの作成・ビルドの割当・テスターの追加まで全部 API でできます（`POST /v1/betaGroups` は `isInternalGroup: true` を受け付ける）。API 不可なのはアプリレコードの作成だけでした。

実機に入ったことの確認は、テスターの `state` が `INVITED` → `INSTALLED` に遷移することで取れます（内部テストは審査不要なので、ここまで数分です）。

> TestFlight で配布を始める前にやるべきだったこと
> ローカル DB のスキーマ移行の仕組みを先に用意しておくべきでした。配布を始めると**実データを持った端末が生まれます**。その状態で初回のスキーマ変更をすると、**そのデータは救えません**。
> 私は「配布前ならスキーマ変更が無料（移行不要）」という前提で進めていましたが、この前提が切れる瞬間が TestFlight の配布開始です。CI でスキーマ変更を検知して「配布前だから凍結ファイルの更新だけ」か「配布後だから移行計画が必要」かを判断できるようにしました。

### その他、この工程で踏んだ罠

- `archive` タスクの中でバリデーションまで走る構成にしていると、`make ios-archive && make ios-validate` のように繋いだ 2 回目の make が新しい出力ディレクトリ（タイムスタンプ付き）を計算して ipa を見失う。別コマンドに跨って使うなら出力先を明示する
- バックグラウンド実行で `| tail` すると終了コードが tail のものになり、失敗が exit 0 で通知される。`make verify 2>&1 | tail -80` で実際に false green を踏みました。**パイプしない**
- シェル変数を全角括弧で囲むと `set -u` で即死する（`「$gname」` の `「` の先頭バイトが変数名に食われる）。`${gname}` で閉じる

---

## バージョン番号の設計

両 OS で番号を共有する構成にしました。

|  | 値 |
|---|---|
| marketing version | `1.0.0`（3桁・両 OS 共通） |
| build 番号 | `git rev-list --count HEAD`（コミット数）を両 OS 共有 |

この方式に至った理由です。

- build 番号は使い捨て。惜しむものではない。ユーザーが見るのは marketing version だけで、テストで 30 本焼いてリリースが build 31 になっても普通です。一度アップロードした番号はビルドを削除しても再利用できないので、節約する意味がありません
- Android の `versionCode` は `versionName` と無関係にグローバル単調増加が必須でリセット不可（iOS はバージョンを上げれば build を 1 に戻せる）。両 OS 共通にするには「一生リセットしない」方式が唯一の解
- 日付時刻方式（`202608221230`）は 12 桁で int32 上限 2,147,483,647 を溢れるので Android では使えない
- 同じコミットから 2 回アップロードすると番号が衝突する → 「アップロード前にコミットする」運用 + 環境変数での上書き口を用意

### 拡張のバージョンはホストアプリと一致必須

XcodeGen で生成していると、拡張ターゲット（Widget / AppIntents）に `CFBundleVersion` の指定がないと既定値 `1.0` / `1` が入ります。生成後の `Widgets/Info.plist` を実際に確認して気づきました。

本体だけ動的化すると不一致で App Store のバリデーションに落ちます。3 ターゲットすべてを `$(MARKETING_VERSION)` / `$(CURRENT_PROJECT_VERSION)` 参照にしました。

### `xcodebuild` を直接叩くと番号が既定値のまま

`project.yml` の既定値のままビルドされ、2 回目のアップロードが「番号の重複」で弾かれます。番号を注入する make ターゲット経由に統一しました。

---

## IAP（アプリ内課金）8種の登録

サブスク 2 種（月額・年額）+ 買い切り 6 種を App Store Connect API で登録しました。ここは罠の密度が高い区間です。

### 登録内容と不可逆な選択

|  | 商品 | 価格 |
|---|---|---|
| 非消費型 6 | `theme.*` | ¥250 / $1.99 |
| サブスク 2 | `pro.yearly`（groupLevel 1）/ `pro.monthly`（groupLevel 2） | ¥2,800・¥280 / $17.99・$1.99 |

- ファミリー共有は全商品オフにした。有効化は後からできるが**無効には戻せない**
- 年額を groupLevel 1 にすると、月額 → 年額が即時アップグレード（日割り精算）、逆は次回更新時になる
- 商品 ID は恒久的で、削除しても再利用できない
- 商品 ID は StoreKit 設定ファイル / iOS 実装 / Android 実装の**3 者一致を確認**した

### availability を設定しないと価格 POST が 409 で落ちる

これが一番時間を溶かしました。

```
POST /v1/subscriptionPrices
→ 409
   pointer: /data/relationships/subscriptionPricePoint/id
   「価格情報の処理中にエラー」
```

エラーが価格ポイントを指しているので「価格ポイント ID が間違っている」と誤読します。別の価格ポイントを何個試しても同じエラーです。

正解は先に availability を作ることでした。

```
POST /v1/subscriptionAvailabilities      # 価格より先
POST /v1/inAppPurchaseAvailabilities     # 価格より先
```

### サブスクは全部埋めても `state` が MISSING_METADATA のまま

これも実害のあった誤読です。

|  | API の `state` | ASC の UI |
|---|---|---|
| 非消費型 | 埋めれば `READY_TO_SUBMIT` になる | 「提出準備中」 |
| サブスク | 永遠に `MISSING_METADATA` | 「提出準備中」+「審査用に追加」ボタン有効 |

UI が正で、API の `state` はサブスクでは当てになりません。これを知らずに「まだ何か足りない」と探し続けました（実際には何も足りていなかった）。

UI ではこう見えます。「下書き」の中に全商品が並び、ステータスが黄色い「提出準備中」で揃った状態が**正常**です。

![App Store Connect のアプリ内購入一覧。「下書き（6）」の中に非消費型の商品 6 件が並び、ステータスがすべて黄色い「提出準備中」になっている](./assets/asc-iap-waiting.png)

「提出準備中」は不備ではなく、**初回は新しいアプリバージョンと一緒に提出する**ためにこの状態で待つのが仕様です。

時間を置けば変わる、でもありません。540 回ポーリングしても不変でした。Sandbox 購入が通っている = 実体は完成しています。

> サブスクの完成判定に API の `state` を使わない。

### ポーリングの終了条件が「空応答」でも成立してしまう

上のポーリングを書いたときの自作スクリプトの穴です。

```python
# 悪い終了条件
{p: state for p in products if state == "MISSING_METADATA"}  # が空になったら完了
```

API が一瞬 `data: []` を返した回でもこの条件が成立し、false green で「完了」と報告しかけました。

> 終了条件は「消えたこと」ではなく「期待する値が期待する件数だけ揃ったこと」で書く（`len(ready) == 2` など）。

### その他の実際に確認メモ

- サブスクは価格を「配信地域すべて」に入れる必要がある（JPN だけでは不足）。JPN の price point → `/v1/subscriptionPricePoints/{id}/equalizations?filter[territory]=USA` で Apple の等価価格を引ける（¥280 → $1.99 / ¥2,800 → $17.99）
- `MISSING_METADATA` を抜けるには「availability + 審査用スクリーンショット」の両方が必要（非消費型で試して確認。片方だけでは変わらない）
- 審査用スクショの有無は `images` ではなく `appStoreReviewScreenshot` で引く。`images` で引くと 0 件に見えて「未投入」と誤読します（実際に踏んだ）
- Sandbox テスターは API で作れない。`POST /v2/sandboxTesters` → 403。ASC の ユーザーとアクセス → Sandbox → テスター から手作業。Gmail の `+alias` で作成できた
- ただし「購入履歴のクリア」は API でできる（`POST /v2/sandboxTestersClearPurchaseHistoryRequest`）。「権利を 1 つも持たない状態」を作れるので、復元系の分岐を実機で検証するのに要ります（Play 側は同じ状態を作れなかったので、この検証は iOS で行いました）
- `POST /v1/inAppPurchasePriceSchedules` の `included` の id は **`${local-id}` 形式**でないと 409（`p1` などは弾かれる）
- 審査用スクショはシミュレータから撮った（テーマ = ストア画面 / Pro = Paywall）。Debug ビルドはモック価格が必ず出るので通常導線で撮れる

### 残るのは「最初のサブスクは version と一緒に提出する」

ASC の画面に出るこのバナーです。

> 最初の自動更新サブスクリプションは新しいアプリバージョンとともに提出する必要があります

= アプリ審査の提出時にサブスクを version に紐付ける必要があります。IAP 単体では審査に出せません。

---

## Billing Grace Period — 不正値 PATCH で enum を列挙させる

サブスクの支払い猶予期間の設定です。ここで**非破壊で API を調べる手口**を確立できたので共有します。

### enum の有効値は「不正値を PATCH してエラーを読む」と分かる

App Store Connect API は enum 違反を `ENTITY_ERROR.ATTRIBUTE.TYPE` で返し、`detail` に有効値を全列挙します。値が不正なので**設定は変更されません** = 安全に調べられます。

```
HTTP 409  ENTITY_ERROR.ATTRIBUTE.TYPE
'BOGUS_VALUE' is not a valid value for the attribute 'duration'.
Expected one of: 'THREE_DAYS', 'SIXTEEN_DAYS', 'TWENTY_EIGHT_DAYS'
```

公式ドキュメントを探すより速く、しかも「その環境で実際に受理される値」が分かります（ドキュメントは古い / 地域差があることがある）。GET では enum が分からないので、この手が必要です。

### 「7日」は Apple の選択肢に存在しなかった

私のタスクリストには長く「Grace Period を 7 日に設定」と書かれていました。上の手法で実際に確認した結果がこれです。

| 属性 | 有効値（実際に確認したもの） |
|---|---|
| `duration` | `THREE_DAYS` / `SIXTEEN_DAYS` / `TWENTY_EIGHT_DAYS` |
| `renewalType` | `ALL_RENEWALS` / `PAID_TO_PAID_ONLY` |

「7日」の出所は Google Play 側の既存設定（`P7D`）でした。Android に揃えるつもりで書かれた項目が、Apple では実現不能だったわけです。

Play の制約は離散選択ではなく `P0D`〜`P30D` の範囲でした（`P999D` を投げて `Grace period must be between P0D and P30D, inclusive` を得た）。だから **Play を Apple 側に寄せる**ことで両ストアを 16 日で一致させられました。

> チェックリストの文言が実現不能なことがある。着手前に API で選択肢を確認する。

### Apple は prod だけ有効にできない

```
optIn: true + sandboxOptIn: false
→ 409 STATE_ERROR.SUBSCRIPTION_GRACE_PERIOD_OPT_IN_SETTING_INVALID
   "cannot be enabled in prod and disabled in sandbox"
```

**両方 true が必須**です。言い換えると **Sandbox でも猶予期間が効く**ので、課金検証のシナリオがその分変わります。

### 猶予期間より先は揃えられない

Play は猶予後に「**アカウントの保留（`P30D`）**」がありますが、Apple にこの概念はなく Billing Retry が最大 60 日動きます。**両 OS で揃えられるのは猶予期間だけ**です。

> **補足**: リソース ID は appId と同一（`/v1/subscriptionGracePeriods/1234567890`）。grace period は**アプリ単位**で、サブスクグループ単位ではありません。

### この手法は「エンドポイントの存在確認」にも使える

未作成の to-one リレーションシップは GET が `data: null` を返します（404 ではない）。

```
GET /v1/appStoreVersions/{id}/appStoreReviewDetail
→ {"data": null}
```

**これを「API 非対応」と誤読しやすい**（逆に `appAvailabilityV2` は未設定だと 404 を返す、という非対称もある）。

属性を空にした POST を投げると、こう返ります。

```
409 ENTITY_ERROR.RELATIONSHIP.REQUIRED
    You must provide a value for the relationship 'appStoreVersion'
```

何も作成されないまま、エンドポイントの存在だけが確定します。

### この手法自体が壊れていた

余談ですが、私の API ラッパーは `ex.read().decode()[:500]` でエラー出力を切っていました。そのため 1 レスポンスに複数 `errors` が並んだときに 2 件目以降が読めていませんでした。

`appStoreReviewDetails` の POST は `contactPhone` について `ATTRIBUTE.REQUIRED` と `ATTRIBUTE.INVALID` の **2 件**を返すのに後者が消え、「電話を足せば通るのか、形式も違うのか」が判断できませんでした。

エラー本文を切り詰めるロギングは、この種の調査を静かに無効化します。

---

## メタデータと英語ロケール

配信が日本 + 米国なので、ja と en-US の 2 ロケールを埋めました。

### 英語ロケールを作った瞬間に、提出条件が増える

最初、`appInfoLocalizations` は `ja` だけでした。プライバシーポリシー URL の英語版を入れるために `en-US` を新規作成したところ、提出時に en-US の説明文 / キーワード / スクリーンショットもすべて必須になりました。

しかも当時は **ja の説明文すら未入力**でした。「英語 URL を 1 つ入れる」つもりの作業が、掲載文 4 面（ASC の ja / en-US、Play の ja-JP / en-US）とスクショの英語版まで巻き込んだわけです。

> ロケールの追加は「1 項目の追加」ではなく「提出要件の追加」。

### ASC は罫線素片 U+2500 を禁止文字として弾く

説明文の区切り線に `─`（U+2500）を使っていたら 409 になりました。

```
ENTITY_ERROR.ATTRIBUTE.INVALID.INVALID_CHARACTERS
```

`■` / `【】` / `※` は問題なく通ります。**装飾記号は使う前に 1 回投げて確かめる**のが早いです。

### en-US のバージョンローカライズは POST すると 409

`appInfoLocalizations` に `en-US` を作った時点で、バージョン側（`appStoreVersionLocalizations`）も自動生成されていました。なので POST ではなく **PATCH が正解**です。

### 掲載文の方針で決めたこと

内容そのものは各アプリの話ですが、判断の型としてメモしておきます。

- iOS 用と Android 用で本文を分けた。Siri / Spotlight / iCloud は Android 非対応なので、Android 側は「ウィジェット / ランチャー長押し / 他アプリからの共有」「復元専用バックアップ」に差し替えた
- Free / Pro の記述に誤りがあったのを実装から試して修正した。当初 Pro 欄に書いていた機能が実際は Free だった（= 誤った有料表示）
- 価格は設定値を API で自分で動かして書いた。Apple Guideline 3.1.2 が自動更新サブスクの「名称・期間・価格」明記を要求するので省略できません
- 文案は**機械検査した**（文字数 / 片 OS 専用語の混入 / 価格明記 / 日本語欄への英単語混入）

---

## 年齢レーティングとプライバシー申告

### 17+ になっていた真因はアプリ内ブラウザだった

仕様書には 2 か所で「このアプリは 4+」と書いてあったのに、実際の申告は **17+** になっていました。しかも Play 側は ESRB E / 3+ で、両ストアで食い違っていました。

真因は iOS だけが `SFSafariViewController` で任意 URL をアプリ内に開いていたことです。Apple はこれを Unrestricted Web Access と見なします。Android は OS のブラウザに渡していたので 3+ で通っていました。**片 OS だけの非対称**です。

3 箇所を `openURL`（外部ブラウザに渡す）に統一した結果:

```
SEVENTEEN_PLUS → FOUR_PLUS
ブラジル SIXTEEN → L
フランス EIGHTEEN → なし
```

`socialMedia = true` でも 4+ のままでした。Apple の "first appears at the 13+ level" という表現はレーティングの床ではありません。

再発防止として「アプリ内ブラウザを持ち込んだら CI が落ちる」ガードを入れました（Android 側は `androidx.browser` が Gradle 依存に入っただけでも落とす）。

### contentRightsDeclaration は提出前必須項目

`contentRightsDeclaration` が未申告（`None`）のままでした。**提出前必須項目**です。

該当したのは以下でした。

- **ブランドロゴ 7 種**（GitHub / Google / Google Calendar / Google Sheets / Instagram / LINE / X）
- **関連 URL から取得した OG 画像**（URL を DB に保存して他ユーザーにも表示されるため）

→ `USES_THIRD_PARTY_CONTENT` を設定しました。

この確認の過程で、**同梱アセットの出所台帳**（画像 81 件の出所を 1 ファイルに記録）も作りました。ブランドロゴについては、こう整理しました。

- 「当時どこから落としたか」の記録が残っていなかったので、「今の同梱物が公式配布物と一致するか」で担保する方針に切り替えた
- X / LINE は公式配布物と **sha256 完全一致**を確認。LINE は PNG を直接配布せず ZIP のみなので、そこまで辿らないとガイドラインの「公式サイトからダウンロード」要求を満たす証拠にならない
- GitHub のロゴは公式 octicons と **SVG path がバイト完全一致**
- Google 系は公式配布物と**色構成が完全一致**（着色改変なしの証明）

「permitted logo で連携を示す」用途はガイドラインが明示的に許可しているので、**無改変であることを証明できれば問題ない**という整理です。

### App Privacy（プライバシー栄養ラベル）は実装から起こす

マニフェストが **2 種類しか宣言していない**のに、実際はもっと多くのデータを外部に送っていました。実際に確認して **9 種**に是正しました。

| 分類 | 種別 |
|---|---|
| 関連付けあり（7種） | 名前 / メールアドレス / 写真またはビデオ / その他のユーザコンテンツ / ユーザ ID / デバイス ID / 購入 |
| 関連付けなし（2種） | 製品の操作（アナリティクス）/ クラッシュデータ |
| トラッキング | 全 9 種で「なし」 |

判断の根拠を試した結果作った点がポイントです。

- **非リンクの根拠は実際に確認**: `setUserID` / `setUserId` / `setUserProperty` が iOS も Android も 1 箇所も無いので、アナリティクスとクラッシュログは UID に紐づかない
- **支払い情報は申告しない**（StoreKit / Play Billing が処理し、アプリは保持しない）
- ASC の栄養ラベルは API では扱えない（`appDataUsages` という relationship は**存在しない**ことを確認しました）。目視確認しかないので、プレビューが「関連付けあり N カテゴリ / なし M カテゴリ / トラッキング欄なし」になることで検証した

### ブロック / 報告の要件はレーティングではなく審査の合否

UGC（ユーザー生成コンテンツ）を扱うアプリには報告・ブロック機能が必須ですが、iOS の `ageRatingDeclaration` には全 29 項目のどこにもブロック / 報告 / モデレーションの設問がありません（すべて content 系 + 機能開示）。

Apple のブロック・報告要件は Guideline 1.2（審査の合否）であって、レーティングの入力項目ではないということです。実装しても年齢レーティングは変わりません。

---

## 審査ノート（App Review Information）

### Console 手動だと思っていたが API で作れた

`appStoreReviewDetails` は API で作成できます。前述のとおり、未作成のときの GET が `data: null` を返すだけなので「API 非対応」と誤読していました。

```
POST /v1/appStoreReviewDetails        # appStoreVersion に紐付け
GET  /v1/appStoreVersions/{id}/appStoreReviewDetail   # 反映を検証
```

- `notes` は **4,000 字上限**（3,938 字で投入）
- `contactPhone` は必須 + `+` と国番号から始まる形式が必須（`080…` は不可 / `+81 80 …` は可）。省くと `ATTRIBUTE.REQUIRED` と `ATTRIBUTE.INVALID` の 2 件が返る
- 審査連絡先は推測しない。`/v1/users` で `ACCOUNT_HOLDER` の `firstName` / `lastName` / `username`（= メール）が読める

### 導線は「英語 UI の実ラベル」で書く

これは効きます。

レビューノートに書く操作手順は、アプリの英語ローカライズから実際のラベルを引用して書きました。日本語のラベルで書くと、英語 UI で審査するレビュアーが画面上の文字と照合できません。

### 9 セクション構成（そのまま使える骨格）

書いてみて分かったのですが、レビューノートは「レビュアーが詰まりそうな順」に並べると書きやすいです。3,938 字（上限 4,000 字）で次の 9 節にしました。

| # | 節 | 何を書くか |
|---|---|---|
| 1 | サインイン | 認証方法と、デモアカウントが不要である理由（`demoAccountRequired: false` と整合させる） |
| 2 | 特殊な仕様の先出し | レビュアーが「バグ」と誤解しうる挙動（私は単一端末リース）を先に説明して「意図した動作」と明記 |
| 3 | UGC の安全機能（Guideline 1.2） | 規約同意 / 報告 / ブロックの具体的な導線と、1 アカウントで確認できる経路 |
| 4 | ソーシャル機能 | 何が 1 アカウントでできて、何に 2 アカウントが必要かを切り分けて書く |
| 5 | 課金 | 商品名・期間・価格・上限の変化。「Pro に含まれないもの」も明記（別売りの買い切りがある場合） |
| 6 | 任意の外部連携 | 「審査に不要」と明記（連携しないと機能が確認できないと誤解されないように） |
| 7 | 通知 | ローカル通知とリモート通知の区別。起動時に許可を求めない設計であること |
| 8 | ウィジェット / ショートカット | 何が動いて、何が非対応か |
| 9 | データとプライバシー | 保存先とアップロード範囲、削除機能の場所、規約・ポリシーへの導線 |

**連絡先の氏名は英語表記で入れました**（アカウントの登録は日本語表記ですが、**レビュアーが読む欄**なので）。

### 意図的に入れた 3 点

私が「これを書いておかないと落とされる」と判断して入れたのは次の 3 つです。

1. **単一端末リースの説明** — レビュアーが 2 台目で「この端末はロックされています」に当たっても審査ブロックにならないように、復帰操作を案内した
2. Guideline 1.2 の報告 / ブロック / 規約同意を 1 アカウントで確認できる導線 — フレンド機能は双方の承認が必要で 2 アカウント要るので、安全機能は公開リスト経由で検証してもらう手順にした
3. **音声アシスタントでの発話が非対応である旨** — 実装上の裁定を先に伝えて「機能しない」と判定されるのを防ぐ

**3 は特に重要**でした。「実装したが仕様として非対応にした」ものは、黙っていると欠陥と見なされます。

---

## 提出 — API から出すと不足項目が全部返ってくる

これがこの記事で一番役に立つ話かもしれません。

提出は 3 手順です。

```
POST  /v1/reviewSubmissions           # app に紐付け
POST  /v1/reviewSubmissionItems       # appStoreVersion に紐付け
PATCH /v1/reviewSubmissions/{id}      # {"submitted": true}
```

このうち 2 番目が 409 を返し、`meta.associatedErrors` に「何が足りないか」を全部列挙します。

```
/v1/appInfos/… : ENTITY_ERROR.RELATIONSHIP.REQUIRED   primaryCategory 未設定
/v2/appPrices/ : STATE_ERROR.APP_PRICING_REQUIRED     価格未設定
```

どちらも ASC の UI では「空」だと気づけませんでした。

- **プライマリカテゴリ未設定** — 画面上は選択欄が空でもエラー表示が出ない
- **価格未設定** — 無料アプリでも「無料」を明示的に設定する必要がある

> 提出は API から試すのが、穴の見つけ方として優れている。
> UI をどれだけ眺めても分からない不足が、409 で名前付きで返ってきます。

対処はこうでした。

```
PATCH /v1/appInfos/{id}          primaryCategory = LIFESTYLE
POST  /v1/appPriceSchedules      無料（baseTerritory=JPN + customerPrice=0 の appPricePoint）
```

価格ポイント ID は `/v1/apps/{id}/appPricePoints?filter[territory]=JPN` から引きます。

**カテゴリは Play 側と揃えました**（Play の「ライフスタイル」に対して `LIFESTYLE`）。両ストアで違うカテゴリを申告すると、後の申告更新で食い違いの原因になります。

### 提出前に ipa を実際に確かめる

提出ボタンを押す前に、成果物が本当に本番構成かを ipa に対して確認しました。

| 項目 | 期待値 |
|---|---|
| 署名 | `Apple Distribution: …` |
| `get-task-allow` | `False`（デバッガをアタッチできない） |
| `aps-environment` | `production` |
| クラウドコンテナの環境 | `Production`（Development のまま出すと本番データが見えない） |
| `ProvisionedDevices` | 0 件（配布用プロファイル） |
| ビルド番号 / 表示名 | 意図した値 |

クラウド同期の環境（`icloud-container-environment`）は特に注意が必要でした。entitlements に明示していないと**署名に従って決まる**ので、手元の開発用署名で入れた実機では Development、TestFlight 配布では Production になります。本番のクラウド同期は TestFlight でしか検証できません。

### 「承認 → 自動公開」にしない

両ストアを同日に公開したかったので、**承認されても自動で公開されない**設定にしました。

|  | 公開の抑止 |
|---|---|
| iOS | `releaseType = MANUAL` |
| Play | 「管理対象の公開」を ON |

これで両ストアとも「承認 → 待機 → 手動で公開」になります。片方が先に通っても勝手に出ないので、**揃ってから 2 つボタンを押すだけ**にできました。

### 両ストアでビルド番号を揃えた

iOS と Play で **同じビルド番号（コミット数）** で提出しました。そのためには**提出まで 1 コミットもしない**必要があります（コミット数がビルド番号なので）。地味ですが、両ストアの提出物を後から突合するとき、番号が揃っていると圧倒的に楽です。

---

## 提出直前にやったこと

「提出前に毎回確認する」性質のものは、手作業のチェックリストではなく**1 コマンドのゲート**に落としました（読み取り専用・14 チェック・異常時 exit 1）。検証で端末を触ると状態が変わるので、**再提出のたびに再発しうる**からです。

以下、提出直前に見つかったものです。

### 最大 Dynamic Type で「離脱不可の画面」が詰んでいた

これは**発見が遅れたら確実にリジェクトされていた**ブロッカーです。

```sh
xcrun simctl ui <device> content_size accessibility-extra-extra-extra-large
```

これで試してみると、**離脱経路のない全画面**が最大文字サイズで詰んでいました。

| 画面 | 症状 |
|---|---|
| 規約同意シート | 素の `VStack` + `Spacer()` で CTA が画面外。`.interactiveDismissDisabled(true)` なのでアプリが完全に詰む |
| 端末ロック画面 | 同型（root ごと差し替えなので逃げ道なし） |
| オンボーディングの規約ページ | 下部確保が固定 pt で本文が CTA の下に潜る |
| Paywall | すでに `ScrollView` を持っていて無傷 |

修正は 3 画面とも同じ形です。

```swift
GeometryReader { proxy in
    ScrollView {
        content
            .frame(minHeight: proxy.size.height, alignment: .center)
    }
    .scrollBounceBehavior(.basedOnSize)
}
```

通常サイズでの中央寄せを維持しつつ、大きい文字ではスクロールできるようになります。

> 「離脱できない画面」だけは最大 Dynamic Type で必ず手を動かして確認する。逃げ道のある画面は多少崩れても致命的になりません。

### シミュレータの defaults を汚すと UI テストが落ちる

上の検証で規約モーダルを出すために UserDefaults を書き換えたら、**UI テストが落ちました**（`-skipOnboarding` 起動でも離脱不可モーダルが被る）。

検証が終わったら必ず元に戻します。`terms.acceptedVersion` / `onboarding.hasFinished` を元に戻し、デバッグキーを削除しました。

プログラムからの任意タップができないので、モーダルは UserDefaults 操作で自動表示させて撮るのが確実です。ディープリンク（`simctl openurl`）はシステム確認ダイアログが被るので画面撮影には使えません。

### 審査導線が「実データに依存」していた

レビューノートに「公開リストのタブから報告 / ブロックを確認できる」と書いたのに、本番の公開リストが実際に確認 0 件でした。

理由は 2 つの重なりです。

- 公開リストは**端末言語と一致する投稿だけ**を出す仕様（言語フォールバックは意図的に持たない設計）
- **自分の投稿は自己除外される**

「他人名義の公開投稿」が本番に 1 件もないと、タップできる行が存在せず、`…` メニューが一度も描画されない = 報告 / ブロック導線に到達できない = Guideline 1.2 の定型リジェクトになります。

対処として、審査用に「他人名義の公開投稿」を ja / en 各 1 件、本番に常設しました。そして、同時に進めていた「テストデータの掃除」と正面衝突しました。掃除対象に入っていた公開投稿を消すと導線が消えるので、**常設する分を掃除対象から明示的に除外**し、レビューノートにその投稿名を書きました。

提出直前に審査アカウントでサインインして、公開リストが実際に何件返すかを目視する——これ自体を受け入れ条件にしました。0 件ならリジェクト確定なので。

### 検証で触った端末が審査アカウントのセッションを握っていた

単一端末リースの設計だと、**最後に審査アカウントで触った端末**がアクティブ端末として残ります。この状態でレビュアーがサインインすると、**最初に見る画面がロック画面になります**（復帰操作はできるので機能的には詰みませんが、印象が悪い）。

リースは端末単位なので、検証のたびに再発します。だから「提出直前に毎回確認する」項目にしました。

### 提出時に触ってはいけないもの

私が「審査中は絶対に触らない」とメモしたのは次の 3 つです。

1. **審査導線に必要な本番データ**（上記の公開投稿）
2. **審査アカウントへの権利付与の設定**（外すと「有料コンテンツも審査できる」という申告が偽になる）
3. **審査アカウントで端末にサインインしない**（リースが移ってレビュアーがロック画面に当たる）

---

## Apple から届くメール一覧

「いま何が終わったのか」はメールでしか分からない場面が多いので、提出までに実際に届いたものを順番に並べておきます。判定に使えるものと、無視してよいものが混ざっています。

| 順 | 件名 | 送信元 | 届くタイミング |
|---|---|---|---|
| 1 | ご注文ありがとうございます（+ 注文番号） | `order_acknowledgment@orders.apple.com` | Developer Program の支払い直後。**加入完了ではありません** |
| 2 | Apple Developer Programへようこそ | `developer@email.apple.com` | 加入がアクティベートされた時 |
| 3 | Welcome to App Store Connect. | `no_reply@email.apple.com` | 2 と同時。**この 2 通が揃って初めて有効な有料メンバーシップ** |
| 4 | App Store Connect API Access Request Approved | `no_reply@email.apple.com` | API キーの利用申請が通った時 |
| 5 | Apple Duns Update | `noreply-appledev@email.apple.com` | D&B のデータが Apple 側に反映された時 |
| 6 | Action needed: The uploaded build for〈アプリ名〉has one or more issues. | `no_reply@email.apple.com` | **アップロードは成功したのに処理で弾かれた時** |
| 7 | 〈アプリ名〉1.0.0 (〈ビルド番号〉) for iOS is now available to test. | `testflight_no_reply@email.apple.com` | TestFlight に配信された時。ビルドごとに 1 通 |
| 8 | Thank You for Submitting Your App | `no_reply@email.apple.com` | 審査に提出した直後 |

読み方で注意した点が 3 つあります。

- **1 は加入完了ではありません。** 注文確認だけで、2 と 3 が来るまでメンバーシップは有効になりません。私は 1 を見て「入れた」と判断しかけました
- **6 が「アップロード失敗」と紛らわしいです。** `altool` は `UPLOAD SUCCEEDED` を返していて、転送は本当に成功しています。その後の処理で弾かれたことがメールでしか通知されません。API のビルド一覧が 0 件のままだった真因はこれでした
- **7 はビルド番号が入っているので、どの版が配信されたかの判定に使えます。** Console の表示より確実です

審査に出した後は 8 が届きます。その先の審査結果メールについては、結果が出たので別記事にまとめます。

---

## 参考リンク

画面名で探すと迷うものが多かったので、実際に使った入口を置いておきます。**すべて公開時点で開けることを確認しました**（コンソール系はサインインへリダイレクトされます）。

| 用途 | URL |
|---|---|
| 「ビジネス」（有料アプリ契約 / 税務情報 / 銀行口座） | <https://appstoreconnect.apple.com/business> |
| API キーの発行（Integrations → App Store Connect API） | <https://appstoreconnect.apple.com/access/integrations/api> |
| 配布証明書の一覧 | <https://developer.apple.com/account/resources/certificates/list> |
| Identifiers（App ID / App Group / iCloud Container） | <https://developer.apple.com/account/resources/identifiers/list> |
| デバイスの登録（実機の UDID） | <https://developer.apple.com/account/resources/devices/add> |
| 年齢レーティングの各値の定義 | <https://developer.apple.com/help/app-store-connect/reference/app-information/age-ratings-values-and-definitions/> |
| 配信する国や地域の管理 | <https://developer.apple.com/help/app-store-connect/manage-your-apps-availability/manage-availability-for-your-app/> |
| アプリ内課金を審査に提出する | <https://developer.apple.com/help/app-store-connect/manage-submissions-to-app-review/submit-an-in-app-purchase> |
| プライバシーマニフェスト | <https://developer.apple.com/documentation/bundleresources/privacy-manifest-files> |
| App Store Review Guidelines | <https://developer.apple.com/app-store/review/guidelines/> |
| Apple Developer Program の加入 | <https://developer.apple.com/programs/enroll/> |
| D-U-N-S を無料で取る専用ページ | <https://developer.apple.com/enroll/duns-lookup/> |

---

## まとめ — Apple 側チェックリスト

順序つきで置いておきます。

```
【契約・名義】
□ Individual / Organization を決める（Individual は本名が公開される・変更不可）
□ D-U-N-S が必要なら Apple の無料ルートを先に確認する
□ 屋号・住所のローマ字表記を 1 か所に固定する（4 か所で完全一致）
□ Apple Developer Program に加入 → メール 2 通のペアで有効化を確認
□ Team ID を Membership 詳細で確認（登録ID / 注文番号とは別物）

【identifier・署名】
□ App Group / iCloud Container を App ID より先に作る（グローバル一意）
□ archive → export（destination=export）で発行 → ipa で署名を実測
□ Cloud Managed なので .p12 バックアップは不要と理解する

【契約・税務】
□ App Store Connect → 「ビジネス」で有料アプリ契約を有効化（IAP はこれが前提）
□ W-8BEN: Line 10 は空欄 / 6a に Foreign TIN / Line 8 は MM-DD-YYYY
□ 銀行口座: 国の既定を直す / 口座名義は銀行の表記どおり / 種類は「普通」
□ 配信地域を絞って EU DSA 申告を回避（availableInNewTerritories = false）

【アプリレコード】
□ Web UI でアプリを作成（API では不可）→ 名前が予約される
□ CFBundleName を表示名として明示設定する
□ contentRightsDeclaration を設定（提出前必須）
□ プライマリカテゴリを設定（UI では空でも気づけない）
□ 価格を設定（無料でも明示が必要）

【課金】
□ availability を価格より先に作る
□ サブスクの state は MISSING_METADATA のままが正常
□ 審査用スクショは appStoreReviewScreenshot で確認する
□ ファミリー共有 / groupLevel は不可逆な選択と理解する
□ Sandbox テスターは画面から手作業で作る
□ Grace Period の選択肢を不正値 PATCH で確認する

【ビルド配信】
□ 転送 / 受理 / 配信の 3 段をそれぞれ判定する
□ ビルド存在確認は filter[version] で引く（limit の一覧は使わない）
□ ベータグループへの割当を自動化する（後から自動配信には変えられない）
□ 拡張ターゲットのバージョンを本体と一致させる

【申告】
□ App Privacy を実装から起こす（マニフェストの宣言を信じない）
□ 年齢レーティング: アプリ内ブラウザがあると 17+ になる
□ 審査ノートは英語 UI の実ラベルで書く / contactPhone は +国番号

【提出直前】
□ 提出を API から試して 409 の associatedErrors を読む
□ 最大 Dynamic Type で離脱不可画面を実測する
□ 審査導線が実データに依存していないか目視する
□ 検証で触った端末のセッションを解放する
□ 提出直前チェックを 1 コマンドのゲートに落とす
```

Apple 側で一番の学びは、「成功レスポンスと画面表示のどちらも信用できない」ということでした。`UPLOAD SUCCEEDED` は受理ではなく、`VALID` は配信ではなく、サブスクの `MISSING_METADATA` は不備でもありません。**判定は必ず「その工程の実データ」で行う**——これが唯一まともに機能した方針です。

Google Play 側は、また別の種類の罠がありました。そちらは次の記事にまとめています。

## 関連記事

このシリーズの他の記事です。

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

## アプリそのものについて

作ったアプリの中身と、実装で踏んだ罠は別のシリーズにまとめました。

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
